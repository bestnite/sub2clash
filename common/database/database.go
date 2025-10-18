package database

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/bestnite/sub2clash/common"
	"github.com/bestnite/sub2clash/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func ConnectDB() (*Database, error) {
	path := filepath.Join("data", "sub2clash.db")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Discard,
	})
	if err != nil {
		return nil, common.NewDatabaseConnectError(err)
	}

	if err = db.AutoMigrate(&model.ShortLink{}); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) FindShortLinkByID(id string) (model.ShortLink, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return gorm.G[model.ShortLink](d.db).Where("id = ?", id).First(ctx)
}

func (d *Database) CreateShortLink(shortLink *model.ShortLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return gorm.G[model.ShortLink](d.db).Create(ctx, shortLink)
}

func (d *Database) UpdataShortLink(id string, name string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := gorm.G[model.ShortLink](d.db).Where("id = ?", id).Update(ctx, name, value)
	return err
}

func (d *Database) CheckShortLinkIDExists(id string) (bool, error) {
	_, err := d.FindShortLinkByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (d *Database) DeleteShortLink(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := gorm.G[model.ShortLink](d.db).Where("id = ?", id).Delete(ctx)
	return err
}
