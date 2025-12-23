package sql

// import (
// 	"context"
// 	"log/slog"
// 	"sync"

// 	"github.com/semanggilab/webcore-go/app/config"
// 	libdb "github.com/semanggilab/webcore-go/app/libs/db"
// )

// var (
// 	dbConns = map[string]libdb.Database{}
// 	lockDB  sync.Mutex
// )

// // func WrapMonggoClient(conn *mongoDB.Client) Database {
// // 	driver := "mongodb"
// // 	if _, ok := dbConns[driver]; !ok {
// // 		lockDB.Lock()
// // 		defer lockDB.Unlock()
// // 		dbConns[driver] = WrapMonggoConnection(config.DatabaseConfig{
// // 			Driver: "mongodb",
// // 		}, conn)
// // 	}
// // 	return dbConns[driver]
// // }

// func GetConnection(config *config.DatabaseConfig) (libdb.Database, error) {
// 	if _, ok := dbConns[config.Driver]; !ok {
// 		lockDB.Lock()
// 		defer lockDB.Unlock()
// 		db, err := loadDB(config)
// 		if err != nil {
// 			logger.Error("Database", "server", config.Name, "error", err)
// 			return nil, err
// 		}
// 		dbConns[config.Driver] = db
// 	}
// 	return dbConns[config.Driver], nil
// }

// func loadDB(config *config.DatabaseConfig) (libdb.Database, error) {
// 	var db libdb.Database
// 	var err error
// 	driver := config.Driver
// 	if driver == "mongodb" {
// 		db = NewMongoDB(*config)
// 	} else {
// 		db, err = NewSQLDatabase(*config)
// 	}

// 	return db, err
// }

// func CloseConnection(ctx context.Context) {
// 	for _, dbConn := range dbConns {
// 		if dbConn != nil {
// 			if err := dbConn.Close(ctx); err != nil {
// 				logger.Error("got an error while disconnecting database", "server", dbConn.GetName(), "error", err)
// 			}
// 		}
// 	}
// }
