package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uploadDir   = "./web_data/uploads"
	outputDir   = "./web_data/reports"
	maxLogFileSize  = 10 << 20
	maxArchiveSize  = 30 << 20
	adminPassword   = "admin888"
)

var state = NewAppState()

type UploadResponse struct {
	SessionID string   `json:"session_id"`
	Files     []string `json:"files"`
	IsDir     bool     `json:"is_dir"`
}

type StartAnalysisRequest struct {
	SessionID    string         `json:"session_id"`
	CustomerName string         `json:"customer_name"`
	Options      AnalyzeOptions `json:"options"`
	FontName     string         `json:"font_name"`
}

type TaskResponse struct {
	Task  *AnalysisTask `json:"task"`
	Stats Stats         `json:"stats"`
}

func init() {
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(outputDir, 0755)
}

func handleUpload(c *gin.Context) {
	sessionID := c.PostForm("session_id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	state.RegisterSession(sessionID)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析上传请求失败: " + err.Error()})
		return
	}

	sessionUploadDir := filepath.Join(uploadDir, sessionID)
	os.MkdirAll(sessionUploadDir, 0755)

	uploadedFiles := handleFileUploads(c, form, sessionUploadDir)
	if len(uploadedFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未接收到有效文件"})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		SessionID: sessionID,
		Files:     uploadedFiles,
		IsDir:     false,
	})
}

func handleFileUploads(c *gin.Context, form *multipart.Form, dstDir string) []string {
	files := form.File["files"]
	var uploaded []string
	for _, f := range files {
		if f.Size > maxLogFileSize {
			continue
		}
		dst := filepath.Join(dstDir, f.Filename)
		if err := c.SaveUploadedFile(f, dst); err != nil {
			continue
		}
		uploaded = append(uploaded, f.Filename)
	}
	return uploaded
}

func handleStartAnalysis(c *gin.Context) {
	var req StartAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}
	if req.SessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 session_id"})
		return
	}

	customerName := strings.TrimSpace(req.CustomerName)
	if customerName == "" {
		customerName = "未命名客户"
	}
	fontName := strings.TrimSpace(req.FontName)
	if fontName == "" {
		fontName = "微软雅黑"
	}

	state.SetSessionCustomer(req.SessionID, customerName)

	sessionUploadDir := filepath.Join(uploadDir, req.SessionID)
	entries, err := os.ReadDir(sessionUploadDir)
	if err != nil || len(entries) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件，请先上传"})
		return
	}

	var uploadPath string
	if len(entries) == 1 && entries[0].IsDir() {
		uploadPath = filepath.Join(sessionUploadDir, entries[0].Name())
	} else {
		uploadPath = sessionUploadDir
	}

	workDir, _ := os.Getwd()
	taskID := uuid.New().String()
	sessionOutputDir := filepath.Join(outputDir, taskID)
	os.MkdirAll(sessionOutputDir, 0755)

	cancelCh := make(chan struct{})
	task := &AnalysisTask{
		ID:           taskID,
		CustomerName: customerName,
		FileName:     filepath.Base(uploadPath),
		Status:       StatusPending,
		Progress:     "等待中...",
		Options:      req.Options,
		CreatedAt:    time.Now(),
		outputDir:    sessionOutputDir,
		cancel:       cancelCh,
	}
	state.AddTask(task)

	go func() {
		state.UpdateTaskStatus(taskID, StatusRunning, "正在执行分析...", nil, "", "")

		reportFiles, err := runAnalysis(workDir, uploadPath, sessionOutputDir, taskID, req.Options, fontName, cancelCh)
		if err != nil {
			state.UpdateTaskStatus(taskID, StatusFailed, "", nil, "", err.Error())
			return
		}

		t := state.GetTask(taskID)
		if t != nil && t.Status == StatusCancelled {
			state.CleanupTaskFiles(taskID)
			return
		}

		zipName := fmt.Sprintf("%s_MySQL巡检报告_%s.zip", customerName, time.Now().Format("20060102_150405"))
		zipPath := filepath.Join(sessionOutputDir, zipName)
		if err := packageToZip(sessionOutputDir, reportFiles, zipPath); err != nil {
			state.UpdateTaskStatus(taskID, StatusFailed, "", nil, "", "打包失败: "+err.Error())
			return
		}

		state.UpdateTaskStatus(taskID, StatusCompleted, "分析完成", reportFiles, zipName, "")
	}()

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "分析任务已提交",
	})
}

func handleCancelTask(c *gin.Context) {
	taskID := c.Param("id")
	if ok := state.CancelTask(taskID); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法取消：任务不存在或已结束"})
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		state.CleanupTaskFiles(taskID)
	}()
	c.JSON(http.StatusOK, gin.H{"message": "任务已取消，文件已清理"})
}

func handleDeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	task := state.GetTask(taskID)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}
	if task.Status == StatusRunning || task.Status == StatusPending {
		state.CancelTask(taskID)
	}
	state.CleanupTaskFiles(taskID)
	state.DeleteTask(taskID)
	c.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

func handleTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	task := state.GetTask(taskID)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}
	if sid := c.Query("session_id"); sid != "" {
		state.RefreshSession(sid)
	}
	c.JSON(http.StatusOK, TaskResponse{
		Task:  task,
		Stats: state.GetStats(),
	})
}

