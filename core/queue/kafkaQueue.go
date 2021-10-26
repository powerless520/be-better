package queue

import (
	"be-better/core/global"
	"bytes"
	"context"
	"encoding/json"
	"gopkg.in/Shopify/sarama.v1"
	"reflect"
	"sync"
	"time"
)

type KafkaHandler struct {
	KafkaClient sarama.ConsumerGroup
	finish      chan struct{}
	timeOut     time.Duration
	wg          sync.WaitGroup
	JobNum      int
}

var _ Handler = (*KafkaHandler)(nil)

func (k *KafkaHandler) Push(queue *Queue) error {
	panic("implement me")
}

func (k *KafkaHandler) Consume(c *Consumer) error {
	global.GlobalLogger.Infof("Kafka consumer handler %s - %s start...", c.QueueName, c.groupId)
	var exit = false
	k.finish = make(chan struct{}, 1)
	k.wg.Add(1)
	go func() {
		<-c.ch
		exit = true
		k.finish <- struct{}{}
		k.wg.Done()
	}()

	for i := 0; i < k.JobNum; i++ {
		k.wg.Add(1)
		go k.consume(c)
	}

	consume := KafkaConsumer{
		ready:    make(chan bool),
		consumer: c,
	}

	if c.mode == 1 {
		consume.handler = k.batchConsume
	} else {
		consume.handler = k.singleConsume
	}

	for !exit {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		err := k.KafkaClient.Consume(ctx, []string{c.QueueName}, &consume)
		if ctx.Err() != nil {
			cancelFunc()
			continue
		}
		if err != nil {
			global.GlobalLogger.Errorf("KafkaConsumeHandle read message err: %v\n", err)
			cancelFunc()
		}

		consume.ready = make(chan bool)
		cancelFunc()
	}

	k.wg.Wait()
	k.KafkaClient.Close()
	global.GlobalLogger.Infof("Kafka consumer handler %s - %s stop...", c.QueueName, c.groupId)
	return nil
}

func (k *KafkaHandler) singleConsume(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, c *Consumer) {
	for message := range claim.Messages() {
		if message.Offset <= 0 {
			continue
		}
		messageJson, _ := json.Marshal(message)

		messageData := map[string]interface{}{}
		dec := json.NewDecoder(bytes.NewBuffer(message.Value))
		dec.UseNumber()
		if err := dec.Decode(&messageData); err != nil {
			global.GlobalLogger.Errorf("KafkaConsumerHandle json decode message err: %v, data: %v\n", err, string(message.Value))
			continue
		}

		eventName := messageData["event"].(string)
		if _, ok := c.fun[eventName]; ok {
			global.GlobalLogger.Debug("KafkaConsumeHandle message: ", string(messageJson)+"\n"+string(message.Value))
			continue
		}

		if _, ok := c.fun[eventName]; ok {
			global.GlobalLogger.Errorf("KafkaConsumeHandle message: ", string(messageJson)+"\n"+string(message.Value))
		}

		propertiesByte, _ := json.Marshal(messageData["properties"].(map[string]interface{}))
		c.kafkaChan <- map[string][]byte{
			eventName: propertiesByte,
		}
		session.MarkMessage(message, "")
	}

}

func (k *KafkaHandler) batchConsume(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, c *Consumer) {
	messageMap := make(map[string]map[string][]interface{})
	for message := range claim.Messages() {
		if message.Offset <= 0 {
			continue
		}
		messageJson, _ := json.Marshal(message)
		messageData := map[string]interface{}{}
		dec := json.NewDecoder(bytes.NewBuffer(message.Value))
		dec.UseNumber()
		if err := dec.Decode(&messageData); err != nil {
			global.GlobalLogger.Errorf("KafkaConsumeHandle json decode message err: %v, data: %v\n", err, string(message.Value))
			continue
		}

		eventName := messageData["event"].(string)
		if eventName == "" {
			session.MarkMessage(message, "")
			continue
		}

		if _, ok := c.fun[eventName]; !ok {
			session.MarkMessage(message, "")
			continue
		}

		global.GlobalLogger.Debug("KafkaConsumeHandle batch message: ", string(messageJson)+"\n"+string(message.Value))

		properties := messageData["properties"].(map[string]interface{})
		biz, ok := properties["official_bizid"]
		if !ok || biz == nil {
			session.MarkMessage(message, "")
			continue
		}

		bizId := biz.(string)

		if len(messageMap[eventName]) == 0 {
			messageMap[eventName] = make(map[string][]interface{}, 0)
		}
		if len(messageMap[eventName][bizId]) == 0 {
			messageMap[eventName][bizId] = make([]interface{}, 0)
		}
		messageMap[eventName][bizId] = append(messageMap[eventName][bizId], properties)

		for eventName, data := range messageMap {
			//当任何一个组的数据达到totalMessage的时候，把所有组发送出去
			flag := false
			for _, li := range data {
				if len(li) >= c.totalMessage {
					flag = true
					break
				}
			}
			if flag {
				k.sendAllToChannel(eventName, data, c)
				delete(messageMap, eventName)
			}
		}
		session.MarkMessage(message, "")
	}
	k.sendEventAllToChannel(messageMap, c)
}

func (k *KafkaHandler) sendEventAllToChannel(eventDataList map[string]map[string][]interface{}, c *Consumer) {
	if len(eventDataList) == 0 {
		return
	}
	for eventName, data := range eventDataList {
		k.sendAllToChannel(eventName, data, c)
	}
}

func (k *KafkaHandler) sendAllToChannel(eventName string, data map[string][]interface{}, c *Consumer) {
	if len(data) == 0 {
		return
	}

	for _, v := range data {
		if len(v) != 0 {
			d, _ := json.Marshal(v)
			c.kafkaChan <- map[string][]byte{
				eventName: d,
			}
		}
	}
}

func (k *KafkaHandler) consume(c *Consumer) {
	global.GlobalLogger.Infof("handler consume start ...")
	var finish = false

	for !finish {
		select {
		case data := <-c.kafkaChan:
			k.consumeHandle(data, c)
		case <-k.finish:
			k.flushAll(c)
			finish = true
		}
	}
	global.GlobalLogger.Infof("handler consume stop...")
	k.wg.Done()
}

func (k *KafkaHandler) consumeHandle(data map[string][]byte, c *Consumer) {
	for k1, v1 := range data {
		job := Job{
			JobHandler: nil,
			JobName:    k1,
			JobData:    string(v1),
		}

		handler, exist := c.fun[job.JobName]
		if !exist {
			continue
		}

		global.GlobalLogger.Debugf("KafkaConsumeHandle job data: %s", string(v1))
		job.JobHandler = reflect.ValueOf(handler).Elem().Interface()
		if err := job.Execute(); err != nil {
			global.GlobalLogger.Errorf("KafkaConumeHandle execute job err: %v\n", err)
		}

	}
}

func (k *KafkaHandler) flushAll(c *Consumer) {
	for len(c.kafkaChan) != 0 {
		select {
		case data := <-c.kafkaChan:
			k.consumeHandle(data, c)
		}
	}
}
