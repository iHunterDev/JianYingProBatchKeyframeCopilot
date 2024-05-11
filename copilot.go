package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// 请求前记录日志
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		// 调用下一个处理器
		next(w, r)
	}
}

// 启动 http 服务
func (a *App) StartHTTPServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))

		// json 格式的日志
		runtime.EventsEmit(a.ctx, "logs", `{"type":"log","message":"Hello, World!"}`)

	}))

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
