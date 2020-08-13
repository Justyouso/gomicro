package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"gomicro/utils"
	"strconv"
)

type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
}
type UserResponse struct {
	Result string `json:"result"`
}

//限流装饰器
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (
			response interface{}, err error) {
			if !limit.Allow() {
				//return nil, errors.New("too many request")
				//自定义错误信息
				return nil, utils.NewMyError(429, "to many requests")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEndpoint(userService IuserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "noting"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelName(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("%d用户删除成功", r.Uid)
			}
		}

		return UserResponse{Result: result}, nil
	}
}
