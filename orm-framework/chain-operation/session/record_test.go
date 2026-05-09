package session_test

import (
	databasesql "ormframework/chain-operation"
	"ormframework/chain-operation/session"
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

func TestSession_Limit(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	err := s.Limit(1).Find(&users)
	if err != nil || len(users) != 1 {
		t.Fatal("failed to query with limit condition")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 30)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 30 {
		t.Fatal("failed to update")
	}
}

func TestSession_DeleteAndCount(t *testing.T) {
	s := testRecordInit(t)
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	count, _ := s.Where("Name != ?", "Tom").Count()

	if affected != 1 || count != 1 {
		t.Fatal("failed to delete or count")
	}
}
