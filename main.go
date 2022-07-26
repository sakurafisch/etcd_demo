package main

import (
	"log"
	"net/http"
	"os/exec"

	"git.tdology.com/sakurafisch/etcd_demo/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/ginS"
	"go.etcd.io/etcd/server/v3/etcdmain"
)

func init() {
	clean()
	go startEtcd()
}

func main() {
	defer util.Close()
	setup_ginS()
	ginS.Run(":9503")
}

func clean() {
	clean := exec.Command("python", "clean.py")
	err := clean.Run()
	if err != nil {
		log.Fatalln("fail to run clean.py")
	}
}

func startEtcd() {
	etcdmain.Main([]string{""})
}

func setup_ginS() {
	ginS.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Etcd!")
	})
	ginS.Any("/set", func(c *gin.Context) {
		key := checkParamEmpty("key", c)
		if key == "" {
			c.Abort()
			return
		}
		value := checkParamEmpty("value", c)
		if value == "" {
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, util.SetKey(key, value))
		c.Abort()
	})
	ginS.Any("/get", func(c *gin.Context) {
		key := checkParamEmpty("key", c)
		if key == "" {
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, util.GetKey(key))
		c.Abort()
	})
	ginS.Any("/all", func(c *gin.Context) {
		c.JSON(http.StatusOK, util.GetAll())
		c.Abort()
	})
	ginS.Any("/delete", func(c *gin.Context) {
		key := checkParamEmpty("key", c)
		if key == "" {
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, util.DeteleKey(key))
	})
}

func checkParamEmpty(paramName string, c *gin.Context) string {
	param := c.Param(paramName)
	if param == "" {
		param = c.Query(paramName)
	}
	if param == "" {
		param = c.PostForm(paramName)
	}
	if param == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": paramName + " can not be empty",
		})
		c.Abort()
		return ""
	}
	return param
}
