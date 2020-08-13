package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	routermux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
	. "gomicro/Services"
	"gomicro/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// 使用命令行解析运行多个服务
	name := flag.String("name", "", "服务名称")
	port := flag.Int("p", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	// 设置服务名和端口
	utils.SetServiceNameAndPort(*name, *port)

	//创建user和endpiont
	user := UserService{}
	//endp := GenUserEndpoint(user)

	//有限流功能的endpoint
	limit := rate.NewLimiter(1, 5)
	//endpoint处理包含了错误信息
	endp := RateLimit(limit)(GenUserEndpoint(user))

	//自定义处理error
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(MyErrorEncoder),
	}
	// 创建httpHandler
	serverHandler := httptransport.NewServer(endp, DecodeUserRequest,
		EncodeUserResponse, options...)

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
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r)
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
