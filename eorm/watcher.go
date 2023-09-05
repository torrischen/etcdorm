package eorm

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	ctx context.Context
	rch clientv3.WatchChan
}

//func StartWatchEtcdEvent() context.CancelFunc {
//	ctx, cancel := context.WithCancel(context.Background())
//	w := &watcher{
//		ctx: ctx,
//		rch: etcdc.Watch(ctx, "", clientv3.WithFromKey(), clientv3.WithPrevKV()),
//	}
//	w.Run()
//
//	return cancel
//}

func (w *watcher) Run() {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case msg := <-w.rch:
				for _, event := range msg.Events {
					if event.Type == mvccpb.PUT {
						go handlePutEvent(event)
					}
					if event.Type == mvccpb.DELETE {
						go handleDeleteEvent(event)
					}
				}
			}
		}
	}()
}

func handlePutEvent(event *clientv3.Event) {
	fmt.Printf("%s is put, version is %d\n", event.Kv.Key, event.Kv.Version)
}

func handleDeleteEvent(event *clientv3.Event) {
	fmt.Printf("%s is deleted\n", event.Kv.Key)
}
