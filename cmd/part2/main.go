package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

/**
作业要求
1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200
*/

func main() {

	http.HandleFunc("/header", headerFunc)

	http.HandleFunc("/version", versionFunc)

	http.HandleFunc("/log", logFunc)

	http.HandleFunc("/healthz", healthzFunc)

	log.Fatal(http.ListenAndServe(":80", nil))

}

func headerFunc(w http.ResponseWriter, r *http.Request) {

	headers := r.Header

	for k, v := range headers {
		fmt.Println(k, "----->", v)
		r := strings.Join(v, ";")
		w.Header().Set(k, r)
	}
}

// VERSION=V1  go run cmd/part2/main.go
func versionFunc(w http.ResponseWriter, r *http.Request) {

	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)

}

func logFunc(w http.ResponseWriter, r *http.Request) {

	ip, _ := getIP(r)
	statusCode := 200 // 服务器端默认状态码为 200, 根据业务规则来, 比如 4xx, 5xx等等
	log.Printf("client ip : %s , http code : %d", ip, statusCode)

}

// https://www.cnblogs.com/GaiHeiluKamei/p/13731791.html
// 暂时不纠结, 先不做 严谨测试 ，只需要获取到一个 IP
func getIP(r *http.Request) (string, error) {

	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")

}

func healthzFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
