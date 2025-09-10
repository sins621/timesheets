package database

import (
	"fmt"
	"os"
	"path/filepath"

	"ts_mcp/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase(db *gorm.DB) *GormDatabase {
	return &GormDatabase{db: db}
}

func (g *GormDatabase) CreateUser(user *models.User) (*models.User, error) {
	err := g.db.Create(user).Error

	return user, err
}

func (g *GormDatabase) SelectUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := g.db.Where("email = ?", email).First(&user).Error

	return &user, err
}

func (g *GormDatabase) UpdateUser(user *models.User) (*models.User, error) {
	err := g.db.Save(&user).Error

	return user, err
}

func InitializeGormDB() *gorm.DB {
	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v\n", err))
	}
	exeDir := filepath.Dir(exePath)
	dbPath := filepath.Join(exeDir, "timesheets.db")

	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Database failed to initialize")
	}

	err = gormDB.AutoMigrate(&models.User{})
	if err != nil {
		panic("Error running migration")
	}

	return gormDB
}
