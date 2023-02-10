package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

func main() {
	r := gin.Default()

	// 请求路由
	r.GET("/get", func(c *gin.Context) { c.JSON(200, "get") })
	r.POST("/post", func(c *gin.Context) { c.JSON(200, "post") })
	r.Any("/any", func(c *gin.Context) { c.JSON(200, "any") })

	// 绑定静态文件夹
	r.Static("/assets", "./assets")
	r.StaticFS("/static", http.Dir("static"))
	r.StaticFile("/favicon.ico", "./favicon.ico")

	// 获取路径参数
	r.GET("/:name/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"name": c.Param("name"), "id": c.Param("id")})
	})

	// 泛绑定
	r.GET("/user/*action", func(c *gin.Context) { c.String(200, "hello world") })
	r.POST("/*action", func(c *gin.Context) { c.String(200, "hello world") })

	// 获取请求参数,获取 GET 参数
	r.GET("/test", func(c *gin.Context) {
		firstName := c.Query("first_name")
		lastName := c.DefaultQuery("last_name", "default_last_name")
		c.String(200, "%s\n%s", firstName, lastName)
	})

	// 获取 Body 内容
	// curl -X POST 'http://127.0.0.1:8080/test' -d "first_name=wang&last_name=kai"
	r.POST("/test", func(c *gin.Context) {
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "[ReadAll] Error:", err)
			c.Abort()
		}

		// 若想拿到 first_name 和 last_name,需要将 bytes 回写至 body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		firstName := c.PostForm("first_name")
		lastName := c.DefaultPostForm("last_name", "default_last_name")

		c.String(http.StatusOK, "%s,%s\n%s", firstName, lastName, bodyBytes)
	})

	// 获取 bind 参数
	// curl -H "Content-Type:application/json"
	// -X POST 'http://127.0.0.1:8080/testing'
	// -d '{"name":"hc","address":"shanghai","birthday":"1999-01-01"}'
	r.POST("/bindingParams", bindingParams)

	// 验证请求参数

	// 结构体验证参数

	r.GET("/bindingParams", func(c *gin.Context) {
		people := People{}

		if err := c.ShouldBind(&people); err != nil {
			c.String(500, "people bind error:%v\n", err)
			c.Abort() // 此处停止,防止程序继续往下走
			return
		}
		c.String(200, "%v", people)
	})

	// 自定义验证规则
	// curl -X GET 'http://127.0.0.1:8080/bookable?check_in=2021-06-30&check_out=2021-07-01'
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookableDate", bookableDate)
	}
	r.GET("/bookable", func(c *gin.Context) {
		b := Booking{}
		if err := c.ShouldBindWith(&b, binding.Query); err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message": "ok!",
			"booking": b,
		})
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

func bindingParams(c *gin.Context) {
	person := Person{}

	// according to content-type to do different binding option
	if err := c.ShouldBind(&person); err != nil {
		c.String(500, "binding error:%v", err)
	}
	c.String(200, "%v", person)
}

func customFunc(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	if date, ok := field.Interface().(time.Time); ok {
		if date.Unix() > time.Now().Unix() {
			return true
		}
	}
	return false
}

func bookableDate(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	if date, ok := field.Interface().(time.Time); ok {
		if date.Unix() > time.Now().Unix() {
			return true
		}
	}
	return false
}
