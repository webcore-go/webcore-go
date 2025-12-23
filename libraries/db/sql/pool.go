package sql

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/semanggilab/webcore-go/app/config"
// 	"github.com/semanggilab/webcore-go/app/helper"
// 	"github.com/semanggilab/webcore-go/app/logger"
// 	"gorm.io/gorm"
// )

// // DatabasePool manages database connection pools
// type DatabasePool struct {
// 	master   *SQLDatabase
// 	slaves   []*SQLDatabase
// 	config   config.DatabaseConfig
// 	strategy string // "round-robin", "random", "weight"
// }

// // NewDatabasePool creates a new database connection pool
// func NewDatabasePool(config config.DatabaseConfig) (*DatabasePool, error) {
// 	pool := &DatabasePool{
// 		config:   config,
// 		strategy: "round-robin",
// 	}

// 	// Connect to master database
// 	if err := pool.connectMaster(); err != nil {
// 		return nil, fmt.Errorf("failed to connect to master database: %v", err)
// 	}

// 	// Connect to slave databases if configured
// 	if len(config.SlaveHosts) > 0 {
// 		if err := pool.connectSlaves(); err != nil {
// 			logger.Warn("Failed to connect to some slave databases, continuing with master only")
// 		}
// 	}

// 	return pool, nil
// }

// // connectMaster connects to the master database
// func (p *DatabasePool) connectMaster() error {
// 	// dsn := BuildDSN(p.config)

// 	var err error
// 	var ok bool
// 	// p.master, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
// 	// 	Logger: glogger.Default.LogMode(glogger.Info),
// 	// })
// 	db, err := NewSQLDatabase(p.config)
// 	if err != nil {
// 		return err
// 	}

// 	p.master, ok = db.(*SQLDatabase)
// 	if !ok {
// 		return fmt.Errorf("failed to convert to SQLDatabase")
// 	}

// 	// Get underlying sql.DB for connection pool configuration
// 	sqlDB, err := p.master.DB.DB()
// 	if err != nil {
// 		return err
// 	}

// 	// Configure connection pool
// 	sqlDB.SetMaxIdleConns(p.config.MaxIdleConns)
// 	sqlDB.SetMaxOpenConns(p.config.MaxOpenConns)
// 	sqlDB.SetConnMaxLifetime(time.Duration(p.config.ConnMaxLifetime) * time.Second)

// 	logger.Info("Successfully connected to master database")
// 	return nil
// }

// // connectSlaves connects to slave databases
// func (p *DatabasePool) connectSlaves() error {
// 	p.slaves = make([]*SQLDatabase, 0, len(p.config.SlaveHosts))

// 	for _, slaveConfig := range p.config.SlaveHosts {
// 		slaveConfig.Driver = p.config.Driver // override driver from master config

// 		// dsn := BuildDSN(slaveConfig)
// 		// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		// 	Logger: glogger.Default.LogMode(glogger.Info),
// 		// })
// 		db, err := NewSQLDatabase(p.config)

// 		if err != nil {
// 			logger.Warn(fmt.Sprintf("Failed to connect to slave database %s: %v", slaveConfig.Host, err))
// 			continue
// 		}

// 		db2, ok := db.(*SQLDatabase)
// 		if !ok {
// 			logger.Warn("failed to convert to SQLDatabase")
// 			continue
// 		}

// 		// Configure connection pool for slave
// 		sqlDB, err := db2.DB.DB()
// 		if err != nil {
// 			logger.Warn(fmt.Sprintf("Failed to get sql.DB for slave %s: %v", slaveConfig.Host, err))
// 			continue
// 		}

// 		sqlDB.SetMaxIdleConns(slaveConfig.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(slaveConfig.MaxOpenConns)
// 		sqlDB.SetConnMaxLifetime(time.Duration(slaveConfig.ConnMaxLifetime) * time.Second)

// 		p.slaves = append(p.slaves, db2)
// 		logger.Info(fmt.Sprintf("Successfully connected to slave database: %s", slaveConfig.Host))
// 	}

// 	return nil
// }

// // GetDB returns a database connection based on the configured strategy
// func (p *DatabasePool) GetDB() *gorm.DB {
// 	if len(p.slaves) == 0 {
// 		return p.master.DB
// 	}

// 	switch p.strategy {
// 	case "round-robin":
// 		return p.getRoundRobinSlave()
// 	case "random":
// 		return p.getRandomSlave()
// 	case "weight":
// 		return p.getWeightedSlave()
// 	default:
// 		return p.getRoundRobinSlave()
// 	}
// }

// // getRoundRobinSlave returns a slave database connection using round-robin strategy
// func (p *DatabasePool) getRoundRobinSlave() *gorm.DB {
// 	// This is a simplified implementation
// 	// In a real scenario, you would maintain a counter for round-robin
// 	return p.slaves[0].DB
// }

