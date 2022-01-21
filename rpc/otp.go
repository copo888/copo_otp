package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"

	"github.com/copo888/copo_otp/rpc/internal/config"
	"github.com/copo888/copo_otp/rpc/internal/server"
	"github.com/copo888/copo_otp/rpc/internal/svc"
	"github.com/copo888/copo_otp/rpc/otp"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	configFile = flag.String("f", "etc/otp.yaml", "the config file")
	envFile    = flag.String("env", "etc/.env", "the env file")
)

func main() {
	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	var c config.Config
	conf.MustLoad(*configFile, &c,  conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	srv := server.NewOtpServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		otp.RegisterOtpServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
