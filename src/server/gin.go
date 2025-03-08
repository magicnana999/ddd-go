package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magicnana999/ddd-go/controllers/user"
	"net/http"
	"reflect"
)

func Start() {
	r := gin.Default()

	userGroup := r.Group("/user")
	userGroup.POST("/login", bind(user.Login))

	r.Run(":8080")
}

func bind(controller any) gin.HandlerFunc {
	return func(c *gin.Context) {
		controllerType := reflect.TypeOf(controller)
		if controllerType.Kind() != reflect.Func {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "controller must be a function"})
			return
		}

		if controllerType.NumIn() < 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "controller must have at least one input parameter"})
			return
		}

		reqType := controllerType.In(0) // 获取第一个参数类型
		if reqType.Kind() != reflect.Ptr {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "controller input must be a pointer to a struct"})
			return
		}

		reqValue := reflect.New(reqType.Elem()).Interface() // 创建 `req` 结构体实例

		// 绑定数据
		if err := c.ShouldBind(reqValue); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to bind request: %v", err)})
			return
		}

		results := reflect.ValueOf(controller).Call([]reflect.Value{reflect.ValueOf(reqValue)})

		var response any
		var err error
		if len(results) > 0 {
			response = results[0].Interface() // 第一个返回值
		}
		if len(results) > 1 && !results[1].IsNil() {
			err = results[1].Interface().(error) // 第二个返回值
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 返回 `response`
		c.JSON(http.StatusOK, response)
	}
}
