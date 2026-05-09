package databasesql_test

import (
	databasesql "ormframework/save-query-orm"
	"ormframework/save-query-orm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	engine, err := databasesql.NewEngine("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	_ = s.DropTable()
	if err := s.CreateTable(); err != nil {
		t.Fatal(err)
	}
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
func TestSession_Model(t *testing.T) {
	engine, err := databasesql.NewEngine("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()
	s := engine.NewSession().Model(&User{})
	table := s.RefTable()
	s.Model(&session.Session{})
	if table.Name != "User" || s.RefTable().Name != "Session" {
		t.Fatal("Failed to change model")
	}
}
