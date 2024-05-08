package mqttutil

import (
	"fmt"
	"testing"
)

func TestMqttUtil(t *testing.T) {
	client, isConnect := NewMqttClient("127.0.0.1:1883", "", "", "mytest", 60)
	if !isConnect {
		return
	}
	fmt.Println("连接成功")
	client.DoPublish("test", "hello123")
	client.DoSubscribe("test", DefaultJsonPrint())
}
