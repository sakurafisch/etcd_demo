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
		key := c.Param("key")
		if key == "" {
			key = c.Query("key")
		}
		if key == "" {
			key = c.PostForm("key")
		}
		if key == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": "key can not be empty",
			})
			c.Abort()
			return
		}
		value := c.Param("value")
		if value == "" {
			value = c.Query("value")
		}
		if value == "" {
			value = c.PostForm("value")
		}
		if value == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": "value can not be empty",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, util.SetKey(key, value))
		c.Abort()
	})
	ginS.Any("/get", func(c *gin.Context) {
		key := c.Param("key")
		if key == "" {
			key = c.Query("key")
		}
		if key == "" {
			key = c.PostForm("key")
		}
		if key == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "key can not be empty",
			})
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
}
