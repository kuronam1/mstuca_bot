package processor

import (
	"database/sql"
	"mstuca_schedule/internal/models"
)

type Processor interface {
	SaveProfile(*models.User) error
	EditProfile(*models.User) error
	GetProfile(id int64) (*models.User, error)
	IsExist(id int64) bool
}

type processor struct {
	sqliteClient *sql.DB
}

func New() (Processor, error) {

	return nil, nil

	db, err := sql.Open("sqlite3", "mstuca_users")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &processor{
		sqliteClient: db,
	}, nil
}

func (p *processor) SaveProfile(*models.User) error {
	return nil
}

func (p *processor) EditProfile(*models.User) error {
	return nil
}

func (p *processor) GetProfile(id int64) (*models.User, error) {
	return nil, nil
}

func (p *processor) IsExist(id int64) bool {
	return false
}
