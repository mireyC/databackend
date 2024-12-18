package main

import (
	"flag"
	"fmt"
	"geoserver/api/etc"
	"geoserver/api/internal/config"
	"geoserver/api/internal/handler"
	"geoserver/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/geoserver-api-linux.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	//fmt.Printf("config: %+v\n", c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etc.LoadImageAndCreateContainer(c)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
	fmt.Println("server stop")
}
