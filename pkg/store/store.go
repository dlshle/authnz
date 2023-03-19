package store

type PBEntityStore interface {
	Get(id string) (*PBEntity, error)
	Put(*PBEntity) (*PBEntity, error)
	Delete(id string) error
}
