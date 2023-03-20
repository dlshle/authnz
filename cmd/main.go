package main

import (
	"context"
	"fmt"

	"github.com/dlshle/authnz/internal/config"
	"github.com/dlshle/authnz/internal/contract"
	"github.com/dlshle/authnz/internal/group"
	"github.com/dlshle/authnz/internal/migration"
	"github.com/dlshle/authnz/internal/policy"
	"github.com/dlshle/authnz/internal/server"
	"github.com/dlshle/authnz/internal/subject"
	pb "github.com/dlshle/authnz/proto"
	"github.com/dlshle/gommon/logging"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.Load("./etc/config.yaml")
	if err != nil {
		panic(err)
	}
	grpcServer, err := initGrpcServer(cfg)
	if err != nil {
		panic(err)
	}
	err = server.StartServer(cfg.Server, grpcServer)
	if err != nil {
		panic(err)
	}
}

func initGrpcServer(config config.Config) (pb.AuthNZServer, error) {
	dbConfig := config.Database
	db_connect := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Pass, dbConfig.DBName)
	logging.GlobalLogger.Infof(context.Background(), "db connection info: %s", db_connect)
	db, err := sqlx.Open("postgres", db_connect)
	if err != nil {
		return nil, err
	}
	err = execMigrationScript(db)
	if err != nil {
		return nil, err
	}
	groupSQLStore := group.NewSQLStore(db)
	policySQLStore := policy.NewSQLStore(db)
	subjectSQLStore := subject.NewSQLStore(db)
	contractSQLStore := contract.NewContractStore(db)

	groupHandler := group.NewHandler(groupSQLStore)
	policyHandler := policy.NewHandler(policySQLStore)
	subjectHandler := subject.NewHandler(subjectSQLStore, contractSQLStore, groupSQLStore)
	contractHandler := contract.NewHandler(contractSQLStore)

	grpcServer := server.NewGRPCServer(subjectHandler, groupHandler, policyHandler, contractHandler)
	return grpcServer, nil
}

func execMigrationScript(db *sqlx.DB) error {
	return migration.ExecMigration(db)
}
