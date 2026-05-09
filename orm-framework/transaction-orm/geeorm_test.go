package databasesql_test

import (
	"errors"
	databasesql "ormframework/transaction-orm"
	"ormframework/transaction-orm/session"
	"testing"
)

func OpenDB(t *testing.T) *databasesql.Engine {
	t.Helper()
	engine, err := databasesql.NewEngine("sqlite3", "")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func Test_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, err := s.Insert(&User{"Tom", 18})
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_ = s.Model(&User{}).CreateTable()
		_, err := s.Insert(&User{"Tom", 18})
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
