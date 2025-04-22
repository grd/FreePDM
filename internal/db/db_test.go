package db_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/grd/FreePDM/internal/db"
	"github.com/grd/FreePDM/internal/server"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	gormdb, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory SQLite DB: %v", err)
	}

	// automatic migration
	err = gormdb.AutoMigrate(&db.PdmMaterial{}, &db.PdmModel{}, &db.PdmUser{}, &db.PdmProject{}, &db.PdmItem{}) // etc.
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return gormdb
}

func createTestUser(gormdb *gorm.DB) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	user := db.PdmUser{
		UserName:     "testuser",
		PasswordHash: string(hashedPassword),
		EmailAddress: "test@example.com",
	}
	return gormdb.Create(&user).Error
}

func TestValidLogin(t *testing.T) {
	// Arrange
	gormdb := setupTestDB(t)
	repo := db.NewUserRepo(gormdb)

	server := server.NewServer(repo)

	// create test user
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	_ = gormdb.Create(&db.PdmUser{
		UserName:     "jdoe",
		PasswordHash: string(hashed),
	})

	form := url.Values{}
	form.Add("username", "jdoe")
	form.Add("password", "secret")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Act
	server.HandleLogin(w, req)

	// Assert
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusSeeOther {
		t.Errorf("expected redirect, got %d", res.StatusCode)
	}
}

// func TestLoginWithValidUser(t *testing.T) {
// 	db := setupTestDB(t)

// 	// voeg gebruiker toe, test login etc...
// }
