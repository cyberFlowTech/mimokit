package leader

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestElectLeader(t *testing.T) {
	etcdEndpoints := []string{"192.168.31.57:2379"}
	// 创建 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: etcdEndpoints,
	})
	if err != nil {
		t.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()
	type args struct {
		config ElectionConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "elect",
			args: args{config: ElectionConfig{
				EtcdClient:     etcdClient,
				ElectionKey:    "test",
				LeaseTTL:       5,
				OnBecameLeader: nil,
			}},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ElectLeader(tt.args.config)
		})
	}
}
