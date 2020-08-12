package utils

import (
	consulapi "github.com/hashicorp/consul/api"
	"log"
)

var ConsulClient *consulapi.Client

func init() {
	config := consulapi.DefaultConfig()
	//consul服务器信息
	config.Address = "119.3.230.228:8500"

	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client

}

func RegService() {

	//被注册服务信息
	res := consulapi.AgentServiceRegistration{}
	res.ID = "userservice"
	res.Name = "userservice"
	res.Address = "169.254.175.50"
	res.Port = 8080
	res.Tags = []string{"primary"}

	//心跳Api
	check := consulapi.AgentServiceCheck{}
	check.Interval = "5s"
	check.HTTP = "http://169.254.175.50:8080/health"

	//将心跳放入注册服务信息
	res.Check = &check

	//将客户端注册
	err := ConsulClient.Agent().ServiceRegister(&res)
	if err != nil {
		log.Fatal(err)
	}
}

//反注册
func UnregService() {
	ConsulClient.Agent().ServiceDeregister("userservice")
}
