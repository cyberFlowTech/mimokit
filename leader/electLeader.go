package leader

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
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
	ctx := context.TODO()
	if err := e.Campaign(ctx, "candidate-id"); err != nil {
		logx.Errorf("Failed to create session: %v", err)
		return
	}

	config.OnBecameLeader()
}
