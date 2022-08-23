package cnt

import (
	"encoding/json"
	"expvar"
	"github.com/pkg/errors"
	"sync"
)

var kafka *counterKafka

type counterKafka struct {
	kafka *kafkaStruct
	m     *sync.RWMutex
}

type kafkaStruct struct {
	Total   uint `json:"total"`
	Error   uint `json:"error"`
	Success uint `json:"success"`
}

func init() {
	kafka = &counterKafka{
		m:     &sync.RWMutex{},
		kafka: &kafkaStruct{},
	}
	expvar.Publish("Consumer", kafka)
}

func (c *counterKafka) incTotal() {
	c.m.Lock()
	defer c.m.Unlock()
	c.kafka.Total++
}

func (c *counterKafka) incError() {
	c.m.Lock()
	defer c.m.Unlock()
	c.kafka.Error++
}

func (c *counterKafka) incSuccess() {
	c.m.Lock()
	defer c.m.Unlock()
	c.kafka.Success++
}

func (c *counterKafka) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	data, err := json.Marshal(c.kafka)
	if err != nil {
		return errors.Wrap(err, "bad json parse").Error()
	}

	return string(data)
}

func IncSuccess() {
	kafka.incSuccess()
}

func IncError() {
	kafka.incError()
}

func IncTotal() {
	kafka.incTotal()
}
