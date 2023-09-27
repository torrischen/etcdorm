package eorm

import (
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type processor func(*Etcdb)

func put() processor {
	return func(db *Etcdb) {
		txn := db.client.Txn(db.ctx)
		ifClause := make([]clientv3.Cmp, 0)
		thenClause := make([]clientv3.Op, 0)

		stmt := db.Statement
		for i := 0; i < len(stmt.EntryList); i++ {
			op := clientv3.OpPut(stmt.EntryList[i].Key, stmt.EntryList[i].Value, stmt.OpOptions...)
			thenClause = append(thenClause, op)
		}

		resp, err := txn.If(ifClause...).Then(thenClause...).Commit()
		if err != nil {
			db.Error = err
			return
		}
		if resp.Succeeded != true {
			db.Error = errors.New("put error")
			return
		}
		return
	}
}

func get() processor {
	return func(db *Etcdb) {
		stmt := db.Statement
		resp, err := db.client.KV.Get(db.ctx, stmt.Cond, stmt.OpOptions...)
		if err != nil {
			db.Error = err
			return
		}
		for _, kv := range resp.Kvs {
			e := Entry{
				Key:            btos(kv.Key),
				Value:          btos(kv.Value),
				Version:        kv.Version,
				CreateRevision: kv.CreateRevision,
				ModRevision:    kv.ModRevision,
				Lease:          kv.Lease,
			}
			stmt.EntryList = append(stmt.EntryList, e)
		}
		return
	}
}

func drop() processor {
	return func(db *Etcdb) {
		stmt := db.Statement
		_, err := db.client.KV.Delete(db.ctx, stmt.Cond, stmt.OpOptions...)
		if err != nil {
			db.Error = err
			return
		}
	}
}

func grant() processor {
	return func(db *Etcdb) {
		stmt := db.Statement
		resp, err := db.client.Lease.Grant(db.ctx, stmt.TTL)
		if err != nil {
			db.Error = err
			return
		}
		stmt.LeaseID = int64(resp.ID)
	}
}
