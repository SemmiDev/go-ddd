package persistence

import (
	"DDD/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveUser_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = entity.User{}
	user.Email = "sammidev@gmail.com"
	user.FirstName = "Sammi"
	user.LastName = "Dev"
	user.Password = "sammidev"

	repo := NewUserRepository(conn)

	u, saveErr := repo.SaveUser(&user)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, u.Email, "sammidev@gmail.com")
	assert.EqualValues(t, u.FirstName, "Sammi")
	assert.EqualValues(t, u.LastName, "Dev")
	assert.NotEqual(t, u.Password, "password")
}

func TestSaveUser_Failure(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	_, err = seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	var user = entity.User{}
	user.Email = "sammidev@gmail.com"
	user.FirstName = "Sammi"
	user.LastName = "Dev"
	user.Password = "sammidev"

	repo := NewUserRepository(conn)
	u, saveErr := repo.SaveUser(&user)

	dbMsg := map[string]string{
		"email_taken": "email already taken",
	}

	assert.Nil(t, u)
	assert.EqualValues(t, dbMsg, saveErr)
}

func TestGetUser_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	user, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewUserRepository(conn)
	u, getErr := repo.GetUser(user.ID)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, "sammidev@gmail.com")
	assert.EqualValues(t, u.FirstName, "Sammi")
	assert.EqualValues(t, u.LastName, "Dev")
}

func TestGetUsers_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	_, err = seedUsers(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewUserRepository(conn)
	users, getErr := repo.GetUsers()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	u, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = &entity.User{
		Email:    "sammidev@gmail.com",
		Password: "sammidev",
	}
	repo := NewUserRepository(conn)
	u, getErr := repo.GetUserByEmailAndPassword(user)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, user.Email)
	assert.NotEqual(t, u.Password, user.Password)
}