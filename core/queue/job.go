package queue

import (
	"encoding/json"
	"errors"
	"reflect"
)

type BaseJob interface {
	Handle() error
	New([]byte) (interface{}, error)
}

type Job struct {
	JobHandler interface{}
	JobData    string `json:"job_data"`
	JobName    string `json:"job_name"`
	Attempts   int    `json:"attempts"`
	MaxTries   int    `json:"max_tries"`
}

func (j Job) Execute() error {
	var err error

	for {
		if j.Attempts > j.MaxTries {
			break
		}
		f := reflect.ValueOf(j.JobHandler).MethodByName("New")

		if f.Kind().String() == "invalid" || f.Kind().String() != "func" {
			err = errors.New("Job execute no method err: " + j.JobName + "|handle")
			j.Attempts++
			continue
		}

		params := make([]reflect.Value, 1)

		params[0] = reflect.ValueOf([]byte(j.JobData))
		callResultValue := f.Call(params)

		if callResultValue == nil || len(callResultValue) != 2 {
			err = errors.New("Job execute interface return error: " + j.JobName)
			j.Attempts++
			continue
		}

		baseJob := callResultValue[0].Interface()
		if baseJob != nil {
			if err = baseJob.(BaseJob).Handle(); err == nil {
				break
			}
		} else {
			err = callResultValue[1].Interface().(error)
		}

		j.Attempts++

	}
	return err
}

func NewJob(jobName string, data interface{}) Job {
	b, _ := json.Marshal(data)
	return Job{
		JobHandler: nil,
		JobData:    string(b),
		JobName:    jobName,
		MaxTries:   1,
	}
}

func (j *Job) SetMaxTries(num int) {
	j.MaxTries = num
}
