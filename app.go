package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	httpServer *http.Server
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time! ", name)
// }

// SelectedDirectory 选择目录
func (a *App) SelectedDirectory() (string, error) {
	dialog, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if dialog == "" {
		return "", nil // 用户取消了对话框
	}
	return dialog, nil
}

// 设置草稿根目录
func (a *App) SetDraftRootPath(draftRootPath string) (string, error) {
	config.DraftRootPath = draftRootPath
	err := configStore.SaveConfig(config)

	if err != nil {
		return "", err
	}
	return draftRootPath, nil
}

func (a *App) SendLogsToPage(message string) {
	// json 格式的日志
	log := Log{
		Type:    "log",
		Message: message,
	}
	jsonData, err := json.Marshal(log)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	runtime.EventsEmit(a.ctx, "logs", string(jsonData))
}
