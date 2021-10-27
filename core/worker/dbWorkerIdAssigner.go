package worker

import (
	"be-better/app/models"
	"be-better/core/global"
	"be-better/core/global/model"
	"be-better/utils/dbutil"
	idGenerator2 "be-better/utils/idGenerator"
	"fmt"
	"net"
	"strconv"
	"time"
)

type DbWorkerIdAssigner struct {
	Namespace string
	quit          chan struct{}
}

func NewApplicationIdGenerator() *model.IdGenerators {
	facmWorker := DbWorkerIdAssigner{ Namespace: "facm"}

	userIdGenerator, _ := idGenerator2.NewGeneratorAssigner("users", facmWorker, global.GlobalConfig.DataCenterId)
	sessionIdGenerator, _ := idGenerator2.NewGeneratorAssigner("sessions", facmWorker, global.GlobalConfig.DataCenterId)

	go facmWorker.Heartbeat()

	return &model.IdGenerators{
		UserIdGenerator:    userIdGenerator,
		SessionIdGenerator: sessionIdGenerator,
	}
}

func (t DbWorkerIdAssigner) Heartbeat() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("DbWorkerIdAssigner Heartbeat recover err:%+v\n", err)
			return
		}
	}()

	ticker := time.NewTicker(time.Duration(5) * time.Minute)
	for {
		select {
		case <-t.quit:
			fmt.Println("DbWorkerIdAssigner exit ...")
			return
		case <-ticker.C:
			hostname := getLocalIP()
			workerNode := models.WorkerNode{}

			var err = global.GlobalDatabase.Where("namespace = ?", t.Namespace).Where("host_name", hostname).First(&workerNode).Error
			if err == nil{
				nowStr :=  time.Now().Format("2006-01-02 15:04:05")
				workerNode.UpdatedAt = &nowStr
				workerNode.AvailableAt = &nowStr
				err = global.GVA_DB.Save(workerNode).Error
				if err != nil{
					global.GVA_LOG.Error("DbWorkerIdAssigner heartbeat update error", err)
				}
			}
		}
	}
}

func (t DbWorkerIdAssigner) Close() {
	t.quit <- struct{}{}
}

func (t DbWorkerIdAssigner)AssignWorkerId() int64 {

	workerNode := models.WorkerNode{}
	hostname := getLocalIP()
	nowStr :=  time.Now().Format("2006-01-02 15:04:05")
	pre72HourStr :=  time.Now().Add(-time.Hour * 72).Format("2006-01-02 15:04:05")
	var err = global.GVA_DB.Where("namespace = ?", t.Namespace).Where("host_name", hostname).First(&workerNode).Error
	if err != nil{
		if dbutil.IsNoRecord(err) {

			var workerList []models.WorkerNode
			workerId := int64(-1)
			existsWorkIdMap := map[int64]int{}
			err = global.GVA_DB.Where("namespace = ?", t.Namespace).Where("available_at > ? or available_at is null", pre72HourStr).Find(&workerList).Error
			if err != nil{
				panic("assign worker id error: " + err.Error())
			}else{

				if len(workerList) > 0{

					for _, v := range workerList{
						existsWorkIdMap[v.WorkId] = 1
					}

					//这里仅找0-127(现在worker_id范围为0-127)
					for i := 0; i <= 127; i++ {
						i64 := int64(i)
						if _, ok := existsWorkIdMap[i64]; !ok {
							workerId = i64
							break
						}
					}
				}else{
					workerId = int64(0)
				}
			}

			if workerId <0 {
				panic("assign worker id error, work_id: " + strconv.FormatInt(workerId, 10))
			}

			err = global.GVA_DB.Where("namespace = ?", t.Namespace).Where("work_id = ?",  workerId).First(&workerNode).Error
			if err != nil{
				if dbutil.IsNoRecord(err) {
					//new work_id
					workerNode.HostName = hostname
					workerNode.Namespace = t.Namespace
					workerNode.CreatedAt = &nowStr
					workerNode.UpdatedAt = &nowStr
					workerNode.LaunchDate = &nowStr
					workerNode.AvailableAt = &nowStr
					workerNode.WorkId = workerId
					err := global.GVA_DB.Create(&workerNode).Error
					if err != nil{
						panic("assign worker id error: " + err.Error())
					}
				}
			}else{
				//存在记录，更新现有记录
				preHost := workerNode.HostName
				preAvailableAt := ""

				if workerNode.AvailableAt != nil{
					preAvailableAt = *workerNode.AvailableAt
				}

				workerNode.HostName = hostname
				workerNode.UpdatedAt = &nowStr
				workerNode.LaunchDate = &nowStr
				workerNode.AvailableAt = &nowStr
				err = global.GVA_DB.Save(workerNode).Error
				if err != nil{
					panic("assign worker id error: " + err.Error())
				}

				global.GVA_LOG.Info("DbWorkerIdAssigner host: " + preHost + " is not available, allocate the worker id:" + strconv.FormatInt(workerId, 10) +" to " + hostname +", pre available time is: " + preAvailableAt, err)
			}

			return workerNode.WorkId
		}
		panic("assign worker id error: " + err.Error())
	}

	//查询到结果，更新launchDate
	workerNode.UpdatedAt = &nowStr
	workerNode.LaunchDate = &nowStr
	workerNode.AvailableAt = &nowStr
	err = global.GVA_DB.Save(workerNode).Error
	if err != nil{
		panic("assign worker id error: " + err.Error())
	}

	return workerNode.WorkId
}

func (t DbWorkerIdAssigner) GetNamespace() string {
	return t.Namespace
}

func getLocalIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Get Local IP Error:" + err.Error())
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
