package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"climax.com/mqtt.test.sub/Sub"

	//import the Paho Go MQTT library

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var mycount int32

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
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

	total := 10000
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(total)

	for i := 1; i <= total; i++ {
		tpanel, _, err := topicGenerator(i)
		fmt.Println(i, tpanel)
		go Sub.SubTestTopic(c, tpanel, &wg)

		if err != nil {
			fmt.Println(err)
		}
	}

	wg.Wait()
	fmt.Println("sub finished")
	fmt.Println("elapsed time: ", time.Since(start))

	<-make(chan int)
}

var macPrefix = "11:11:11:"
var panelTopic = "panel"
var userTopic = "user"
var ffffffNum = 16777215

func topicGenerator(num int) (tpanel string, tuser string, err error) {
	post, err := numberToMac(num)
	if err != nil {
		return "", "", err
	}
	tpanel = macPrefix + post + "_" + panelTopic
	tuser = macPrefix + post + "_" + userTopic
	return tpanel, tuser, nil
}

func userGenerator(num int) (user string, passwd string, err error) {
	post, err := numberToMac(num)
	if err != nil {
		return "", "", err
	}
	user = macPrefix + post + ":" + user
	passwd = macPrefix + post + ":" + passwd
	return user, passwd, nil
}

func numberToMac(num int) (string, error) {
	if num > 16777215 {
		return "", errors.New("number is greater than 16777215")
	}
	hexnum := fmt.Sprintf("%06x", num)
	postmac := fmt.Sprintf("%s:%s:%s", hexnum[0:2], hexnum[2:4], hexnum[4:6])
	//	log.Println("postmac:", postmac)
	return postmac, nil
}
