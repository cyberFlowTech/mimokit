package leader

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
)

type ElectionConfig struct {
	EtcdClient     *clientv3.Client
	ElectionKey    string
	LeaseTTL       int
	OnBecameLeader func()
}

func ElectLeader(config ElectionConfig) {
	s, err := concurrency.NewSession(config.EtcdClient, concurrency.WithTTL(config.LeaseTTL))
	if err != nil {
		logx.Errorf("Failed to create session: %v", err)
		return
	}

	defer s.Close()

	mutex := concurrency.NewMutex(s, config.ElectionKey)
	// 尝试获取锁
	ctx := context.TODO()
	if err := mutex.Lock(ctx); err != nil {
		log.Fatalf("Failed to acquire lock: %v", err)
	}
	defer mutex.Unlock(ctx)
	// 如果成功获取锁，说明当前服务成为了 leader
	fmt.Println("You have been elected as leader")
	// 调用回调函数（如果存在）
	if config.OnBecameLeader != nil {
		config.OnBecameLeader()
	}

	// 保持 leader 状态
	select {}
}
