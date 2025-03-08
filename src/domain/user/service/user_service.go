package service

import (
	"github.com/magicnana999/ddd-go/domain/user/cache"
	"github.com/magicnana999/ddd-go/domain/user/repository"
	"github.com/magicnana999/ddd-go/dto/user"
	"sync"
)

type UserService interface {
	Login(req *user.LoginRequest) (*user.LoginResponse, error)
}

var (
	UserServiceInstance UserService
	usiOnce             sync.Once
)

type DefaultUserService struct {
	userCache      *cache.UserCache
	userRepository *repository.UserRepository
}

func InitUserService() UserService {

	usiOnce.Do(func() {
		UserServiceInstance = &DefaultUserService{
			userCache:      cache.InitUserCache(),
			userRepository: repository.InitUserRepository(),
		}
	})
	return UserServiceInstance
}

func (s *DefaultUserService) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	return &user.LoginResponse{Token: "token"}, nil
}
