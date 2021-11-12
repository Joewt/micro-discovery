## micro-discovery
微服务发现插件


#### kratos服务发现插件，适配go-micro

1. 安装
```
go get github.com/yinrenxin/micro-discovery/kratos
```
2. example
```
    client, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	service := kratos.NewGoMicro(client)

	r, err := service.GetService(ctx, "serviceName")
	if err != nil {
		panic(err)
	}
```