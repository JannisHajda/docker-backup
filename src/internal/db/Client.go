package db

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

var (
	dbPath = "test.db"
)

var db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

type Project struct {
	gorm.Model
	Name    string `gorm:"unique"`
	Backups []Backup
}

type Backup struct {
	gorm.Model
	ProjectID uint
	Project   Project
}

func init() {
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	db.AutoMigrate(&Project{}, &Backup{})
}

func NewClient() *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (c *Client) CreateProject(name string) (*Project, error) {
	project := Project{
		Name: name,
	}

	result := c.db.Create(&project)
	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (c *Client) GetProject(name string) (*Project, error) {
	var project Project
	result := c.db.Where("name = ?", name).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (p *Project) CreateBackup() (*Backup, error) {
	backup := Backup{
		ProjectID: p.ID,
	}

	result := db.Create(&backup)
	if result.Error != nil {
		return nil, result.Error
	}

	return &backup, nil
}

func (p *Project) GetBackup(name string) (*Backup, error) {
	var backup Backup
	result := db.Where("project_id = ? AND name = ?", p.ID, name).First(&backup)
	if result.Error != nil {
		return nil, result.Error
	}

	return &backup, nil
}

func (p *Project) GetBackups() ([]Backup, error) {
	var backups []Backup
	result := db.Where("project_id = ?", p.ID).Find(&backups)
	if result.Error != nil {
		return nil, result.Error
	}

	return backups, nil
}
