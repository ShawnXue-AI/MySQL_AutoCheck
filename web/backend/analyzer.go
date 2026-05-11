package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var pythonCmd string

func init() {
	pythonCmd = detectPython()
}

func detectPython() string {
	candidates := []string{"python3", "python", "py"}
	if runtime.GOOS == "windows" {
		candidates = []string{"python", "py", "python3"}
	}
	for _, name := range candidates {
		_, err := exec.LookPath(name)
		if err == nil {
			return name
		}
	}
	return "python"
}

func runAnalysis(workDir, uploadPath, outputDir, taskID string, opts AnalyzeOptions, fontName string, cancel <-chan struct{}) ([]string, error) {
	uploadPath = filepath.Clean(uploadPath)
	info, err := os.Stat(uploadPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问上传路径: %w", err)
	}

	select {
	case <-cancel:
		return nil, nil
	default:
	}

	if info.IsDir() {
		return runBatchAnalysis(workDir, uploadPath, outputDir, opts, fontName, cancel)
	}
	return runSingleAnalysis(workDir, uploadPath, outputDir, opts, fontName, cancel)
}

func runSingleAnalysis(workDir, logFile, outputDir string, opts AnalyzeOptions, fontName string, cancel <-chan struct{}) ([]string, error) {
	args := []string{"-u", "log_analyzer.py", logFile}
	if fontName != "" && fontName != "微软雅黑" {
		args = []string{"-u", "log_analyzer.py", "--font", fontName, logFile}
	}
	cmd := exec.Command(pythonCmd, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	errCh := make(chan error, 1)
	go func() {
		_, err := cmd.CombinedOutput()
		errCh <- err
	}()

	select {
	case <-cancel:
		cmd.Process.Kill()
		return nil, nil
	case err := <-errCh:
		if err != nil {
			output, _ := cmd.CombinedOutput()
			return nil, fmt.Errorf("单实例分析失败: %w\n输出: %s", err, string(output))
		}
	}

	reportsDir := filepath.Join(workDir, "reports")
	var files []string
	if entries, err := os.ReadDir(reportsDir); err == nil {
		os.MkdirAll(outputDir, 0755)
		for _, e := range entries {
			if !e.IsDir() {
				src := filepath.Join(reportsDir, e.Name())
				dst := filepath.Join(outputDir, e.Name())
				copyFile(src, dst)
				if opts.GenerateWord || !strings.HasSuffix(e.Name(), ".docx") {
					files = append(files, e.Name())
				}
			}
		}
	}
	return files, nil
}

func runBatchAnalysis(workDir, logDir, outputDir string, opts AnalyzeOptions, fontName string, cancel <-chan struct{}) ([]string, error) {
	args := []string{"-u", "batch_analyze.py", logDir}
	if fontName != "" && fontName != "微软雅黑" {
		args = append(args, "--font", fontName)
	}
	if opts.GenerateExcel {
		args = append(args, "--excel")
	}
	if !opts.GenerateWord {
		args = append(args, "-f", "txt")
	}

	cmd := exec.Command(pythonCmd, args...)
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	errCh := make(chan error, 1)
	go func() {
		_, err := cmd.CombinedOutput()
		errCh <- err
	}()

	select {
	case <-cancel:
		cmd.Process.Kill()
		return nil, nil
	case err := <-errCh:
		if err != nil {
			output, _ := cmd.CombinedOutput()
			return nil, fmt.Errorf("批量分析失败: %w\n输出: %s", err, string(output))
		}
	}

	patterns, _ := filepath.Glob(filepath.Join(workDir, "batch_reports_*"))
	var latestDir string
	if len(patterns) > 0 {
		latestDir = patterns[len(patterns)-1]
	}

	var files []string
	if latestDir != "" {
		os.MkdirAll(outputDir, 0755)
		entries, _ := os.ReadDir(latestDir)
		for _, e := range entries {
			if !e.IsDir() {
				src := filepath.Join(latestDir, e.Name())
				dst := filepath.Join(outputDir, e.Name())
				copyFile(src, dst)
				if opts.GenerateWord || !strings.HasSuffix(e.Name(), ".docx") {
					files = append(files, e.Name())
				}
			}
		}
	}
	return files, nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func packageToZip(srcDir string, files []string, zipPath string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	for _, name := range files {
		src := filepath.Join(srcDir, name)
		info, err := os.Stat(src)
		if err != nil {
			continue
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			continue
		}
		hdr.Name = name
		hdr.Method = zip.Deflate
		writer, err := w.CreateHeader(hdr)
		if err != nil {
			continue
		}
		r, err := os.Open(src)
		if err != nil {
			continue
		}
		io.Copy(writer, r)
		r.Close()
	}
	return nil
}
