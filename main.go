package main

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	routermux "github.com/gorilla/mux"
	. "gomicro/Service"
	"gomicro/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//创建user和endpiont
	user := UserService{}
	endp := GenUserEndpoint(user)

	// 创建httpHandler
	serverHandler := httptransport.NewServer(endp, DecodeUserRequest,
		EncodeUserResponse)

	//创建一个路由
	r := routermux.NewRouter()
	//将handler绑定路由
	//r.Handle(`/user/{uid:\d+}`, serverHandler)
	{
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		})
	}

	//启动时需要注册服务,退出时需要反注册服务
	errChan := make(chan error)
	//注册服务
	go (func() {
		utils.RegService()
		// 创建http服务
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	})()
	//停止服务
	go (func() {
		sig_c := make(chan os.Signal)
		//监听关闭信号，stop和ctrl+c
		signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig_c)
	})()
	////注册服务
	//utils.RegService()
	//// 创建http服务
	//http.ListenAndServe(":8080", r)
	getErr := <-errChan
	utils.UnregService()
	log.Println(getErr)

}
