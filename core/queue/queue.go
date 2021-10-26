package queue

import (
	"be-better/core/global"
	"errors"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/Shopify/sarama.v1"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Queue struct {
	Id         string `json:"id"`
	queueType  string `json:"queue_type"`
	QueueName  string `json:"queue_name"`
	JobHandler Job    `json:"job_handler"`
	delay      time.Duration
	CreateAt   int64 `json:"create_at"`
}

type Consumer struct {
	Queue
	wg           sync.WaitGroup
	chNum        int
	ch           chan int
	fun          map[string]reflect.Type
	kafkaChan    chan map[string][]byte
	frequency    int
	totalMessage int
	groupId      string
	mode         int //1： 一次处理多条，默认0：处理一条
	handler      Handler
}

type Handler interface {
	Push(queue *Queue) error
	Consume(consumer *Consumer) error
}

func NewQueue(queueName string, job Job) *Queue {
	queue := new(Queue)
	queue.JobHandler = job
	queue.QueueName = queueName
	queue.CreateAt = time.Now().Unix()
	queue.Id = "id"
	queue.queueType = "redis"

	return queue
}

func (q *Queue) SetDelay(delay time.Duration) *Queue {
	if delay != 0 {
		q.delay = delay
	}
	return q
}

func (q *Queue) Push() error {
	var handler Handler

	switch q.queueType {
	case "redis":
		handler = &RedisHandler{}
	case "kafka":
		handler = &KafkaHandler{}
	default:
		return errors.New("不支持的类型")
	}
	return handler.Push(q)
}

func (q *Queue) generateId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func (c *Consumer) RegisterJob(event string, job BaseJob) {
	t := reflect.TypeOf(job).Elem()
	c.fun[event] = t
}

func NewConsumer(queueType string, queueName string, channelNum int, handler Handler) *Consumer {
	consumer := &Consumer{}
	consumer.QueueName = queueName
	consumer.queueType = queueType
	consumer.ch = make(chan int, channelNum)
	consumer.chNum = channelNum
	consumer.fun = make(map[string]reflect.Type)
	consumer.handler = handler

	switch queueType {
	case "kafka":
		consumer.kafkaChan = make(chan map[string][]byte, 200)
	}

	return consumer
}

func NewKafkaConsumer(topic string, groupId string, channelNum int, handlerChannelNum int) *Consumer {
	kfkCfg := global.GlobalConfig.Kafka
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kfkCfg.Username
	config.Net.SASL.Password = kfkCfg.Password
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = -2
	version, _ := sarama.ParseKafkaVersion("2.1.1")
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	client, err := sarama.NewConsumerGroup(kfkCfg.Url, groupId, config)
	if err != nil {
		panic("New Kafka client err")
	}
	kafkaHandler := &KafkaHandler{}

	return NewConsumer("kafka", topic, channelNum, kafkaHandler)
}
func NewRedisConsumer(queueName string, channelNum int) *Consumer {
	handler := &RedisHandler{}
	return NewConsumer("redis", queueName, channelNum, handler)
}

func (c *Consumer) SetFrequency(frequency int) *Consumer {
	c.frequency = frequency
	return c
}

func (c *Consumer) SetMode(mode, num int) *Consumer {
	c.mode = mode
	if mode == 1 {
		c.totalMessage = num
	} else {
		c.totalMessage = 0
	}
	return c
}

func (c *Consumer) SetGroupId(groupId string) *Consumer {
	c.groupId = groupId
	return c
}

func (c *Consumer) Consume() error {
	for i := 0; i < c.chNum; i++ {
		c.wg.Add(1)
		go func(fun Handler) {
			err := fun.Consume(c)
			if err != nil {
				global.GlobalLogger.Errorf("Queue consumer err: %v", err)
			}
			c.wg.Done()
		}(c.handler)
	}
	c.wg.Wait()
	return nil
}

func (c *Consumer) ShutDown(sig int) {
	if sig == 1 {
		for i := 0; i < c.chNum; i++ {
			c.ch <- 1
		}
	}
}
