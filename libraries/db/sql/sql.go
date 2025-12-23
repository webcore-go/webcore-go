package sql

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type SQLDB interface {
	loader.IDatabase
}

// SQLDatabase represents shared database connection
type SQLDatabase struct {
	Context context.Context
	Config  config.DatabaseConfig
	DB      *gorm.DB
	Dialect gorm.Dialector
}

func (d *SQLDatabase) SetDialect(dialect gorm.Dialector) {
	d.Dialect = dialect
}

// Install library
func (d *SQLDatabase) Install(args ...any) error {
	d.Context = args[0].(context.Context)
	d.Config = args[1].(config.DatabaseConfig)

	if d.Dialect == nil {
		return fmt.Errorf("Gorm Dialect is not set")
	}

	if d.Config.Driver != d.Dialect.Name() {
		return fmt.Errorf("Driver(%s) and Dialect(%s) does not match", d.Config.Driver, d.Dialect.Name())
	}

	return nil
}

// Connect establishes a database connection
func (d *SQLDatabase) Connect() error {
	gconf := gorm.Config{Logger: glogger.Default.LogMode(glogger.Info)}
	db, err := gorm.Open(d.Dialect, &gconf)
	// switch d.Config.Driver {
	// case "postgres":
	// 	db, err = gorm.Open(postgres.Open(dsn), &gconf)
	// case "mysql":
	// 	db, err = gorm.Open(mysql.Open(dsn), &gconf)
	// case "sqlite":
	// 	dsn := d.Config.Name
	// 	db, err = gorm.Open(sqlite.Open(dsn), &gconf)
	// default:
	// 	return fmt.Errorf("unsupported database driver: %s", d.Config.Driver)
	// }

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	d.DB = db

	return nil
}

// Close closes the database connection
func (d *SQLDatabase) Disconnect() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Connect establishes a database connection
func (d *SQLDatabase) Uninstall() error {
	// Connection is already established in NewSQLDatabase
	return nil
}

// Ping checks if the database connection is alive
func (d *SQLDatabase) Ping(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetConnection returns the underlying database connection
func (d *SQLDatabase) GetConnection() any {
	return d.DB
}

// GetDriver returns the database driver name
func (d *SQLDatabase) GetDriver() string {
	return d.Config.Driver
}

// GetName returns the database name
func (d *SQLDatabase) GetName() string {
	return d.Config.Name
}

// Count counts records in a table with optional filtering
func (d *SQLDatabase) Count(ctx context.Context, table string, filter loader.DbMap) (int64, error) {
	var count int64
	query := d.DB.Table(table)

	if filter != nil {
		query = query.Where(filter)
	}

	result := query.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// Find retrieves records from a table with optional filtering, sorting, and pagination
func (d *SQLDatabase) Find(ctx context.Context, table string, column []string, filter loader.DbMap, sort map[string]int, limit int64, skip int64) ([]loader.DbMap, error) {
	var results []loader.DbMap
	query := d.DB.Table(table)

	if len(column) > 0 {
		query = query.Select(column)
	}

	if filter != nil {
		query = query.Where(filter)
	}

	for field, direction := range sort {
		if direction > 0 {
			query = query.Order(fmt.Sprintf("%s ASC", field))
		} else {
			query = query.Order(fmt.Sprintf("%s DESC", field))
		}
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	if skip > 0 {
		query = query.Offset(int(skip))
	}

	result := query.Find(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	return results, nil
}

// FindOne retrieves a single record from a table
func (d *SQLDatabase) FindOne(ctx context.Context, result any, table string, column []string, filter loader.DbMap, sort map[string]int) error {
	query := d.DB.Table(table)

	if len(column) > 0 {
		query = query.Select(column)
	}

	if filter != nil {
		query = query.Where(filter)
	}

	for field, direction := range sort {
		if direction > 0 {
			query = query.Order(fmt.Sprintf("%s ASC", field))
		} else {
			query = query.Order(fmt.Sprintf("%s DESC", field))
		}
	}

	return query.First(result).Error
}

// InsertOne inserts a single record into a table
func (d *SQLDatabase) InsertOne(ctx context.Context, table string, data any) (any, error) {
	result := d.DB.Table(table).Create(data)
	if result.Error != nil {
		return nil, result.Error
	}

	// Return the inserted data with ID
	if result.RowsAffected > 0 {
		return data, nil
	}

	return nil, fmt.Errorf("failed to insert record")
}

// Update updates records in a table with optional filtering
func (d *SQLDatabase) Update(ctx context.Context, table string, filter loader.DbMap, data any) (int64, error) {
	query := d.DB.Table(table).Where(filter)
	result := query.Updates(data)
	return result.RowsAffected, result.Error
}

// UpdateOne updates a single record in a table
func (d *SQLDatabase) UpdateOne(ctx context.Context, table string, filter loader.DbMap, data any) (int64, error) {
	result := d.DB.Table(table).Where(filter).First(&loader.DbMap{}).Updates(data)
	return result.RowsAffected, result.Error
}

// Delete deletes records from a table with optional filtering
func (d *SQLDatabase) Delete(ctx context.Context, table string, filter loader.DbMap) (any, error) {
	var result *gorm.DB
	if filter != nil {
		result = d.DB.Table(table).Where(filter).Delete(&loader.DbMap{})
	} else {
		result = d.DB.Table(table).Delete(&loader.DbMap{})
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"deleted": result.RowsAffected,
	}, nil
}

// DeleteOne deletes a single record from a table
func (d *SQLDatabase) DeleteOne(ctx context.Context, table string, filter loader.DbMap) (any, error) {
	result := d.DB.Table(table).Where(filter).First(&loader.DbMap{}).Delete(&loader.DbMap{})
	if result.Error != nil {
		return nil, result.Error
	}

	return map[string]interface{}{
		"deleted": result.RowsAffected,
	}, nil
}

// BuildDSN builds a database connection string from configuration
func BuildDSN(config config.DatabaseConfig) string {
	var dsn string
	switch config.Driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			config.Host, config.User, config.Password, config.Name, config.Port, config.SSLMode)
		if len(config.Attributes) > 0 {
			queryParams := url.Values{}
			for key, value := range config.Attributes {
				queryParams.Add(key, value)
			}
			if len(queryParams) > 0 {
				dsn += "?" + queryParams.Encode()
			}
		}
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User, config.Password, config.Host, config.Port, config.Name)

		attributes := config.Attributes
		if _, exists := attributes["charset"]; !exists {
			attributes["charset"] = "utf8mb4"
		}
		if _, exists := attributes["parseTime"]; !exists {
			attributes["parseTime"] = "True"
		}
		if _, exists := attributes["loc"]; !exists {
			attributes["loc"] = "Local"
		}

		queryParams := []string{}
		for key, value := range attributes {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
		}
		if len(queryParams) > 0 {
			dsn += "?" + strings.Join(queryParams, "&")
		}
	case "sqlite":
		dsn = config.Name
		if len(config.Attributes) > 0 {
			queryParams := []string{}
			for key, value := range config.Attributes {
				queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
			}
			if len(queryParams) > 0 {
				dsn += "?" + strings.Join(queryParams, "&")
			}
		}
	}

	return dsn
}
