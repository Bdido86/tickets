package kafka

import (
	"fmt"
	"log"
)

func (k *Kafka) checkAsyncProducer() {

	go func() {
		for msg := range k.asyncProducer.Errors() {
			log.Printf("%v", msg)
		}
	}()

	go func() {
		for res := range k.asyncProducer.Successes() {
			fmt.Printf("succ: %+v\n", res)
		}
	}()
}
