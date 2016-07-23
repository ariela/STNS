package backend

import (
	"fmt"

	"github.com/STNS/STNS/stns"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const UNKNOWN_DATABASE = 1049

type Mysql struct {
	config *stns.Config
}

type Users struct {
	Name      string `gorm:"primary_key"`
	Password  string `gorm:"size:1024"`
	HashType  string
	GroupId   int
	Directory string
	Shell     string
	Gecos     string
}
type Sudoers struct {
	*Users
}

func (m *Mysql) Migrate() error {
	db, err := gorm.Open("mysql", m.connectInfo("stns"))
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if ok {
			if me.Number == UNKNOWN_DATABASE {
				err := m.createDatabase()
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	if err := db.AutoMigrate(&Users{}).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&Sudoers{}).Error; err != nil {
		return err
	}
	return nil
}

func (m *Mysql) connectInfo(db string) string {
	return fmt.Sprintf("%s:%s@tcp([%s]:%s)/%s",
		m.config.Backend.User,
		m.config.Backend.Password,
		m.config.Backend.Host,
		m.config.Backend.Port,
		db,
	)
}

func (m *Mysql) createDatabase() error {
	db, err := gorm.Open("mysql", m.connectInfo(""))
	if err != nil {
		return err
	}

	if err := db.Exec("CREATE DATABASE stns").Error; err != nil {
		return err
	}
	return nil
}
