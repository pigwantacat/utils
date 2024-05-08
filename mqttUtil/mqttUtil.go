package mqttutil

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// MqttClient mqtt客户端
type MqttClient struct {
	// mqtt客户端
	mqttClient *mqtt.Client
}

// NewMqttClient 创建mqtt客户端
// @param broker mqtt服务端地址
// @param username 用户名
// @param password 密码
// @param clientId 客户端id
// @param timeOut 超时时间
// @return *MqttClient, bool mqtt客户端,连接是否成功
func NewMqttClient(broker string, username string, password string, clientId string, timeOut int64) (client *MqttClient, isConnect bool) {
	// 配置客户端
	clientOptions := mqtt.NewClientOptions().
		AddBroker(broker).
		SetUsername(username).
		SetPassword(password).
		SetClientID(clientId).
		SetConnectTimeout(time.Duration(timeOut) * time.Second)
	// 创建客户端
	c := mqtt.NewClient(clientOptions)
	// 测试连接
	if token := c.Connect(); token.WaitTimeout(time.Duration(timeOut)*time.Second) && token.Wait() && token.Error() != nil {
		return nil, false
	}
	return &MqttClient{mqttClient: &c}, true
}

// DoSubscribe 订阅mqtt服务端中指定的topic
// @param topic 订阅的主题
// @param doMessage 消息处理函数
func (client *MqttClient) DoSubscribe(topic string, doMessage mqtt.MessageHandler) {
	for {
		token := (*client.mqttClient).Subscribe(topic, 1, doMessage)
		token.Wait()
	}
}

// DoPublish 往mqtt服务端中指定的topic推送数据
// @param topic 推送的主题
// @param content 推送的数据
func (client *MqttClient) DoPublish(topic string, content interface{}) {
	// 往主题推送数据
	(*client.mqttClient).Publish(topic, 1, false, content)
}

func DefaultJsonPrint() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var resMap map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &resMap); err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(resMap)
		}
	}
}
