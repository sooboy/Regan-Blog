package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type MuxteMap struct{}

func (m *MuxteMap) Get(key string)                    {}
func (m *MuxteMap) Set(key string, value interface{}) {}

type UnitFunc func(m *MuxteMap)
type UnitFuncs []UnitFuncs

func AuthLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in AuthLimit")
	}
}

func News() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in News")
		c.Set(
			"unitFn",
			setUnitFn(c, func(m *MuxteMap) {
				// 获取数据
				m.Set("some key", "some Value")
			}),
		)
	}
}

func Spot() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in Spot")
		c.Set(
			"unitFn",
			setUnitFn(c, func(m *MuxteMap) {
				// 获取数据
				m.Set("some key", "some Value")
			}),
		)
	}
}

func Future() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in Future")
		c.Set(
			"unitFn",
			setUnitFn(c, func(m *MuxteMap) {
				// 获取数据
				m.Set("some key", "some Value")
			}),
		)
	}
}

func Vote() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in Future")
		c.Set(
			"unitFn",
			setUnitFn(c, func(m *MuxteMap) {
				// 获取数据
				m.Set("some key", "some Value")
			}),
		)
	}
}

func HTML() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in HTML")
	}
}

func JSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  这里处理业务需求
		fmt.Println("this is in JSON")
	}
}

func main() {
	engine := gin.Default()

	admin := engine.Group("/admin", AuthLimit())
	{
		admin.GET("/index", News(), Spot(), Future(), HTML())
		admin.GET("/vote", News(), Spot(), Vote(), HTML())
	}

	api := engine.Group("/api", AuthLimit())
	{
		api.GET("/vote", Vote(), JSON())
	}
	engine.Run(":8080")
}

func setUnitFn(c *gin.Context, unit UnitFunc) UnitFuncs {
	return UnitFuncs{}
}
