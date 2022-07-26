package util

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var etcd_client *clientv3.Client

func init() {
	setupClient()
}

func setupClient() {
	var err error
	etcd_client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func Close() {
	etcd_client.Close()
}

func SetKey(key string, value string) *clientv3.PutResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Put(ctx, key, value)
	cancel()
	handleError(err)
	return resp
}

func GetKey(key string) *clientv3.GetResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Get(ctx, key)
	cancel()
	handleError(err)
	return resp
}

func GetAll() *clientv3.GetResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Get(ctx, "", clientv3.WithPrefix())
	cancel()
	handleError(err)
	return resp
}

func DeteleKey(key string) *clientv3.DeleteResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Delete(ctx, key)
	cancel()
	handleError(err)
	return resp
}

func DeleteAll() *clientv3.DeleteResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcd_client.Delete(ctx, "", clientv3.WithPrefix())
	cancel()
	handleError(err)
	return resp
}

func handleError(err error) {
	if err == nil {
		return
	}
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
