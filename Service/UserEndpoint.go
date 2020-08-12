package Service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct {
	Uid    int `json:"uid"`
	Method string
}
type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndpoint(userService IuserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		result := "noting"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid)
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
