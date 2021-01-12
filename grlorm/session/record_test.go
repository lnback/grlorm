package session

import (
	"database/sql"
	"grlorm/dialect"
	"os"
	"testing"
)

type User struct {
	Name string `grlorm:"PRIMARY KEY"`
	Age int
}
var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
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
var (
	TestDB * sql.DB
	TestDia, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M)  {
	TestDB,_ = sql.Open("mysql","root:123456@tcp(192.168.33.30:3306)/test?charset=utf8&parseTime=True&loc=Local")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)

}
func NewSession() *Session {
	return New(TestDB,TestDia)
}