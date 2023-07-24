package logic

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
	"gozerodemo/app/order/internal/model"
	"gozerodemo/app/product/product"
	"gozerodemo/app/user/user"
	"gozerodemo/app/user/userclient"
	"gozerodemo/pkg/snowflake"
	"gozerodemo/pkg/xerr"

	"gozerodemo/app/order/internal/svc"
	"gozerodemo/app/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderDTMLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderDTMLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderDTMLogic {
	return &CreateOrderDTMLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderDTMLogic) CreateOrderDTM(in *order.AddOrderReq) (*order.AddOrderResp, error) {
	var (
		userRpcRes        *user.UserInfoResponse
		productRpcRes     *product.ProductItem
		receiveAddressRes *user.UserReceiveAddress
	)

	//check product
	checkProduct := func() error {
		var err error
		var productReq product.ProductItemRequest
		productReq.ProductId = in.Productid
		productRpcRes, err = l.svcCtx.ProductRpc.Product(l.ctx, &productReq)
		if err != nil {
			return nil
		}
		return nil
	}

	//check user
	checkUser := func() error {
		var err error
		var userReq user.UserInfoRequest
		userReq.Id = in.Userid
		userRpcRes, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &userReq)
		if err != nil {
			return nil
		}
		return nil
	}

	//check user_receive_address
	checkUserReceiveAddress := func() error {
		var err error
		var userReceiveAddressInfoReq userclient.UserReceiveAddressInfoReq
		userReceiveAddressInfoReq.Id = in.ReceiveAddressId
		receiveAddressRes, err = l.svcCtx.UserRpc.GetUserReceiveAddressInfo(l.ctx, &userReceiveAddressInfoReq)
		if err != nil {
			return nil
		}
		return nil
	}

	//parallel call
	err := mr.Finish(checkProduct, checkUser, checkUserReceiveAddress)

	if userRpcRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! user not exist exception : %+v  ", userRpcRes)
	}

	if productRpcRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! not exist exception : %+v  ", productRpcRes)
	}

	//check product stock
	if productRpcRes.Stock <= 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("product understock"), "product understock")
	}

	if receiveAddressRes == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "error! user receive address exception : %+v  ", receiveAddressRes)
	}

	//generate new order id
	orderId := snowflake.GenIDString()

	db, err := sqlx.NewMysql(l.svcCtx.Config.DataSource).RawDB()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	// 获取子事务屏障对象
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	// 开启子事务屏障
	if err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		//create new shipping
		var insertShipping model.Shipping
		insertShipping.Orderid = orderId
		insertShipping.Userid = in.Userid
		insertShipping.ReceiverName = receiveAddressRes.Name
		insertShipping.ReceiverPhone = receiveAddressRes.Phone
		insertShipping.ReceiverMobile = receiveAddressRes.Phone
		insertShipping.ReceiverProvince = receiveAddressRes.Province
		insertShipping.ReceiverCity = receiveAddressRes.City
		insertShipping.ReceiverDistrict = receiveAddressRes.Region
		insertShipping.ReceiverAddress = receiveAddressRes.DetailAddress
		insertShippingRes, err := l.svcCtx.ShippingModel.Insert(l.ctx, &insertShipping)
		if err != nil {
			return err
		}
		newShippingId, err := insertShippingRes.LastInsertId()
		if err != nil {
			return err
		}

		//create new orderitem
		insertOrderitem := model.Orderitem{
			OrderId:      orderId,
			UserId:       in.Userid,
			ProductId:    in.Productid,
			ProductName:  productRpcRes.Name,
			ProductImage: productRpcRes.ImageUrl,
			CurrentPrice: productRpcRes.Price,
			Quantity:     in.Quantity,
			//TotalPrice:   float64(mathin.Quantity * productRpcRes.GetPrice()),
		}
		insertOrderitemRes, err := l.svcCtx.OrderitemModel.Insert(l.ctx, &insertOrderitem)
		if err != nil {
			return err
		}
		_, err = insertOrderitemRes.LastInsertId()
		if err != nil {
			return err
		}

		//create new order
		insertOrder := model.Orders{
			Id:         orderId,
			Userid:     in.Userid,
			Shoppingid: newShippingId,
			Postage:    in.Postage,
		}
		_, err = l.svcCtx.OrderModel.TxInsert(tx, &insertOrder)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DbError), "create new order Exception : err: %v", err)
	}

	return &order.AddOrderResp{
		Id: orderId,
	}, nil
}
