package main

import (
	"context"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/ginS"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/etcdmain"
)

var etcd_client *clientv3.Client

func init() {
	clean()
	go start_etcd()
	setup_client()
}

func main() {
	defer close()
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
				"msg": "keycan not be empty",
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
		c.JSON(http.StatusOK, set_key(key, value))
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
		c.JSON(http.StatusOK, get_key(key))
		c.Abort()
	})
	ginS.Run(":9503")
}

func clean() {
	clean := exec.Command("python", "clean.py")
	err := clean.Run()
	if err != nil {
		log.Fatalln("fail to run clean.py")
	}
}

func start_etcd() {
	etcdmain.Main([]string{""})
}

func setup_client() {
	var err error
	etcd_client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func close() {
	etcd_client.Close()
}

func set_key(key string, value string) *clientv3.PutResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v\n", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
	return resp
}

func get_key(key string) *clientv3.GetResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Get(ctx, key)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v\n", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}
	return resp
}
