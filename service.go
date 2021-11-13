package discovery

import (
	"context"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ registry.Discovery = &Discovery{}

type Option func(o *options)

type options struct {
	prefix   string
	timeout  time.Duration
	maxRetry int
}

func Prefix(ns string) Option {
	return func(o *options) { o.prefix = ns }
}

func Timeout(ns time.Duration) Option {
	return func(o *options) { o.timeout = ns }
}

func MaxRetry(num int) Option {
	return func(o *options) { o.maxRetry = num }
}

type Discovery struct {
	opts   *options
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func NewGoMicro(client *clientv3.Client, opts ...Option) (r *Discovery) {
	options := &options{
		prefix:   prefix,
		timeout:  time.Second * 15,
		maxRetry: 5,
	}
	for _, o := range opts {
		o(options)
	}
	return &Discovery{
		opts:   options,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

func (d *Discovery) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	key := path.Join(d.opts.prefix, strings.Replace(name, "/", "-", -1))+"/"
	resp, err := d.kv.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}
	var items []*registry.ServiceInstance
	for _, kv := range resp.Kvs {
		if sn, err := decode(kv.Value); err == nil {
			var endpoints []string
			for _, v := range sn.Nodes {
				endpoints = append(endpoints, v.Address)
			}

			si := &registry.ServiceInstance{
				Name:      sn.Name,
				Version:   sn.Version,
				Metadata:  sn.Metadata,
				Endpoints: endpoints,
			}
			items = append(items, si)
		}
	}
	return items, nil
}

func (d *Discovery) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := path.Join(d.opts.prefix, strings.Replace(name, "/", "-", -1))+"/"
	return newWatcher(ctx, key, d.client), nil
}

func marshal(si *registry.ServiceInstance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *registry.ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
