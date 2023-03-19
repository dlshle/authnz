package store

type PBEntity struct {
	ID      string `db:"id"`
	Payload []byte `db:"payload"`
}
