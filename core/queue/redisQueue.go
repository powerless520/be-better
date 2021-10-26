package queue

import (
	"be-better/core/global"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"reflect"
	"sync"
	"time"
)

const RedisQueuePrefix = "server:queue:"

type RedisHandler struct {
	wg   sync.WaitGroup
	exit bool
}

var _ Handler = (*RedisHandler)(nil)

func (r *RedisHandler) Push(queue *Queue) error {
	if queue.delay > 0 {
		return r.pushDelay(queue)
	} else {
		return r.push(queue)
	}
}

func (r *RedisHandler) pushDelay(queue *Queue) error {
	key := RedisQueuePrefix + "delay:" + queue.QueueName
	score := float64(time.Now().Unix()) + queue.delay.Seconds()

	mapData, err := json.Marshal(queue)
	if err != nil {
		return err
	}

	_, err = global.GlobalRedis.ZAdd(context.Background(), key, &redis.Z{
		Score:  score,
		Member: mapData,
	}).Result()
	return err
}

func (r *RedisHandler) push(queue *Queue) error {
	queueKey := RedisQueuePrefix + queue.QueueName
	mapData, err := json.Marshal(queue)
	if err != nil {
		return err
	}
	_, err = global.GlobalRedis.RPush(context.Background(), queueKey, string(mapData)).Result()

	return err
}

func (r *RedisHandler) Consume(c *Consumer) error {
	global.GlobalLogger.Infof("Redis consumer handler %s start...", c.QueueName)

	r.exit = false
	r.wg.Add(1)
	go func() {
		<-c.ch
		r.exit = true
		r.wg.Done()
	}()

	queueKey := RedisQueuePrefix + c.QueueName
	queueDelayKey := RedisQueuePrefix + "delay:" + c.QueueName

	r.wg.Add(1)
	go r.consume(queueKey, c)

	r.wg.Add(1)
	go r.consumeDelay(queueDelayKey, queueKey)

	r.wg.Wait()
	global.GlobalLogger.Infof("Redis consumer handler %s stop...", c.QueueName)

	return nil
}

func (r *RedisHandler) consume(queueKey string, c *Consumer) {
	for !r.exit {
		result, err := global.GlobalRedis.BLPop(context.Background(), time.Second, queueKey).Result()
		if err != nil && err != redis.Nil {
			global.GlobalLogger.Errorf("RedisConsumeHandler lpop read data err: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		if err == redis.Nil {
			continue
		}

		if len(result) != 2 {
			global.GlobalLogger.Errorf("RedisConsumerHandler lpop data err: %v\n, data: %v", err, result)
			continue
		}

		global.GlobalLogger.Debugf("RedisConsumerHandler read message data: %v", result)
		if err = r.handle([]byte(result[1]), c); err != nil {
			global.GlobalLogger.Errorf("RedisConsumerHandler handle list data err: %v, data: %v,consumer: %v\n", err, result, c)
			continue
		}
	}

	r.wg.Done()
}

func (r *RedisHandler) consumeDelay(queueDelayKey string, queueKey string) {
	for !r.exit {
		timeNowInt := time.Now().Unix()
		lua := `--get all of the jobs with an expired "score"...
local val = redis.call('zrangebyscore',KEYS[1],'-inf',ARGV[1])
-- If wh have values in the array,we will remove them from the first queue
-- and add them onto the destination queue in chunks of 100,which moves
-- all of the appropriate jobs onto the destination queue very safely.
if (next(val) ~= nil) then
    redis.call('zremrangebyrank',KEYS[1],0,#val - 1)

    for i = 1, #val,100 do
        redis.call('rpush',KEYS[2],unpack(val,i,math.min(i+99,#val)))
    end
end

return val`
		err := global.GlobalRedis.Eval(context.Background(), lua, []string{queueDelayKey, queueKey}, timeNowInt).Err()
		time.Sleep(time.Second)
		if err != nil {
			global.GlobalLogger.Errorf("RedisConsumerHandler zrange data err: %v\n", err)
		}
	}
	r.wg.Done()

}

func (r *RedisHandler) handle(data []byte, c *Consumer) error {
	job, err := r.getJobHandler(data, c)
	if err != nil {
		global.GlobalLogger.Errorf("RedisConsumerHandler get job handler err: %v, data: %v, consumer: %v\n", err, string(data), c)
	}
	if job.JobHandler == nil {
		return nil
	}
	return job.Execute()
}

func (r *RedisHandler) getJobHandler(result []byte, c *Consumer) (Job, error) {
	job := Job{}
	queueData := map[string]interface{}{}
	if err := json.Unmarshal(result, &queueData); err != nil {
		return job, err
	}

	jobData, exist := queueData["job_handler"].(map[string]interface{})
	if !exist || jobData == nil {
		return job, errors.New("job handler not found in queue")
	}

	jobByte, err := json.Marshal(jobData)
	if err != nil {
		return job, err
	}
	if err = json.Unmarshal(jobByte, &job); err != nil {
		return job, err
	}

	handler, exist := c.fun[job.JobName]
	if !exist {
		return job, nil
	}

	job.JobHandler = reflect.New(handler).Elem().Interface()
	return job, nil
}
