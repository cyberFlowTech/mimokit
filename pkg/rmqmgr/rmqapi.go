package rmqmgr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type RmqMgr struct {
	conn  rmq.Connection
	queue rmq.Queue
}

type RmqEvent struct {
	UserID int64             `json:"user_id"`
	Event  string            `json:"event"`
	ETime  int64             `json:"etime"`
	Ext    map[string]string `json:"ext"`
}

const MIMOMQKey = "MimoRmqKey"

var errChan = make(chan error)

func NewRmqMgr(addrs string) *RmqMgr {
	fmt.Println("rmq addr:", addrs)
	r := &RmqMgr{}
	adds := strings.Split(addrs, ",")
	var err error
	if len(adds) == 1 {
		r.conn, err = rmq.OpenConnection("mimo-userups", "tcp", addrs, 1, errChan)
	} else {
		redisClusterOptions := &redis.ClusterOptions{
			Addrs: adds,
		}
		redisClusterClient := redis.NewClusterClient(redisClusterOptions)
		r.conn, err = rmq.OpenConnectionWithRedisClient("mimo-userups", redisClusterClient, errChan)
	}
	if err != nil {
		logx.Severef("rmq openConnection err:%s", err.Error())
	}
	r.OpenQueue()
	return r
}

func (r *RmqMgr) OpenQueue() (err error) {

	r.queue, err = r.conn.OpenQueue(MIMOMQKey)
	return err
}

func (r *RmqMgr) PublicEvent(msg RmqEvent) error {
	// create task
	taskBytes, err := json.Marshal(msg)
	if err != nil {
		// handle error
		return fmt.Errorf("marshal err %s", err.Error())
	}

	if r.queue == nil {
		r.OpenQueue()
	}
	if r.queue == nil {
		return fmt.Errorf("invalid queue")
	}
	return r.queue.PublishBytes(taskBytes)
}

func (r *RmqMgr) ConsumeEvent(consumer rmq.Consumer) error {
	if r.queue == nil {
		r.OpenQueue()
	}
	if r.queue == nil {
		return fmt.Errorf("invalid queue")
	}

	r.queue.StartConsuming(10, time.Second)

	//taskConsumer := &TaskConsumer{}
	name, err := r.queue.AddConsumer("task-consumer", consumer)
	if err != nil {
		logx.Errorf("addConsumer err:%s", err.Error())
		return err
	}
	logx.Infof("add consumer %s", name)
	return nil
}
