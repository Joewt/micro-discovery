package discovery

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var _ registry.Watcher = &watcher{}


type watcher struct {
	key       string
	ctx       context.Context
	cancel    context.CancelFunc
	watchChan clientv3.WatchChan
	watcher   clientv3.Watcher
	kv        clientv3.KV

	// is first use Next() ?
	first bool
}

func newWatcher(ctx context.Context, key string, client *clientv3.Client) *watcher {
	w := &watcher{
		key:     key,
		watcher: clientv3.NewWatcher(client),
		kv:      clientv3.NewKV(client),
		first:   true,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0))
	w.watcher.RequestProgress(context.Background())
	return w
}

func (w *watcher) Next() ([]*registry.ServiceInstance, error) {
	if w.first {
		item, err := w.getInstance()
		w.first = false
		return item, err
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		return w.getInstance()
	}

}

func (w *watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *watcher) getInstance() ([]*registry.ServiceInstance, error) {
	resp, err := w.kv.Get(w.ctx, w.key, clientv3.WithPrefix(), clientv3.WithSerializable())
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
