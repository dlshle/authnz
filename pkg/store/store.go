package store

type PBEntityStore interface {
	Get(id string) (*PBEntity, error)
	TxGet(SQLTransactional, string) (*PBEntity, error)
	TxBulkGet(tx SQLTransactional, ids []string) ([]*PBEntity, error)
	Put(*PBEntity) (*PBEntity, error)
	TxPut(SQLTransactional, *PBEntity) (*PBEntity, error)
	Delete(id string) error
	TxDelete(SQLTransactional, string) error
	WithTx(cb func(SQLTransactional) error) error
}
