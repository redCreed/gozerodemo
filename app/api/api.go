package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gozerodemo/app/api/internal/config"
	"gozerodemo/app/api/internal/handler"
	"gozerodemo/app/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api-api.yaml", "the config file")

// 默认输出level为stat的统计日志，用于输出当前资源的使用情况，主要为cpu和内存  可以关闭
func init() {
	//close statis log
	logx.DisableStat()
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
