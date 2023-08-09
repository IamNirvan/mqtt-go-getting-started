package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// * this is how an anonymous function is used in Go.
// * This function takes a client and a message
// * Inside, it displays the topic and the payload in the message
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	// * This 2 lines below (26 & 27) allow the logging behavior in MQTT to be configured.
	// * The log.New() creates a new logger.
	// * Param #1: os.Stdout indicates that the logger must log data into the standard output (console)
	// * Param #2: the "" is the prefix that will accompany every log message. If it is empty like in this case, then no prefix will be added
	// * Param #3: specifies the flag settings. 0 means there are no additional flag settings that will be applied to the output
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// * mqtt.NewClientOptions() returns a client whose parameters can be configured. AKA ClientOption. This method applies default configurations
	// * AddBroker(broker) uses the returned client to configure the brokers the client will work with. Add broker will add the specified broker to the list of all brokers the client must work with. This too returns a ClientOption
	// * SetClientId(id) will use the returned ClientOptions to set a client ID for the client. The client will use this ID when connecting to the MQTT broker
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883").SetClientID("emqx_test_client")

	// * Configuring the keep alive time for the client
	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	// * Configuring the client to use 'f' to handle incoming messages.
	// *
	// * [ CHAT-GPT RESPONSE ]:
	// * This method is used to set a default handler function for incoming published messages.
	// * When the MQTT client publishes a message, it also subscribes to the topic it published to.
	// * This means that it can receive its own message (loopback).
	// * The default publish handler is called whenever the client receives a message on a topic it has subscribed to.
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// * mqtt.NewClient(opts) takes the ClientOptions and uses it to create a new MQTT client
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#NewClient
	c := mqtt.NewClient(opts)

	// * The MQTT client's connect() method must be called before using it. It returns a a value of type token
	// * This is because the connect() method ensures that the necessary resources are ready for the client before it can be used.
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang#NewClient
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// * c.Subscribe() starts a new subscription.
	// * Param #1: topic to subscribe to
	// * Param #2: Quality of Service (QoS) value
	// * Param #3: callback function that must be executed when a message is published in that aforementioned topic.
	// * 	If, like in this case 'nil' is provided, then the default callback function will be used
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang@v1.4.3#Client.Subscribe
	if token := c.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// * c.Publish() is used to publish a message under a specific topic with a specific QoS value. It is also possible to specify the message retention value
	// * Param #1: the topic to publish to
	// * Param #2: the QoS value
	// * Param #3: boolean value indicating whether to apply message retention
	// * Param #4: the payload for the message
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang@v1.4.3#Client.Subscribe
	token := c.Publish("testtopic/1", 0, false, "Hello World")
	token.Wait()

	time.Sleep(6 * time.Second)

	// * c.Unsubscribe() accepts a collection of topics that the client must unsubscribe from
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang@v1.4.3#Client.Subscribe
	if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// * c.Disconnect() first wait for the specified number of milliseconds and then disconnect from the server.
	// ? [ ALSO CHECK ]: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang@v1.4.3#Client.Subscribe
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}
