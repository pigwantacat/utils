package redisutil

import (
	"fmt"
	"testing"
)

func TestRedisClient(t *testing.T) {
	client, isConnect := NewRedisClient("127.0.0.1", 6379, "123456", 0)
	if isConnect {
		client.Set("test", "test")
		fmt.Println(client.GetString("test"))
	}
}
