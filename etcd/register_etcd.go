package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fellowme/gin_common_library/config"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

var cli *clientv3.Client

func InitEtcd() {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(config.ServerConfigSettings.EtcdConfig.EtcdAddress, ";"),
		DialTimeout: config.ServerConfigSettings.EtcdConfig.EtcdDialTimeout * time.Second,
		Logger:      zap.L(),
	})
	if err != nil {
		zap.L().Error("InitEtcd error", zap.Any("error", err), zap.Any("config", config.ServerConfigSettings.EtcdConfig))
	}
}

type Register struct {
	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse
	srvInfo     Server
	srvTTL      int64
	cli         *clientv3.Client
	logger      *zap.Logger
}

// NewRegister create a register base on etcd
func NewRegister(logger *zap.Logger) *Register {
	return &Register{
		logger: logger,
	}
}

// Register a service
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}
	r.cli = cli
	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive()

	return r.closeCh, nil
}

// Stop  register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), config.ServerConfigSettings.EtcdConfig.EtcdDialTimeout*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	r.leasesID = leaseResp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}
