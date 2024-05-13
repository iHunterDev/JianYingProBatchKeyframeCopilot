package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

// 启动 http 服务
func (a *App) StartHTTPServer() error {
	mux := http.NewServeMux()

	a.HandleFuncWarp(mux)

	// 创建 CORS 处理器
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	a.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: corsHandler.Handler(mux),
	}

	// 异步启动服务
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err)
		}
	}()

	log.Println("Server started. Press Ctrl+C to stop.")

	return nil
}

// 关闭 http 服务
func (a *App) StopHTTPServer() error {
	if a.httpServer == nil {
		return nil
	}

	if err := a.httpServer.Shutdown(a.ctx); err != nil {
		return err
	}

	return nil
}
