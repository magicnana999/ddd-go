package user

import (
	"github.com/magicnana999/ddd-go/domain/user/service"
	"github.com/magicnana999/ddd-go/dto/user"
)

func Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	svc := service.InitUserService()
	return svc.Login(req)
}
