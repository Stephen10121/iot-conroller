package messageqeue

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/stephen10121/iot-conroller/config"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func Publish(client mqtt.Client, message string) {
	client.Publish("topic/test", 0, false, message)
}

var subscribers = make(map[string]mqtt.Token)

func Subscribe(client mqtt.Client, topic string, connections map[string]*websocket.Conn) {
	token := client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		for i := range connections {
			err := connections[i].WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+m.Topic()+`","data":"`+string(m.Payload())+`"}`))

			if err != nil {
				log.Println(err)
			}
		}
	})
	token.Wait()
	subscribers[topic] = token

	log.Println("Subscribed to topic:", topic)
}

func Unsubscribe(topic string) {
	subscribers[topic].Done()

	delete(subscribers, topic)

	log.Println("Unubscribed to topic:", topic)
}

func setClientOption() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.BrokerAddress, config.BrokerPort))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(messagePubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts
}

func Initialize() mqtt.Client {
	opts := setClientOption()
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
