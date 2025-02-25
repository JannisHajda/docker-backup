package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

type Client struct {
	*gorm.DB
}

func NewClient() (*Client, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Project{}, &Backup{})
	if err != nil {
		return nil, err
	}

	return &Client{db}, nil
}

func (c *Client) GetProject(name string) (*Project, error) {
	var p Project
	err := c.Where("name = ?", name).First(&p).Error
	return &p, err
}

func (c *Client) CreateProject(p *Project) error {
	return c.Create(p).Error
}

func (c *Client) CreateBackup(b *Backup) error {
	return c.Create(b).Error
}

func (c *Client) GetBackups(projectID uint) ([]Backup, error) {
	var backups []Backup
	err := c.Where("project_id = ?", projectID).Find(&backups).Error
	return backups, err
}

func (c *Client) GetBackup(projectID, backupID uint) (*Backup, error) {
	var b Backup
	err := c.Where("project_id = ? AND id = ?", projectID, backupID).First(&b).Error
	return &b, err
}
