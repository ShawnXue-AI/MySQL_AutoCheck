package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

type AnalyzeOptions struct {
	GenerateWord  bool `json:"generate_word"`
	GenerateExcel bool `json:"generate_excel"`
}

type AnalysisTask struct {
	ID           string         `json:"id"`
	CustomerName string         `json:"customer_name"`
	FileName     string         `json:"file_name"`
	Status       TaskStatus     `json:"status"`
	Progress     string         `json:"progress"`
	Options      AnalyzeOptions `json:"options"`
	CreatedAt    time.Time      `json:"created_at"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`
	ReportFiles  []string       `json:"report_files"`
	ZipFile      string         `json:"zip_file"`
	Error        string         `json:"error,omitempty"`
	outputDir    string
	cancel       chan struct{}
	cmd          *exec.Cmd
}

type Stats struct {
	ActiveUsers    int   `json:"active_users"`
	TotalCompleted int64 `json:"total_completed"`
	RunningTasks   int   `json:"running_tasks"`
}

type AppState struct {
	mu               sync.RWMutex
	tasks            map[string]*AnalysisTask
	activeSessions   map[string]time.Time
	sessionCustomers map[string]string
	totalCompleted   int64
}

func NewAppState() *AppState {
	return &AppState{
		tasks:            make(map[string]*AnalysisTask),
		activeSessions:   make(map[string]time.Time),
		sessionCustomers: make(map[string]string),
	}
}

func (s *AppState) AddTask(task *AnalysisTask) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.ID] = task
	go dbInsertTask(task)
}

func (s *AppState) GetTask(id string) *AnalysisTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks[id]
}

func (s *AppState) UpdateTaskStatus(id string, status TaskStatus, progress string, reportFiles []string, zipFile string, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.tasks[id]; ok {
		task.Status = status
		task.Progress = progress
		if reportFiles != nil {
			task.ReportFiles = reportFiles
		}
		if zipFile != "" {
			task.ZipFile = zipFile
		}
		if errMsg != "" {
			task.Error = errMsg
		}
		if status == StatusCompleted {
			now := time.Now()
			task.CompletedAt = &now
			s.totalCompleted++
		}
		if status == StatusCancelled {
			now := time.Now()
			task.CompletedAt = &now
		}
		go dbInsertTask(task)
	}
}

func (s *AppState) CancelTask(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return false
	}
	if task.Status != StatusRunning && task.Status != StatusPending {
		return false
	}
	if task.cancel != nil {
		close(task.cancel)
	}
	if task.cmd != nil && task.cmd.Process != nil {
		task.cmd.Process.Kill()
	}
	task.Status = StatusCancelled
	task.Progress = "已取消"
	task.Error = "用户手动终止"
	go dbInsertTask(task)
	return true
}

func (s *AppState) SetTaskCmd(id string, cmd *exec.Cmd, cancel chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.tasks[id]; ok {
		task.cmd = cmd
		task.cancel = cancel
		task.outputDir = filepath.Join(outputDir, id)
	}
}

func (s *AppState) CleanupTaskFiles(id string) {
	taskDir := filepath.Join(outputDir, id)
	os.RemoveAll(taskDir)
	uploadTaskDir := filepath.Join(uploadDir, id)
	os.RemoveAll(uploadTaskDir)
}

func (s *AppState) DeleteTask(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		go dbDeleteTask(id)
		return true
	}
	return false
}

func (s *AppState) GetRunningTasks() []*AnalysisTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*AnalysisTask, 0)
	for _, task := range s.tasks {
		if task.Status == StatusRunning || task.Status == StatusPending {
			result = append(result, task)
		}
	}
	return result
}

func (s *AppState) GetAllTasks() []*AnalysisTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*AnalysisTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *AppState) RegisterSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSessions[sessionID] = time.Now()
}

func (s *AppState) UnregisterSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.activeSessions, sessionID)
}

func (s *AppState) GetStats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.clearStaleSessions()
	running := 0
	for _, task := range s.tasks {
		if task.Status == StatusRunning || task.Status == StatusPending {
			running++
		}
	}
	totalCompleted := s.totalCompleted
	if totalCompleted == 0 {
		totalCompleted = dbGetTotalCompleted()
	}
	return Stats{
		ActiveUsers:    len(s.activeSessions),
		TotalCompleted: totalCompleted,
		RunningTasks:   running,
	}
}

func (s *AppState) clearStaleSessions() {
	threshold := time.Now().Add(-30 * time.Second)
	for id, lastSeen := range s.activeSessions {
		if lastSeen.Before(threshold) {
			delete(s.activeSessions, id)
		}
	}
}

func (s *AppState) RefreshSession(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSessions[id] = time.Now()
}

func (s *AppState) SetSessionCustomer(sessionID, customerName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionCustomers[sessionID] = customerName
	go dbInsertSession(sessionID, customerName)
}

func (s *AppState) GetSessionCustomer(sessionID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if name, ok := s.sessionCustomers[sessionID]; ok && name != "" {
		return name
	}
	name := dbGetSessionCustomer(sessionID)
	if name != "" {
		s.mu.RUnlock()
		s.mu.Lock()
		s.sessionCustomers[sessionID] = name
		s.mu.Unlock()
		s.mu.RLock()
		return name
	}
	return ""
}
