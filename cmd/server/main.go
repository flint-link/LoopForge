package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"LoopForge/internal/config"
	"LoopForge/internal/database"
	"LoopForge/internal/handler"
	"LoopForge/internal/middleware"
	"LoopForge/internal/repository"
	"LoopForge/internal/service"
	"LoopForge/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "", "Config file path")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 初始化日志
	log, err := logger.NewLogger(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output, cfg.Log.FilePath)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync(log)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db, err := database.NewPostgres(&cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// 依赖注入
	tenantRepo := repository.NewTenantRepository(db.Pool)
	tenantService := service.NewTenantService(tenantRepo, log)
	tenantHandler := handler.NewTenantHandler(tenantService, log)

	// 创建路由
	router := gin.New()

	// 中间件
	router.Use(middleware.RecoveryMiddleware(log))
	router.Use(middleware.LoggerMiddleware(log))

	// 健康检查
	router.GET("/health", handler.HealthCheck)

	// API 路由组
	api := router.Group("/api/v1")
	tenantHandler.RegisterRoutes(api)

	// 启动服务器
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 优雅关闭
	go func() {
		log.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
