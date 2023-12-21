package main

import (
	"context"
	"kube-demo-fe/goland/kubea-demo/config"
	"kube-demo-fe/goland/kubea-demo/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kube-demo-fe/goland/kubea-demo/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

func main() {
	r := gin.Default()
	//初始化配置k8s client
	service.K8S.Init()
	//初始化路由
	controller.Router.InitApiRouter(r)
	//启动 gin server
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen err: %s", err)
		}
	}()
	//优雅退出server
	// 创建一个用于接收信号的channel。如果没有信号，这个channel将一直阻塞。如果有信号，就继续执行
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 设置ctx5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel 用于释放 ctx
	defer cancel()

	// 关闭服务
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("gin server Shutdown: %s", err)
	}
	logger.Info("Gin server退出成功")

}
