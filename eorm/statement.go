package eorm

import clientv3 "go.etcd.io/etcd/client/v3"

type Statement struct {
	EntryList []Entry
	OpOptions []clientv3.OpOption
	Cond      string
	LeaseID   int64
	TTL       int64
}
