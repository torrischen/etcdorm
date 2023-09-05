package eorm

type Entry struct {
	Key            string `json:"key"`
	Value          string `json:"value"`
	Version        int64  `json:"version"`
	CreateRevision int64  `json:"create_revision"`
	ModRevision    int64  `json:"mod_revision"`
	Lease          int64  `json:"lease"`
}
