package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
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

func handleCreatePersonDay(c *gin.Context) {
	var body struct {
		CustomerName string `json:"customer_name"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		WorkContent  string `json:"work_content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	customerName := strings.TrimSpace(body.CustomerName)
	if customerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "客户名称不能为空"})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04", body.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误，请使用 YYYY-MM-DD HH:mm"})
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04", body.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误，请使用 YYYY-MM-DD HH:mm"})
		return
	}
	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间必须晚于开始时间"})
		return
	}

	holidays, _ := dbGetHolidays()
	record := CalculateAndBuildRecord(customerName, strings.TrimSpace(body.WorkContent), startTime, endTime, holidays)

	id, err := dbInsertPersonDay(record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
		return
	}
	record.ID = int(id)

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func handleListPersonDays(c *gin.Context) {
	records, err := dbGetPersonDays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}
	if records == nil {
		records = []PersonDayRecord{}
	}
	c.JSON(http.StatusOK, gin.H{"records": records})
}

func handleUpdatePersonDay(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var body struct {
		CustomerName string `json:"customer_name"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		WorkContent  string `json:"work_content"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	customerName := strings.TrimSpace(body.CustomerName)
	if customerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "客户名称不能为空"})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04", body.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04", body.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
		return
	}
	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间必须晚于开始时间"})
		return
	}

	holidays, _ := dbGetHolidays()
	record := CalculateAndBuildRecord(customerName, strings.TrimSpace(body.WorkContent), startTime, endTime, holidays)
	record.ID = id

	if err := dbUpdatePersonDay(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"record": record})
}

func handleDeletePersonDay(c *gin.Context) {
	idStr := c.Param("id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	if err := dbDeletePersonDay(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func handleCalculatePersonDay(c *gin.Context) {
	var body struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04", body.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04", body.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
		return
	}
	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间必须晚于开始时间"})
		return
	}

	holidays, _ := dbGetHolidays()
	calc := CalculatePersonDays(startTime, endTime, holidays)

	c.JSON(http.StatusOK, gin.H{
		"person_days":    calc.PersonDays,
		"work_hours":     calc.WorkHours,
		"overtime_hours": calc.OvertimeHours,
		"holiday_hours":  calc.HolidayHours,
		"detail":         calc.Detail,
	})
}

func handleListHolidays(c *gin.Context) {
	holidays, err := dbGetAllHolidaysList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	if holidays == nil {
		holidays = []map[string]string{}
	}
	c.JSON(http.StatusOK, gin.H{"holidays": holidays})
}

func handleAddHoliday(c *gin.Context) {
	var body struct {
		Date string `json:"date"`
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if body.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期不能为空"})
		return
	}
	if _, err := time.Parse("2006-01-02", body.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误，请使用 YYYY-MM-DD"})
		return
	}
	if err := dbInsertHoliday(body.Date, body.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "保存成功"})
}

func handleDeleteHoliday(c *gin.Context) {
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期不能为空"})
		return
	}
	if err := dbDeleteHoliday(date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func handleFetchHolidays(c *gin.Context) {
	var body struct {
		Year int `json:"year"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Year == 0 {
		body.Year = time.Now().Year()
	}

	years := []int{body.Year, body.Year + 1}
	totalImported := 0
	totalFound := 0
	insertErrors := 0
	var errors []string
	var debug []string

	transport := &http.Transport{}
	client := &http.Client{Transport: transport, Timeout: 15 * time.Second}

	for _, year := range years {
		urls := []string{
			fmt.Sprintf("https://timor.tech/api/holiday/year/%d", year),
			fmt.Sprintf("http://timor.tech/api/holiday/year/%d", year),
		}
		var bodyBytes []byte
		var successURL string
		var statusCode int

		for _, url := range urls {
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				continue
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			req.Header.Set("Accept", "application/json, text/plain, */*")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("[HolidayFetch] %s 请求失败: %v\n", url, err)
				continue
			}
			bodyBytes, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			statusCode = resp.StatusCode
			if err != nil {
				fmt.Printf("[HolidayFetch] %s 读取失败: %v\n", url, err)
				continue
			}
			if statusCode == 200 && len(bodyBytes) > 0 && bodyBytes[0] == '{' {
				successURL = url
				break
			}
			bodyPreview := string(bodyBytes[:min(len(bodyBytes), 500)])
			fmt.Printf("[HolidayFetch] %s status=%d body=%s\n", url, statusCode, bodyPreview)
		}

		if successURL == "" {
			msg := fmt.Sprintf("API请求失败(year=%d), 所有URL均未返回有效JSON", year)
			fmt.Printf("[HolidayFetch] %s\n", msg)
			errors = append(errors, msg)
			continue
		}

		bodyStr := string(bodyBytes)
		debug = append(debug, fmt.Sprintf("年份%d(%s): status=%d len=%d", year, successURL, statusCode, len(bodyStr)))
		fmt.Printf("[HolidayFetch] year=%d status=%d len=%d\n", year, statusCode, len(bodyStr))

		var result struct {
			Code    int                     `json:"code"`
			Holiday map[string]*json.RawMessage `json:"holiday"`
		}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			msg := fmt.Sprintf("JSON解析失败(year=%d): %v, body前100字=%s", year, err, bodyStr[:min(len(bodyStr), 100)])
			fmt.Printf("[HolidayFetch] %s\n", msg)
			errors = append(errors, msg)
			continue
		}
		if result.Code != 0 {
			msg := fmt.Sprintf("API错误码(year=%d): code=%d", year, result.Code)
			fmt.Printf("[HolidayFetch] %s\n", msg)
			errors = append(errors, msg)
			continue
		}

		debug = append(debug, fmt.Sprintf("年份%d: API返回%d条记录", year, len(result.Holiday)))

		for dateKey, raw := range result.Holiday {
			var item struct {
				Holiday bool   `json:"holiday"`
				Name    string `json:"name"`
				Date    string `json:"date"`
			}
			if err := json.Unmarshal(*raw, &item); err != nil {
				continue
			}
			if item.Holiday {
				totalFound++
				dateStr := item.Date
				if dateStr == "" {
					dateStr = fmt.Sprintf("%d-%s", year, dateKey)
				}
				if err := dbInsertHoliday(dateStr, item.Name); err != nil {
					insertErrors++
					msg := fmt.Sprintf("插入失败: date=%s name=%s err=%v", dateStr, item.Name, err)
					fmt.Printf("[HolidayFetch] %s\n", msg)
					if insertErrors <= 3 {
						errors = append(errors, msg)
					}
				} else {
					totalImported++
				}
			}
		}
	}

	debug = append(debug, fmt.Sprintf("共找到%d个节假日, 成功导入%d, 插入失败%d", totalFound, totalImported, insertErrors))

	holidays, _ := dbGetAllHolidaysList()
	if holidays == nil {
		holidays = []map[string]string{}
	}

	msg := fmt.Sprintf("已同步 %d 个节假日", totalImported)
	if len(errors) > 0 {
		msg += "\n错误: " + strings.Join(errors, "; ")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  msg,
		"holidays": holidays,
		"errors":   errors,
		"debug":    debug,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
