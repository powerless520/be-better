package idGenerator

import (
	"errors"
	"facm/utils/idGenerator/snowFlake"
	"facm/utils/idGenerator/worker"
)

var namespaceRegistry = make(map[string]IdGenerator)

type IdGenerator struct {
	Worker *snowFlake.Worker
}

//var worker, err = snowFlake.NewWorker(1)
//func NextId() (int64, error) {
//	return worker.GetId(), err
//}

func NewGenerator(namespace string, workerId int64, dataCenterId int64) (*IdGenerator, error){
	_, ok := namespaceRegistry[namespace]
	if ok {
		return nil, errors.New("the namespace " +  namespace + " has been new IdGenerator object!")
	}

	var worker, err = snowFlake.NewWorker(workerId, dataCenterId)
	if err != nil{
		return nil,  err
	}

	idGenerator := IdGenerator{Worker: worker}
	namespaceRegistry[namespace] = idGenerator
	return &idGenerator, nil
}

func NewGeneratorAssigner(namespace string, workerIdAssigner worker.WorkerIdAssigner, dataCenterId int64) (*IdGenerator, error){
	_, ok := namespaceRegistry[namespace]
	if ok {
		return nil, errors.New("the namespace " +  namespace + " has been new IdGenerator object!")
	}

	var worker, err = snowFlake.NewWorkerAssigner(workerIdAssigner, dataCenterId)
	if err != nil{
		return nil,  err
	}

	idGenerator := IdGenerator{Worker: worker}
	namespaceRegistry[namespace] = idGenerator
	return &idGenerator, nil
}

func NsNextId(namespace string) (int64, error) {
	idGenerator, ok := namespaceRegistry[namespace]
	if !ok {
		return 0, errors.New("the namespace " +  namespace + " not init!")
	}

	return idGenerator.NextId()
}

func (t IdGenerator) NextId() (int64, error) {
	return t.Worker.GetId(), nil
}



