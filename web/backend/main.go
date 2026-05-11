package main

import (
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	workDir, err := filepath.Abs("../..")
	if err != nil {
		panic("获取项目根目录失败: " + err.Error())
	}
	if err := os.Chdir(workDir); err != nil {
		panic("切换工作目录失败: " + err.Error())
	}

	initDB()
	dbLoadTasksToMemory(state)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/api/health", handleHealth)
	r.GET("/api/stats", handleStats)
	r.POST("/api/upload", handleUpload)
	r.POST("/api/upload/archive", handleUploadArchive)
	r.POST("/api/analyze", handleStartAnalysis)
	r.GET("/api/tasks", handleListTasks)
	r.GET("/api/tasks/running", handleListRunningTasks)
	r.GET("/api/tasks/:id", handleTaskStatus)
	r.POST("/api/tasks/:id/cancel", handleCancelTask)
	r.DELETE("/api/tasks/:id", handleDeleteTask)
	r.GET("/api/download/:task_id/:filename", handleDownload)
	r.GET("/api/download/zip/:task_id", handleDownloadZip)
	r.POST("/api/admin/auth", handleAuth)
	r.GET("/api/admin/files", handleAdminFiles)
	r.POST("/api/admin/cleanup", handleAdminCleanup)

	r.Static("/api/reports", "./web_data/reports")

	r.Run(":8080")
}
