package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	osruntime "runtime"

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

// 自动检测草稿目录
func (a *App) AutoDetectDraftRootPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}

	var draftPath string
	switch osruntime.GOOS {
	case "darwin":
		draftPath = filepath.Join(homeDir, "Movies", "JianyingPro", "User Data", "Projects", "com.lveditor.draft")
	case "windows":
		draftPath = filepath.Join(homeDir, "AppData", "Local", "JianyingPro", "User Data", "Projects", "com.lveditor.draft")
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", osruntime.GOOS)
	}

	// 检查目录是否存在，并且包含 root_meta_info.json 文件
	if _, err := os.Stat(draftPath); os.IsNotExist(err) {
		return "", fmt.Errorf("草稿目录不存在: %s", draftPath)
	}

	rootMetaInfoPath := filepath.Join(draftPath, "root_meta_info.json")
	if _, err := os.Stat(rootMetaInfoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("草稿目录无效，缺少 root_meta_info.json 文件: %s", draftPath)
	}

	return draftPath, nil
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
