package Service

import "errors"

type IuserService interface {
	GetName(userid int) string
	DelName(userid int) error
}

type UserService struct {
}

func (this UserService) GetName(userid int) string {
	if userid == 101 {
		return "wangchao"
	}
	return "guest"
}

func (this UserService) DelName(userid int) error {
	if userid == 101 {
		return errors.New("无权限删除")
	}
	return nil
}
