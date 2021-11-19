## micro-discovery
微服务发现插件


#### kratos服务发现插件，适配go-micro

1. 安装
```
go get github.com/yinrenxin/micro-discovery
```
2. example
```
    client, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: time.Second, DialOptions: []grpc.DialOption{grpc.WithBlock()}})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	service := discovery.NewGoMicro(client)

	// 获取服务列表
	r, err := service.GetService(ctx, "serviceName")
	if err != nil {
		panic(err)
	}

	// 初始化kratos的grpc client，传入go-micro微服务的服务名
	conn, err := grpc.DialInsecure(
		ctx,
		grpc.WithEndpoint("discovery:///servicename"),
		grpc.WithDiscovery(service),
	)

	// 调用方法

```