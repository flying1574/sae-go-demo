package main

import (
    "fmt"
    "html/template"
    "net/http"
    "os"
    "sync"
)

var (
    hostname   string
    accessCount int
    mu sync.Mutex
)

func main() {
    // 获取主机名
    var err error
    hostname, err = os.Hostname()
    if err != nil {
        fmt.Println("Error getting hostname:", err)
        return
    }

    // 设置HTTP路由
    http.HandleFunc("/", handler)

    // 启动HTTP服务器
    fmt.Println("Server is running on :8080")
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 获取客户端IP地址
    clientIP := getClientIP(r)

    // 增加访问次数
    mu.Lock()
    accessCount++
    mu.Unlock()

    // 定义HTML模板
    htmlTemplate := `
    <html>
    <head>
        <title>Web应用</title>
    </head>
    <body>
        <h1>主机名: {{.Hostname}}</h1>
        <h1>客户端IP地址: {{.ClientIP}}</h1>
        <h1>访问次数: {{.AccessCount}}</h1>
    </body>
    </html>
    `

    // 解析HTML模板
    tmpl, err := template.New("webpage").Parse(htmlTemplate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 数据传递给模板
    data := map[string]interface{}{
        "Hostname":   hostname,
        "ClientIP":   clientIP,
        "AccessCount": accessCount,
    }

    // 渲染HTML页面并将结果写入HTTP响应
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
   }
}

func getClientIP(r *http.Request) string {
    ip := r.Header.Get("X-Real-IP")
    if ip == "" {
        ip = r.Header.Get("X-Forwarded-For")
        if ip == "" {
            ip = r.RemoteAddr
        }
    }
    return ip
}
