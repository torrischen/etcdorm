package eorm

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcdb struct {
	ctx        context.Context
	client     *clientv3.Client
	processors []processor
	Statement  *Statement
	clone      int
	Error      error
}

func Open(host string) *Etcdb {
	c, err := clientv3.New(clientv3.Config{
		Endpoints: []string{host},
	})
	if err != nil {
		panic(err)
	}

	return &Etcdb{
		ctx:        context.Background(),
		client:     c,
		processors: make([]processor, 0),
		Statement:  new(Statement),
		clone:      0,
		Error:      nil,
	}
}

func (etcdb *Etcdb) getInstance() *Etcdb {
	if etcdb.clone == 0 {
		tx := &Etcdb{
			ctx:        etcdb.ctx,
			client:     etcdb.client,
			processors: make([]processor, 0),
			Statement:  new(Statement),
			clone:      1,
			Error:      nil,
		}
		return tx
	} else {
		return etcdb
	}
}

// chainable api

func (etcdb *Etcdb) WithOptions(op ...clientv3.OpOption) *Etcdb {
	tx := etcdb.getInstance()

	tx.Statement.OpOptions = append(tx.Statement.OpOptions, op...)
	return tx
}

func (etcdb *Etcdb) Cond(cond string) *Etcdb {
	tx := etcdb.getInstance()

	tx.Statement.Cond = cond

	return tx
}

// finisher api

func (etcdb *Etcdb) exec() {
	for _, prcs := range etcdb.processors {
		prcs(etcdb)
	}
}

func (etcdb *Etcdb) Put(key, value string) *Etcdb {
	tx := etcdb.getInstance()

	e := Entry{
		Key:   key,
		Value: value,
	}
	tx.Statement.EntryList = append(etcdb.Statement.EntryList, e)
	tx.processors = append(etcdb.processors, put())

	tx.exec()

	return tx
}

func (etcdb *Etcdb) PutBatch(kvs map[string]string) *Etcdb {
	tx := etcdb.getInstance()

	for k, v := range kvs {
		e := Entry{
			Key:   k,
			Value: v,
		}
		tx.Statement.EntryList = append(tx.Statement.EntryList, e)
	}

	tx.processors = append(tx.processors, put())

	tx.exec()

	return tx
}

func (etcdb *Etcdb) Get(e *[]Entry) *Etcdb {
	tx := etcdb.getInstance()

	tx.processors = append(tx.processors, get())

	tx.exec()

	*e = tx.Statement.EntryList

	return tx
}

func (etcdb *Etcdb) Drop() *Etcdb {
	tx := etcdb.getInstance()

	tx.processors = append(tx.processors, drop())

	tx.exec()

	return tx
}
