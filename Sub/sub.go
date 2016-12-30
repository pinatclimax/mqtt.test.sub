package Sub

import (
	"fmt"
	"os"

	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//SubTestTopic ...
func SubTestTopic(c MQTT.Client, topic string, wg *sync.WaitGroup) {
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	fmt.Println("topic:", topic)
	wg.Done()
}