func handleDownload(c *gin.Context) {
	taskID := c.Param("task_id")
	fileName := c.Param("filename")

	task := state.GetTask(taskID)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if fileName == task.ZipFile || strings.HasSuffix(fileName, ".zip") {
		fileName = task.ZipFile
	}

	taskDir := filepath.Join(outputDir, taskID)
	filePath := filepath.Join(taskDir, fileName)

	cleanTaskDir := filepath.Clean(taskDir)
	cleanFilePath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanFilePath, cleanTaskDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "非法路径"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func handleDownloadZip(c *gin.Context) {
	taskID := c.Param("task_id")
	task := state.GetTask(taskID)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}
	if task.ZipFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "压缩包未生成"})
		return
	}
	filePath := filepath.Join(outputDir, taskID, task.ZipFile)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, task.ZipFile))
	c.Header("Content-Type", "application/zip")
	c.File(filePath)
}

func handleStats(c *gin.Context) {
	if sid := c.Query("session_id"); sid != "" {
		state.RefreshSession(sid)
	}
	c.JSON(http.StatusOK, state.GetStats())
}

func handleListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": state.GetAllTasks()})
}

func handleListRunningTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"tasks": state.GetRunningTasks()})
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": "1.0.0"})
}

func handleAuth(c *gin.Context) {
	var body struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "error": "参数错误"})
		return
	}
	password := strings.TrimSpace(body.Password)
	fmt.Printf("[Auth] 收到密码: '%s', 期望: '%s', 匹配: %v\n", password, adminPassword, password == adminPassword)
	if password == adminPassword {
		c.JSON(http.StatusOK, gin.H{"success": true, "token": "authenticated", "message": "认证成功"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": false, "error": "密码错误"})
}

func handleAdminFiles(c *gin.Context) {
	type FileInfo struct {
		Name     string `json:"name"`
		Size     int64  `json:"size"`
		Modified string `json:"modified"`
		Customer string `json:"customer"`
		Path     string `json:"path"`
	}
	results := make([]FileInfo, 0)

	reportsPath := "./web_data/reports"
	uploadPath := "./web_data/uploads"

	for _, basePath := range []string{reportsPath, uploadPath} {
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			fullPath := filepath.Join(basePath, entry.Name())
			var size int64
			if entry.IsDir() {
				size = dirSize(fullPath)
			} else {
				size = info.Size()
			}
			fileInfo := FileInfo{
				Name:     entry.Name(),
				Size:     size,
				Modified: info.ModTime().Format("2006-01-02 15:04:05"),
				Customer: resolveCustomerName(fullPath, entry.Name(), entry.IsDir()),
				Path:     fullPath,
			}
			results = append(results, fileInfo)
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": results})
}

func dirSize(dirPath string) int64 {
	var size int64
	filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})
	return size
}

func resolveCustomerName(_, dirName string, isDir bool) string {
	cust := state.GetSessionCustomer(dirName)
	if cust != "" {
		return cust
	}
	if isDir {
		task := state.GetTask(dirName)
		if task != nil {
			return task.CustomerName
		}
		return "未知"
	}
	if strings.Contains(dirName, "_MySQL巡检报告_") {
		parts := strings.SplitN(dirName, "_MySQL巡检报告_", 2)
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return "未知"
}

func handleAdminCleanup(c *gin.Context) {
	var body struct {
		Paths []string `json:"paths"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	deleted := make([]string, 0)
	errors := make([]string, 0)
	for _, p := range body.Paths {
		cleanPath := filepath.Clean(p)
		if !strings.HasPrefix(cleanPath, "web_data/") {
			errors = append(errors, "非法路径: "+p)
			continue
		}
		if err := os.RemoveAll(cleanPath); err != nil {
			errors = append(errors, "删除失败: "+p+" - "+err.Error())
		} else {
			deleted = append(deleted, p)
			dirName := filepath.Base(p)
			if state.GetTask(dirName) != nil {
				state.DeleteTask(dirName)
			}
			dbCleanupTaskByPathPrefix(dirName)
		}
	}
	c.JSON(http.StatusOK, gin.H{"deleted": deleted, "errors": errors})
}

func handleUploadArchive(c *gin.Context) {
	sessionID := c.PostForm("session_id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	state.RegisterSession(sessionID)

	file, header, err := c.Request.FormFile("archive")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取压缩文件失败: " + err.Error()})
		return
	}
	defer file.Close()

	if header.Size > maxArchiveSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件 %s 超过大小限制(30MB)", header.Filename)})
		return
	}

	sessionUploadDir := filepath.Join(uploadDir, sessionID)
	os.MkdirAll(sessionUploadDir, 0755)

	archivePath := filepath.Join(sessionUploadDir, header.Filename)
	out, err := os.Create(archivePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
		return
	}
	io.Copy(out, file)
	out.Close()

	extractDir := sessionUploadDir

	name := strings.ToLower(header.Filename)
	switch {
	case strings.HasSuffix(name, ".zip"):
		err = unzip(archivePath, extractDir)
	case strings.HasSuffix(name, ".tar.gz") || strings.HasSuffix(name, ".tgz"):
		err = untargz(archivePath, extractDir)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的压缩格式，仅支持 .zip / .tar.gz / .tgz"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解压失败: " + err.Error()})
		return
	}

	os.Remove(archivePath)

	entries, _ := os.ReadDir(extractDir)
	c.JSON(http.StatusOK, UploadResponse{
		SessionID: sessionID,
		Files:     []string{header.Filename},
		IsDir:     len(entries) > 0,
	})
}

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dst, f.Name)
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(dst)) {
			continue
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), 0755)
		rc, err := f.Open()
		if err != nil {
			continue
		}
		out, err := os.Create(fpath)
		if err != nil {
			rc.Close()
			continue
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func untargz(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		fpath := filepath.Join(dst, header.Name)
		if !strings.HasPrefix(filepath.Clean(fpath), filepath.Clean(dst)) {
			continue
		}
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(fpath, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(fpath), 0755)
			out, err := os.Create(fpath)
			if err != nil {
				continue
			}
			io.Copy(out, tr)
			out.Close()
		}
	}
	return nil
}
