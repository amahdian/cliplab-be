package pg

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/amahdian/cliplab-be/pkg/logger"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type LogLevel string

const (
	Silent LogLevel = "silent"
	Error           = "error"
	Warn            = "warn"
	Info            = "info"
)

func OpenGormDb(dsn string, logLevel LogLevel) (*gorm.DB, error) {
	logger.Infof("trying to open connection to database %q", dsn)
	// create db object
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel(logLevel)),
	})
	if err != nil {
		return nil, err
	}

	// ping the db
	sqlDb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get underlying sql db from gorm")
	}
	err = sqlDb.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	logger.Info("successfully opened database connection")
	return db, nil
}

func EnsureDatabaseExists(dsn string) error {
	u, err := dburl.Parse(dsn)
	if err != nil {
		return fmt.Errorf("couldn't parse database address, %w", err)
	}
	dbname := strings.TrimPrefix(u.Path, "/")
	if dbname != "" {
		u.Path = ""
		db, err := gorm.Open(postgres.Open(u.String()), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})
		if err != nil {
			return errors.Wrap(err, "failed to open connection to the db")
		}

		err = db.Exec("CREATE DATABASE " + dbname).Error
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				return nil
			}
			return errors.Wrap(err, "failed to create database")
		}
	}
	return nil
}

func OpenGormTestDb(t *testing.T, dsn string, logLevel LogLevel) *gorm.DB {
	db, err := OpenGormDb(dsn, logLevel)
	if err != nil {
		t.Fatalf("could not connect to the test db: %v", err)
	}
	cleanup := attachDeleteCreatedEntitiesHook(db)

	t.Cleanup(func() {
		// t.Log("cleaning up and closing test db...")
		cleanup()

		// t.Log("closing the db connection...")
		sqlDb, err := db.DB()
		if err != nil {
			t.Fatalf("could close connection to test db: %v", err)
		}
		sqlDb.Close()
	})

	return db
}

func attachDeleteCreatedEntitiesHook(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyName string
		keyVal  interface{}
	}
	var entries []entity
	hookName := "test_cleanup_hook"

	db.Callback().Create().After("gorm:create").Register(hookName, func(db *gorm.DB) {
		if db.DryRun || db.Error != nil {
			return
		}

		tableName := db.Statement.Table

		tryExtractingCreatedEntry := func(statementVal reflect.Value) {
			// use the first available field if created entry does not have a primary key
			if db.Statement.Schema.PrioritizedPrimaryField == nil {
				for _, field := range db.Statement.Schema.Fields {
					if fieldVal, isZero := field.ValueOf(db.Statement.Context, statementVal); !isZero {
						fieldName := field.DBName
						entries = append(entries, entity{table: tableName, keyName: fieldName, keyVal: fieldVal})
						return
					}
				}
			}

			// use primary key if present
			if fieldVal, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.Context, statementVal); !isZero {
				fieldName := db.Statement.Schema.PrioritizedPrimaryField.DBName
				entries = append(entries, entity{table: tableName, keyName: fieldName, keyVal: fieldVal})
				return
			}
		}

		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				rv := db.Statement.ReflectValue.Index(i)
				if reflect.Indirect(rv).Kind() != reflect.Struct {
					break
				}
				tryExtractingCreatedEntry(rv)
			}
		case reflect.Struct:
			tryExtractingCreatedEntry(db.Statement.ReflectValue)
		}
	})

	return func() {
		// since I'm closing the connection, this is not needed.
		// remove our hook once we're done.
		// defer db.Callback().Create().Remove(hookName)

		// begin a transaction if we are not already in one
		_, inTransaction := db.Statement.ConnPool.(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}

		// delete the entries in the reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			// fmt.Printf("Deleting entity from '%s' table with key: %v\n", entry.table, entry.keyVal)
			tx.Table(entry.table).Where(entry.keyName+" = ?", entry.keyVal).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
}

func gormLogLevel(level LogLevel) gormLogger.LogLevel {
	switch level {
	case Silent:
		return gormLogger.Silent
	case Error:
		return gormLogger.Error
	case Warn:
		return gormLogger.Warn
	case Info:
		return gormLogger.Info
	}
	return gormLogger.Error
}
