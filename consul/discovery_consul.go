package consul

import (
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"time"
)

func discoveryServices(serviceName string, q *consul.QueryOptions) ([]consul.AgentServiceChecksInfo, error) {
	if clientAgent == nil {
		zap.L().Error("DiscoveryService consul yaml error ")
		return nil, errors.New(consulYamlError)
	}
	_, data, err := clientAgent.AgentHealthServiceByNameOpts(serviceName, q)
	return data, err
}

// Resolver for grpc client
type Resolver struct {
	closeCh            chan struct{}
	agent              *consul.Agent
	serviceName        string
	serviceAddressList []resolver.Address
	cc                 resolver.ClientConn
	logger             *zap.Logger
	lastTime           time.Time
}

// NewResolver create a new resolver.Builder base on etcd
func NewResolver(serviceName string, logger *zap.Logger) *Resolver {
	return &Resolver{
		serviceName: serviceName,
		logger:      logger,
	}
}

// Scheme returns the scheme supported by this resolver.
func (r *Resolver) Scheme() string {
	return ""
}

// Build creates a new resolver.Resolver for the given target
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// ResolveNow resolver.Resolver interface
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {}

// Close resolver.Resolver interface
func (r *Resolver) Close() {
}

// start
func (r *Resolver) start() (chan<- struct{}, error) {
	r.agent = clientAgent
	resolver.Register(r)
	r.closeCh = make(chan struct{})
	if time.Now().Sub(r.lastTime).Seconds() >= 5.0 {
		err := r.sync()
		if err != nil {
			return r.closeCh, err
		}
	}
	return r.closeCh, nil
}

// sync 同步获取所有地址信息
func (r *Resolver) sync() error {
	_, res, err := r.agent.AgentHealthServiceByName(r.serviceName)
	if err != nil {
		return err
	}
	r.serviceAddressList = []resolver.Address{}
	for _, v := range res {
		addr := resolver.Address{Addr: fmt.Sprintf("%s:%d", v.Service.Address, v.Service.Port)}
		r.serviceAddressList = append(r.serviceAddressList, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: r.serviceAddressList})
	r.lastTime = time.Now()
	return nil
}
