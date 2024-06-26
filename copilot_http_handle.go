package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func (a *App) HandleFuncWarp(mux *http.ServeMux) {
	a.HandleHome(mux)
	a.HandleDrafts(mux)
	a.HandleDraft(mux)
}

func (a *App) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// 请求前记录日志
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		// 调用下一个处理器
		next(w, r)

		a.SendLogsToPage(r.Method + " " + r.URL.Path)
	}
}

func (a *App) HandleHome(mux *http.ServeMux) {
	mux.HandleFunc("/", a.loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, My Friend!\nThis is the JianyingPro Batch Keyframe Copilot Server.\nPlease visit https://github.com/iHunterDev/JianYingProBatchKeyframeCopilot"))
	}))
}

func (a *App) HandleDrafts(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/drafts", a.loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// 读取 drafts 信息
		file, err := os.Open(path.Join(config.DraftRootPath, "root_meta_info.json"))
		if err != nil {
			fmt.Println("无法打开文件:", err)
			return
		}
		defer file.Close()

		// 读取文件内容
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("无法读取文件内容:", err)
			return
		}

		// 解析 JSON 数据
		var metaInfo RootMetaInfo
		err = json.Unmarshal(content, &metaInfo)
		if err != nil {
			fmt.Println("无法解析 JSON 数据:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Code:    0,
			Message: "success",
			Data:    metaInfo.AllDraftStore,
		}
		jsonData, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		w.Write(jsonData)
	}))
}

func (a *App) HandleDraft(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/draft", a.loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// 读取请求方式
		method := r.Method

		// 判断请求方式
		if method == "GET" {
			DraftInfoAction(w, r)
		} else if method == "POST" {
			DraftSaveAction(w, r)
		}
	}))
}

// 草稿信息读取操作
func DraftInfoAction(w http.ResponseWriter, r *http.Request) {
	// 读取 draft 信息
	draftPath := r.URL.Query().Get("draft_json_file")
	file, err := os.Open(draftPath)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("无法读取文件内容:", err)
		return
	}

	// 修复再 window 下的路径问题
	draftPath = strings.Replace(draftPath, "\\", "\\\\", -1)
	// 输出 draft 信息
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code":0,"message":"success","data":{"draft_json_file":"` + draftPath + `","draft_info":` + string(content) + `}}`))
}

// 草稿信息保存操作
func DraftSaveAction(w http.ResponseWriter, r *http.Request) {
	// 读取 draft 信息
	// 去读请求的 body json 格式数据
	type DraftSaveRequest struct {
		DraftJSONFile string                 `json:"draft_json_file"`
		DraftInfo     map[string]interface{} `json:"draft_info"`
	}
	var data DraftSaveRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println("无法解析 JSON 数据:", err)
		return
	}

	// 把数据转换为 json 字符串
	jsonData, err := json.Marshal(data.DraftInfo)
	if err != nil {
		fmt.Println("无法转换为 json 字符串:", err)
		return
	}

	// 写入文件
	err = os.WriteFile(data.DraftJSONFile, jsonData, 0644)
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return
	}

	// 修复再 window 下的路径问题
	draftPath := strings.Replace(data.DraftJSONFile, "\\", "\\\\", -1)
	// 输出 draft 信息
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code":0,"message":"success","data":{"draft_json_file":"` + draftPath + `"}}`))
}

//+--------
// types
//+--------

type Log struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RootMetaInfo struct {
	AllDraftStore []Draft `json:"all_draft_store"`
	DraftIds      int     `json:"draft_ids"`
	RootPath      string  `json:"root_path"`
}

type Draft struct {
	DraftCloudLastActionDownload   bool   `json:"draft_cloud_last_action_download"`
	DraftCloudPurchaseInfo         string `json:"draft_cloud_purchase_info"`
	DraftCloudTemplateID           string `json:"draft_cloud_template_id"`
	DraftCloudTutorialInfo         string `json:"draft_cloud_tutorial_info"`
	DraftCloudVideocutPurchaseInfo string `json:"draft_cloud_videocut_purchase_info"`
	DraftCover                     string `json:"draft_cover"`
	DraftFoldPath                  string `json:"draft_fold_path"`
	DraftID                        string `json:"draft_id"`
	DraftIsAIShorts                bool   `json:"draft_is_ai_shorts"`
	DraftIsInvisible               bool   `json:"draft_is_invisible"`
	DraftJSONFile                  string `json:"draft_json_file"`
	DraftName                      string `json:"draft_name"`
	DraftNewVersion                string `json:"draft_new_version"`
	DraftRootPath                  string `json:"draft_root_path"`
	DraftTimelineMaterialsSize     int    `json:"draft_timeline_materials_size"`
	DraftType                      string `json:"draft_type"`
	TMDraftCloudCompleted          string `json:"tmdraft_cloud_completed"`
	TMDraftCloudModified           int    `json:"tmdraft_cloud_modified"`
	TMDraftCreate                  int    `json:"tmdraft_create"`
	TMDraftModified                int    `json:"tmdraft_modified"`
	TMDraftRemoved                 int    `json:"tmdraft_removed"`
	TMDuration                     int    `json:"tmduration"`
}
