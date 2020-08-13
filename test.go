package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

//阻塞等待
func waitn() {
	r := rate.NewLimiter(1, 5)
	ctx := context.Background()
	for {
		err := r.WaitN(ctx, 3)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second)
	}
}

func allow() {
	r := rate.NewLimiter(1, 5)
	for {
		if r.AllowN(time.Now(), 2) {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("to many request")
		}
		time.Sleep(time.Second)
	}
}

var r = rate.NewLimiter(1, 5)

func MyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !r.Allow() {
			http.Error(writer, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func main() {
	//阻塞等待
	//waitn()
	//allow()
	//限流模式
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK!!!"))
	})
	http.ListenAndServe(":8080", MyLimit(mux))
}
