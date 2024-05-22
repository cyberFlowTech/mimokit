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

	e := concurrency.NewElection(s, config.ElectionKey)
	// 查询已选举的对象列表
	resp, err := config.EtcdClient.Get(context.Background(), config.ElectionKey, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Failed to get elected objects: %v", err)
	}

	// 输出已选举的对象列表
	fmt.Println("Elected objects before election:")
	for _, kv := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", kv.Key, kv.Value)
	}

	ctx := context.TODO()
	if err := e.Campaign(ctx, "candidate-id"); err != nil {
		logx.Errorf("Failed to create session: %v", err)
		return
	}
	fmt.Println("you have been elected as leader")
	if config.OnBecameLeader != nil {
		config.OnBecameLeader()
	}

}
