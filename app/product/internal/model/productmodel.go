package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	// ProductModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductModel.
	ProductModel interface {
		productModel
		CategoryProducts(ctx context.Context, ctime string, cateid, limit int64) ([]*Product, error)
		UpdateProductStock(ctx context.Context, pid, num int64) error
		TxUpdateStock(tx *sql.Tx, id int64, delta int) (sql.Result, error)
	}

	customProductModel struct {
		*defaultProductModel
	}
)

// UpdateProductStock 删除数据或者更新数据的时候，以id为key的行记录缓存会被删除
func (c customProductModel) UpdateProductStock(ctx context.Context, pid, num int64) error {
	productProductIdKey := fmt.Sprintf("%s%v", cacheProductProductIdPrefix, pid)
	_, err := c.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, fmt.Sprintf("UPDATE %s SET stock = stock - ? WHERE id = ? and stock > 0", c.table), num, pid)
	}, productProductIdKey)
	return err
}

func (c customProductModel) TxUpdateStock(tx *sql.Tx, id int64, delta int) (sql.Result, error) {
	productIdKey := fmt.Sprintf("%s%v", cacheProductProductIdPrefix, id)
	return c.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set stock = stock + ? where stock >= -? and id=?", c.table)
		return tx.Exec(query, delta, delta, id)
	}, productIdKey)
}

func (c customProductModel) CategoryProducts(ctx context.Context, ctime string, cateid, limit int64) ([]*Product, error) {
	var products []*Product
	err := c.QueryRowsNoCacheCtx(ctx, &products, fmt.Sprintf("select %s from %s where cateid=? and status=1 and create_time<? order by create_time desc limit ?", productRows, c.table), cateid, ctime, limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// NewProductModel returns a model for the database table.
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn, c, opts...),
	}
}
