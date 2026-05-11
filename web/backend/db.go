package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type dbTask struct {
	ID           string         `json:"id"`
	CustomerName sql.NullString `json:"customer_name"`
	FileName     sql.NullString `json:"file_name"`
	Status       sql.NullString `json:"status"`
	Options      sql.NullString `json:"options"`
	ReportFiles  sql.NullString `json:"report_files"`
	ZipFile      sql.NullString `json:"zip_file"`
	Error        sql.NullString `json:"error"`
	CreatedAt    time.Time      `json:"created_at"`
	CompletedAt  sql.NullTime   `json:"completed_at"`
}

func initDB() {
	dsn := "inspection_user:admin888@tcp(127.0.0.1:3306)/inspection_db?charset=utf8mb4&parseTime=true&loc=Local"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)
	if err = db.Ping(); err != nil {
		log.Fatalf("MySQL Ping 失败: %v", err)
	}
	log.Println("MySQL 连接成功: inspection_db")

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) NOT NULL COMMENT '任务ID',
			customer_name VARCHAR(100) DEFAULT NULL COMMENT '客户名称',
			file_name VARCHAR(255) DEFAULT NULL COMMENT '文件名',
			status VARCHAR(20) DEFAULT NULL COMMENT '任务状态',
			options JSON DEFAULT NULL COMMENT '任务选项',
			report_files JSON DEFAULT NULL COMMENT '报告文件列表',
			zip_file VARCHAR(500) DEFAULT NULL COMMENT '打包文件路径',
			error TEXT DEFAULT NULL COMMENT '错误信息',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			completed_at DATETIME DEFAULT NULL COMMENT '完成时间',
			PRIMARY KEY (id),
			KEY idx_customer_name (customer_name),
			KEY idx_status (status),
			KEY idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		log.Fatalf("创建 tasks 表失败: %v", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			session_id VARCHAR(128) NOT NULL COMMENT '会话ID',
			customer_name VARCHAR(100) DEFAULT NULL COMMENT '客户名称',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			PRIMARY KEY (session_id),
			KEY idx_customer_name (customer_name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		log.Fatalf("创建 sessions 表失败: %v", err)
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS person_days (
			id INT AUTO_INCREMENT PRIMARY KEY,
			customer_name VARCHAR(100) NOT NULL COMMENT '客户名称',
			start_time DATETIME NOT NULL COMMENT '开始时间',
			end_time DATETIME NOT NULL COMMENT '结束时间',
			person_days DECIMAL(10,2) NOT NULL COMMENT '人天数',
			work_hours DECIMAL(10,2) NOT NULL COMMENT '工作时间(小时)',
			overtime_hours DECIMAL(10,2) NOT NULL COMMENT '加班时间(小时)',
			holiday_hours DECIMAL(10,2) NOT NULL COMMENT '节假日时间(小时)',
			work_content TEXT DEFAULT NULL COMMENT '工作内容',
			detail TEXT DEFAULT NULL COMMENT '计算明细',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			KEY idx_customer_name (customer_name),
			KEY idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		log.Fatalf("创建 person_days 表失败: %v", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS holidays (
			date DATE NOT NULL PRIMARY KEY COMMENT '节假日日期',
			name VARCHAR(50) DEFAULT NULL COMMENT '节假日名称',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`); err != nil {
		log.Fatalf("创建 holidays 表失败: %v", err)
	}

	log.Println("数据库表初始化完成")
}

func dbInsertTask(task *AnalysisTask) {
	optsJSON, _ := json.Marshal(task.Options)
	reportsJSON, _ := json.Marshal(task.ReportFiles)
	var completedAt interface{}
	if task.CompletedAt != nil {
		completedAt = task.CompletedAt
	}
	_, err := db.Exec(
		`INSERT INTO tasks (id, customer_name, file_name, status, options, report_files, zip_file, error, created_at, completed_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE status=VALUES(status), options=VALUES(options), report_files=VALUES(report_files),
		 zip_file=VALUES(zip_file), error=VALUES(error), completed_at=VALUES(completed_at)`,
		task.ID, task.CustomerName, task.FileName, string(task.Status),
		string(optsJSON), string(reportsJSON), task.ZipFile, task.Error,
		task.CreatedAt, completedAt,
	)
	if err != nil {
		log.Printf("[DB] 写入任务失败 (%s): %v", task.ID, err)
	}
}

func dbGetAllTasks() []*AnalysisTask {
	rows, err := db.Query(`SELECT id, customer_name, file_name, status, options, report_files, zip_file, error, created_at, completed_at FROM tasks ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var results []*AnalysisTask
	for rows.Next() {
		var dt dbTask
		if err := rows.Scan(&dt.ID, &dt.CustomerName, &dt.FileName, &dt.Status, &dt.Options, &dt.ReportFiles, &dt.ZipFile, &dt.Error, &dt.CreatedAt, &dt.CompletedAt); err != nil {
			continue
		}
		task := dbToTask(&dt)
		results = append(results, task)
	}
	return results
}

func dbGetTaskCountByStatus(status ...string) int64 {
	if db == nil {
		return 0
	}
	query := "SELECT COUNT(*) FROM tasks WHERE status = ?"
	args := make([]interface{}, 0, len(status))
	for i, s := range status {
		if i == 0 {
			args = append(args, s)
		} else {
			query += " OR status = ?"
			args = append(args, s)
		}
	}
	var count int64
	db.QueryRow(query, args...).Scan(&count)
	return count
}

func dbDeleteTask(id string) {
	db.Exec("DELETE FROM tasks WHERE id = ?", id)
}

func dbInsertSession(sessionID, customerName string) {
	db.Exec(
		"INSERT INTO sessions (session_id, customer_name) VALUES (?, ?) ON DUPLICATE KEY UPDATE customer_name=VALUES(customer_name)",
		sessionID, customerName,
	)
}

func dbGetSessionCustomer(sessionID string) string {
	var name sql.NullString
	db.QueryRow("SELECT customer_name FROM sessions WHERE session_id = ?", sessionID).Scan(&name)
	if name.Valid {
		return name.String
	}
	return ""
}

func dbDeleteSession(sessionID string) {
	db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
}

func dbToTask(dt *dbTask) *AnalysisTask {
	var opts AnalyzeOptions
	if dt.Options.Valid && dt.Options.String != "" {
		json.Unmarshal([]byte(dt.Options.String), &opts)
	}
	var reports []string
	if dt.ReportFiles.Valid && dt.ReportFiles.String != "" {
		json.Unmarshal([]byte(dt.ReportFiles.String), &reports)
	}
	task := &AnalysisTask{
		ID:           dt.ID,
		CustomerName: nullableString(dt.CustomerName),
		FileName:     nullableString(dt.FileName),
		Status:       TaskStatus(nullableString(dt.Status)),
		Options:      opts,
		CreatedAt:    dt.CreatedAt,
		ReportFiles:  reports,
		ZipFile:      nullableString(dt.ZipFile),
		Error:        nullableString(dt.Error),
	}
	if dt.CompletedAt.Valid {
		task.CompletedAt = &dt.CompletedAt.Time
	}
	if reports == nil {
		task.ReportFiles = []string{}
	}
	return task
}

func nullableString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func dbGetTotalCompleted() int64 {
	var count int64
	db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = ?", string(StatusCompleted)).Scan(&count)
	return count
}

func dbGetTotalTasks() int64 {
	var count int64
	db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	return count
}

func dbGetAllTasksJSON() []*AnalysisTask {
	return dbGetAllTasks()
}

func dbGetRunningTasks() []*AnalysisTask {
	rows, err := db.Query(`SELECT id, customer_name, file_name, status, options, report_files, zip_file, error, created_at, completed_at FROM tasks WHERE status = ? OR status = ? ORDER BY created_at DESC`, string(StatusRunning), string(StatusPending))
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) []*AnalysisTask {
	var results []*AnalysisTask
	for rows.Next() {
		var dt dbTask
		if err := rows.Scan(&dt.ID, &dt.CustomerName, &dt.FileName, &dt.Status, &dt.Options, &dt.ReportFiles, &dt.ZipFile, &dt.Error, &dt.CreatedAt, &dt.CompletedAt); err != nil {
			continue
		}
		task := dbToTask(&dt)
		results = append(results, task)
	}
	return results
}

func dbGetAllCustomerSessions() map[string]string {
	result := make(map[string]string)
	rows, err := db.Query("SELECT session_id, customer_name FROM sessions")
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var sid, name sql.NullString
		if err := rows.Scan(&sid, &name); err == nil && sid.Valid {
			result[sid.String] = nullableString(name)
		}
	}
	return result
}

func dbLoadTasksToMemory(s *AppState) {
	tasks := dbGetAllTasks()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tasks {
		s.tasks[t.ID] = t
		if t.Status == StatusCompleted {
			s.totalCompleted++
		}
	}
	s.totalCompleted = dbGetTotalCompleted()
	sessions := dbGetAllCustomerSessions()
	for sid, name := range sessions {
		s.sessionCustomers[sid] = name
	}
	loaded := len(tasks)
	completed := s.totalCompleted
	var pending, running int
	for _, t := range tasks {
		switch t.Status {
		case StatusPending:
			pending++
		case StatusRunning:
			running++
		}
	}
	log.Printf("[DB] 从 MySQL 加载了 %d 条任务 (已完成: %d, 运行中: %d, 等待中: %d), %d 个会话映射",
		loaded, completed, running, pending, len(sessions))
}

func dbCleanupTaskByPathPrefix(dirName string) {
	db.Exec("DELETE FROM tasks WHERE id = ?", dirName)
	db.Exec("DELETE FROM sessions WHERE session_id = ?", dirName)
}

func dbCountTaskResults() map[string]int {
	result := map[string]int{"completed": 0, "failed": 0, "running": 0, "pending": 0, "cancelled": 0}
	rows, err := db.Query("SELECT status, COUNT(*) as cnt FROM tasks GROUP BY status")
	if err != nil {
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		var c int
		if err := rows.Scan(&s, &c); err == nil {
			result[s] = c
		}
	}
	return result
}

func dbVerifyConnectivity() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.Ping()
}

type PersonDayRecord struct {
	ID            int       `json:"id"`
	CustomerName  string    `json:"customer_name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	PersonDays    float64   `json:"person_days"`
	WorkHours     float64   `json:"work_hours"`
	OvertimeHours float64   `json:"overtime_hours"`
	HolidayHours  float64   `json:"holiday_hours"`
	WorkContent   string    `json:"work_content"`
	Detail        string    `json:"detail"`
	CreatedAt     time.Time `json:"created_at"`
}

func dbInsertPersonDay(r *PersonDayRecord) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO person_days (customer_name, start_time, end_time, person_days, work_hours, overtime_hours, holiday_hours, work_content, detail)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.CustomerName, r.StartTime, r.EndTime, r.PersonDays, r.WorkHours, r.OvertimeHours, r.HolidayHours, r.WorkContent, r.Detail,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func dbGetPersonDays() ([]PersonDayRecord, error) {
	rows, err := db.Query(`SELECT id, customer_name, start_time, end_time, person_days, work_hours, overtime_hours, holiday_hours, work_content, detail, created_at FROM person_days ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []PersonDayRecord
	for rows.Next() {
		var r PersonDayRecord
		var workContent, detail sql.NullString
		if err := rows.Scan(&r.ID, &r.CustomerName, &r.StartTime, &r.EndTime, &r.PersonDays, &r.WorkHours, &r.OvertimeHours, &r.HolidayHours, &workContent, &detail, &r.CreatedAt); err != nil {
			continue
		}
		if workContent.Valid {
			r.WorkContent = workContent.String
		}
		if detail.Valid {
			r.Detail = detail.String
		}
		results = append(results, r)
	}
	return results, nil
}

func dbUpdatePersonDay(r *PersonDayRecord) error {
	_, err := db.Exec(
		`UPDATE person_days SET customer_name=?, start_time=?, end_time=?, person_days=?, work_hours=?, overtime_hours=?, holiday_hours=?, work_content=?, detail=? WHERE id=?`,
		r.CustomerName, r.StartTime, r.EndTime, r.PersonDays, r.WorkHours, r.OvertimeHours, r.HolidayHours, r.WorkContent, r.Detail, r.ID,
	)
	return err
}

func dbDeletePersonDay(id int) error {
	_, err := db.Exec(`DELETE FROM person_days WHERE id=?`, id)
	return err
}

func dbGetHolidays() (map[string]string, error) {
	rows, err := db.Query(`SELECT date, name FROM holidays ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var date time.Time
		var name sql.NullString
		if err := rows.Scan(&date, &name); err == nil {
			result[date.Format("2006-01-02")] = nullableString(name)
		}
	}
	return result, nil
}

func dbInsertHoliday(date string, name string) error {
	_, err := db.Exec(`INSERT INTO holidays (date, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name=VALUES(name)`, date, name)
	return err
}

func dbDeleteHoliday(date string) error {
	_, err := db.Exec(`DELETE FROM holidays WHERE date=?`, date)
	return err
}

func dbGetAllHolidaysList() ([]map[string]string, error) {
	rows, err := db.Query(`SELECT date, name FROM holidays ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []map[string]string
	for rows.Next() {
		var date time.Time
		var name sql.NullString
		if err := rows.Scan(&date, &name); err == nil {
			results = append(results, map[string]string{
				"date": date.Format("2006-01-02"),
				"name": nullableString(name),
			})
		}
	}
	return results, nil
}