// // getRandomSlave returns a random slave database connection
// func (p *DatabasePool) getRandomSlave() *gorm.DB {
// 	// This is a simplified implementation
// 	// In a real scenario, you would use a random number generator
// 	return p.slaves[0].DB
// }

// // getWeightedSlave returns a slave database connection based on weight
// func (p *DatabasePool) getWeightedSlave() *gorm.DB {
// 	// This is a simplified implementation
// 	// In a real scenario, you would implement weighted selection
// 	return p.slaves[0].DB
// }

// // GetMaster returns the master database connection
// func (p *DatabasePool) GetMaster() *gorm.DB {
// 	return p.master.DB
// }

// // Close closes all database connections
// func (p *DatabasePool) Disconnect() error {
// 	var errors []error

// 	// Close master connection
// 	if p.master != nil {
// 		sqlDB, err := p.master.DB.DB()
// 		if err == nil {
// 			if err := sqlDB.Disconnect(); err != nil {
// 				errors = append(errors, fmt.Errorf("failed to close master connection: %v", err))
// 			}
// 		}
// 	}

// 	// Close slave connections
// 	for i, slave := range p.slaves {
// 		if slave != nil {
// 			sqlDB, err := slave.DB.DB()
// 			if err == nil {
// 				if err := sqlDB.Disconnect(); err != nil {
// 					errors = append(errors, fmt.Errorf("failed to close slave connection %d: %v", i, err))
// 				}
// 			}
// 		}
// 	}

// 	if len(errors) > 0 {
// 		return fmt.Errorf("encountered %d errors while closing database connections", len(errors))
// 	}

// 	return nil
// }

// // Health checks the health of all database connections
// func (p *DatabasePool) Health() map[string]any {
// 	health := make(map[string]any)

// 	// Check master
// 	if p.master != nil {
// 		sqlDB, err := p.master.DB.DB()
// 		if err == nil {
// 			err = sqlDB.Ping()
// 		}
// 		health["master"] = map[string]any{
// 			"status":   "healthy",
// 			"error":    helper.ErrToString(err),
// 			"max_idle": sqlDB.Stats().MaxOpenConnections, // Note: MaxIdleConnections is not available in sql.DBStats
// 			"max_open": sqlDB.Stats().MaxOpenConnections,
// 		}
// 	}

// 	// Check slaves
// 	slaveHealth := make([]map[string]any, 0, len(p.slaves))
// 	for i, slave := range p.slaves {
// 		if slave != nil {
// 			sqlDB, err := slave.DB.DB()
// 			if err == nil {
// 				err = sqlDB.Ping()
// 			}
// 			slaveHealth = append(slaveHealth, map[string]any{
// 				"id":       i,
// 				"status":   "healthy",
// 				"error":    helper.ErrToString(err),
// 				"max_idle": sqlDB.Stats().MaxOpenConnections, // Note: MaxIdleConnections is not available in sql.DBStats
// 				"max_open": sqlDB.Stats().MaxOpenConnections,
// 			})
// 		}
// 	}

// 	health["slaves"] = slaveHealth
// 	health["total_slaves"] = len(slaveHealth)

// 	return health
// }

// // SetStrategy sets the read strategy for slave selection
// func (p *DatabasePool) SetStrategy(strategy string) {
// 	p.strategy = strategy
// }

// // Transaction executes a transaction on the master database
// func (p *DatabasePool) Transaction(fc func(tx *gorm.DB) error) error {
// 	return p.master.DB.Transaction(fc)
// }

// // WithContext executes a query with the given context
// func (p *DatabasePool) WithContext(ctx context.Context) *gorm.DB {
// 	return p.master.DB.WithContext(ctx)
// }

// // Model returns a new query for the given model
// func (p *DatabasePool) Model(value any) *gorm.DB {
// 	return p.master.DB.Model(value)
// }

// // First returns the first record that matches the query
// func (p *DatabasePool) First(dest any, conds ...any) *gorm.DB {
// 	return p.master.DB.First(dest, conds...)
// }

// // Find returns all records that match the query
// func (p *DatabasePool) Find(dest any, conds ...any) *gorm.DB {
// 	return p.master.DB.Find(dest, conds...)
// }

// // Create inserts a new record
// func (p *DatabasePool) Create(value any) *gorm.DB {
// 	return p.master.DB.Create(value)
// }

// // Save updates a record or inserts a new one
// func (p *DatabasePool) Save(value any) *gorm.DB {
// 	return p.master.DB.Save(value)
// }

// // Delete deletes a record
// func (p *DatabasePool) Delete(value any, conds ...any) *gorm.DB {
// 	return p.master.DB.Delete(value, conds...)
// }

// // Updates updates a record with the given values
// func (p *DatabasePool) Updates(value any) *gorm.DB {
// 	return p.master.DB.Updates(value)
// }

// // Update updates a single column
// func (p *DatabasePool) Update(column string, value any) *gorm.DB {
// 	return p.master.DB.Update(column, value)
// }
