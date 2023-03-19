package policy

import (
	"fmt"

	"github.com/dlshle/authnz/pkg/store"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/utils"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store interface {
	Get(id string) (*pb.Policy, error)
	Delete(id string) error
	Put(policy *pb.Policy) (*pb.Policy, error)
}

type SQLPolicyStore struct {
	pbEntityStore store.PBEntityStore
}

func NewSQLStore(db *sqlx.DB) Store {
	return &SQLPolicyStore{pbEntityStore: store.NewSQLPBEntityStore(db, "policies")}
}

func (s *SQLPolicyStore) Get(id string) (*pb.Policy, error) {
	pbEntity, err := s.pbEntityStore.Get(id)
	if err != nil {
		return nil, err
	}
	policy := &pb.Policy{}
	err = proto.Unmarshal(pbEntity.Payload, policy)
	return policy, err
}

func (s *SQLPolicyStore) Put(policy *pb.Policy) (ret *pb.Policy, err error) {
	var (
		payload  []byte
		pbEntity *store.PBEntity
	)
	err = utils.ProcessWithErrors(func() error {
		if policy.Id == "" {
			newID, err := uuid.NewV4()
			if err != nil {
				return err
			}
			policy.Id = newID.String()
		}
		return nil
	}, func() error {
		payload, err = proto.Marshal(policy)
		return err
	}, func() error {
		pbEntity, err = s.pbEntityStore.Put(&store.PBEntity{ID: policy.Id, Payload: payload})
		return err
	}, func() error {
		policy.Id = pbEntity.ID
		return nil
	})
	return policy, err
}

func (s *SQLPolicyStore) Delete(id string) error {
	return s.pbEntityStore.Delete(id)
}

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

var schema = `
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

func test() {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		panic(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec("INSERT INTO place (country, city, telcode) VALUES ($1, $2, $3)", "United States", "New York", "1")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Hong Kong", "852")
	tx.MustExec("INSERT INTO place (country, telcode) VALUES ($1, $2)", "Singapore", "65")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	// Query the database, storing results in a []Person (wrapped in []interface{})
	people := []Person{}
	db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
	jason, john := people[0], people[1]
	fmt.Println(jason, john)
}
