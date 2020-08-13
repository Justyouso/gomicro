package Services

import (
	"context"
	"encoding/json"
	"errors"
	routermux "github.com/gorilla/mux"
	"gomicro/utils"
	"net/http"
	"strconv"
)

// 解码
func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	//if r.URL.Query().Get("uid") != "" {
	//	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	//	return UserRequest{Uid: uid}, nil
	//}
	//通过Vars获取参数map
	vars := routermux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{Uid: uid, Method: r.Method}, nil
	}
	return nil, errors.New("参数错误")
}

// 编码
func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("content-type", contentType)
	if myerr, ok := err.(*utils.MyError); ok {
		w.WriteHeader(myerr.Code)
	} else {
		w.WriteHeader(500)
	}
	//w.WriteHeader(404)
	w.Write(body)
}
