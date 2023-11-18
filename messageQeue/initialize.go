package messageqeue

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/stephen10121/iot-conroller/types"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func Publish(client mqtt.Client, message string) {
	client.Publish("topic/test", 0, false, message)
}

func Subscribe(client mqtt.Client, topic string, responder map[string]types.ResponderMap) {
	token := client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		val, ok := responder[m.Topic()]
		// If the key exists
		if ok {
			for i := range val.Recievers {
				err := val.Recievers[i].WriteMessage(websocket.TextMessage, []byte(`{"error":false,"command":"`+m.Topic()+`","data":`+string(m.Payload())+`}`))

				if err != nil {
					log.Println(err)
				}
			}
		}
	})
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
	fmt.Println()
}

func setClientOption() *mqtt.ClientOptions {
	var broker = "192.168.0.27"
	var port = 1883
	opts := mqtt.NewClientOptions()

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
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

	//Subscribe(client, "topic/test")
	return client
}
