package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"
)

// Go Web开发比较通用的脚手架模板

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings error: %s\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Init logger error: %s\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3. 初始化Mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init logger error: %s\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("Init logger error: %s\n", err)
		return
	}
	redis.Close()
	// 5. 注册路由
	r := routes.Setup()
	// 6. 启动服务（优雅关机）
	srv := http.Server{
		Addr:    fmt.Sprintf("%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		//
		if err := srv.ListenAndServe(); err != nil && err == http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}

		// 等待中断信号优雅关闭，设置5秒超时
		quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
		// kill 默认会发送 syscall.SIGTERM 信号
		// kill -2 发送 syscall.SIGINT 信号（如 ctrl+c）
		// kill -9 发送 syscall.SIGKILL 信号，信号无法被捕获所以无需添加
		// signal.Notify 把收到的 syscall.SIGTERM 或 syscall.SIGINT 信号转发给quit
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit // 阻塞，接受信号后继续执行
		zap.L().Info("Shutdown Server ...")
		// 创建一个5s超时的context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			zap.L().Fatal("Server Shutdown", zap.Error(err))
		}

		zap.L().Info("Server exiting")
	}()
}
