package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Tab struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Config struct {
	SnapshotDir string `json:"snapshot_dir"`
	Port        int    `json:"port"`
}

var config Config

func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		panic("cannot open config.json")
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&config)
}

func getSnapshotFilePath() string {

	now := time.Now()

	// snapshots/2026-05-14
	dateDir := filepath.Join(
		config.SnapshotDir,
		now.Format("2006-01-02"),
	)

	// 自动创建目录
	os.MkdirAll(dateDir, os.ModePerm)

	// 14-32-01.json
	baseName := now.Format("15-04-05")

	filePath := filepath.Join(
		dateDir,
		baseName+".json",
	)

	// 重名处理
	index := 1
	for {
		_, err := os.Stat(filePath)

		if os.IsNotExist(err) {
			break
		}

		filePath = filepath.Join(
			dateDir,
			fmt.Sprintf(
				"%s_%d.json",
				baseName,
				index,
			),
		)

		index++
	}

	return filePath
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	var tabs []Tab

	err := json.NewDecoder(r.Body).Decode(&tabs)
	if err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	filePath := getSnapshotFilePath()

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "cannot create file", 500)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(tabs)
	if err != nil {
		http.Error(w, "cannot write file", 500)
		return
	}

	fmt.Println("saved:", filePath)

	w.Write([]byte("ok"))
}

func main() {

	loadConfig()

	http.HandleFunc("/save", saveHandler)

	addr := fmt.Sprintf(":%d", config.Port)

	fmt.Println("Server running at", addr)

	http.ListenAndServe(addr, nil)
}
