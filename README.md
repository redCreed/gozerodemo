# gozerodemo

* goctl api go -api api.api -dir .
  使用定义的api文件重新生成新的文件

* goctl rpc protoc order.proto --go_out=. --go-grpc_out=. --zrpc_out=.
  使用定义的proto文件重新生成新的相关代码
* goctl model mysql datasource -url="root:123456@tcp(127.0.0.1:3306)/product" -table="*"  -dir="./model" -c
  缓存model数据的功能