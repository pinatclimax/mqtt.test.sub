package main

import (
	"fmt"
	"sync"
	"time"

	"climax.com/mqtt.test.sub/Sub"

	//import the Paho Go MQTT library

	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var mycount int32

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	mycount++
	fmt.Println("total count:", mycount)
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://10.15.8.129:1883")
	//opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	total := 80000
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(total)

	for i := 1; i <= total; i++ {
		go Sub.SubTestTopic(c, strconv.Itoa(i), &wg)
		//go Sub.SubTestTopic(c, strconv.Itoa(i))
	}
	// sillyname := randomdata.SillyName()
	// fmt.Println(sillyname)
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	// if token := c.Subscribe(sillyname, 0, nil); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// 	os.Exit(1)
	// }
	wg.Wait()
	fmt.Println("sub finished")
	fmt.Println("elapsed time: ", time.Since(start))
	//Publish 5 messages to /go-mqtt/sample at qos 1 and wait for the receipt
	//from the server after sending each message
	// time.Sleep(10 * time.Second)
	// for i := 0; i < 5; i++ {
	// 	text := fmt.Sprintf("this is msg #%d!", i)
	// 	token := c.Publish(strconv.Itoa(i), 0, false, text)
	// 	token.Wait()
	// }

	<-make(chan int)
}
