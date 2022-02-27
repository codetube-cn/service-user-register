package bootstrap

import (
	service_user_register "codetube.cn/proto/service-user-register"
	"codetube.cn/service-user-register/server"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var BootErrChan chan error

func Start() {
	defer func() {
		err := recover()
		log.Println(err)
	}()
	BootErrChan = make(chan error)
	//初始化
	//go func() {
	//	err := initApp()
	//	if err != nil {
	//		log.Println("init app fail:", err.Error())
	//		BootErrChan <- err
	//	}
	//}()
	//开启服务
	go func() {
		defer func() {
			err := recover()
			fmt.Println("defer")
			log.Println(err)
		}()
		regServer := grpc.NewServer()
		service_user_register.RegisterUserRegisterServer(regServer, server.NewUserRegisterServer())
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			BootErrChan <- err
		}
		if err = regServer.Serve(lis); err != nil {
			BootErrChan <- err
		}
	}()
	//监听事件
	go func() {
		sigC := make(chan os.Signal)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		BootErrChan <- fmt.Errorf("%", <-sigC)
	}()
	getErr := <-BootErrChan
	log.Println(getErr)
}
