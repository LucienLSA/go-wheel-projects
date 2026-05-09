package session_test

import (
	databasesql "ormframework/save-query-orm"
	"ormframework/save-query-orm/session"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	engine, err := databasesql.NewEngine("sqlite3", ":memory:")
	if err != nil {
		t.Fatal("failed to create engine:", err)
	}
	s := engine.NewSession().Model(&User{})
	if err := s.DropTable(); err != nil {
		t.Fatal("failed to drop table:", err)
	}
	if err := s.CreateTable(); err != nil {
		t.Fatal("failed to create table:", err)
	}
	if _, err := s.Insert(user1, user2); err != nil {
		t.Fatal("failed to insert records:", err)
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}
