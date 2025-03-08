package repository

import (
	"github.com/magicnana999/ddd-go/infrastructure"
	"gorm.io/gorm"
	"sync"
)

var (
	UserRepositoryInstance *UserRepository
	uriOnce                sync.Once
)

type UserRepository struct {
	db *gorm.DB
}

func InitUserRepository() *UserRepository {
	uriOnce.Do(func() {
		UserRepositoryInstance = &UserRepository{
			db: infrastructure.InitGorm(),
		}
	})
	return UserRepositoryInstance
}
