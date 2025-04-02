package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db = setupDB()
	r := mux.NewRouter()

	r.HandleFunc("/api/create", createShortURL).Methods("POST")
	r.HandleFunc("/{shortCode}", redirectURL).Methods("GET")
	r.HandleFunc("/api/stats", getStats).Methods("GET")

	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // 精确匹配前端地址
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type"},
        AllowCredentials: true,
        Debug:            true, // 开启调试日志
    })

	fmt.Println("Server is running on port 8080")
	handler := c.Handler(r)
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var data struct {
		LongURL string `json:"long_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortCode := generateRandomCode(6)
	_, err := db.Exec("INSERT INTO urls (short_code, long_url) VALUES (?, ?)", shortCode, data.LongURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"short_code": shortCode})
}
func redirectURL(w http.ResponseWriter, r *http.Request) {
    shortCode := mux.Vars(r)["shortCode"]
    
    // 打印调试信息
    log.Printf("正在处理短码: %s", shortCode)

    // 1. 查询原始URL
    var longURL string
    query := "SELECT long_url FROM urls WHERE short_code = ?"
    log.Printf("执行查询: %s, 参数: %s", query, shortCode)
    err := db.QueryRow(query, shortCode).Scan(&longURL)
    if err != nil {
        log.Printf("查询失败: %v", err) // 明确打印错误类型
        http.NotFound(w, r)
        return
    }

    // 2. 更新访问次数
    updateQuery := "UPDATE urls SET visit_count = visit_count + 1 WHERE short_code = ?"
    log.Printf("执行更新: %s, 参数: %s", updateQuery, shortCode)
    res, err := db.Exec(updateQuery, shortCode)
    if err != nil {
        log.Printf("更新失败: %v", err)
        http.Error(w, "服务器错误", http.StatusInternalServerError)
        return
    }

    // 检查受影响的行数
    rowsAffected, _ := res.RowsAffected()
    log.Printf("受影响的行数: %d", rowsAffected) // 应该是 1

    // 3. 重定向
    http.Redirect(w, r, longURL, http.StatusFound)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT short_code, visit_count,long_url FROM urls ORDER BY visit_count DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stats []ShortURL
	for rows.Next() {
		var s ShortURL
		if err := rows.Scan(&s.ShortCode, &s.VisitCount, &s.LongURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stats = append(stats, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func generateRandomCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
