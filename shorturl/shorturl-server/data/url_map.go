package data

import (
	"database/sql"
	"fmt"
	"shorturl/pkg/log"
	"shorturl/pkg/zerror"
)

type UrlMapEntity struct {
	ID          int64
	UserID      int64
	ShortKey    string
	OriginalUrl string
	Times       int
	CreatedAt   int64
	UpdatedAt   int64
}

type IUrlMapData interface {
	GenerateID(userID, now int64) (int64, error)
	Update(e UrlMapEntity) error
	GetByID(id int64) (UrlMapEntity, error)
	GetByOriginal(originalUrl string) (UrlMapEntity, error)
	IncrementTimes(id int64, incrementTimes int, now int64) error
}

type urlMapData struct {
	log       log.ILogger
	db        *sql.DB
	tableName string
}

func newUrlMapData(log log.ILogger, db *sql.DB, tableName string) IUrlMapData {
	return &urlMapData{
		log:       log,
		db:        db,
		tableName: tableName,
	}
}

func (d *urlMapData) GenerateID(userID, now int64) (int64, error) {
	if userID != 0 {
		sqlStr := fmt.Sprintf("insert into %s (user_id,create_at,update_at) values (?,?,?)", d.tableName)
		res, err := d.db.Exec(sqlStr, userID, now, now)
		if err != nil {
			d.log.Error(zerror.NewByErr(err))
			return 0, err
		}
		return res.LastInsertId()
	} else {
		sqlStr := fmt.Sprintf("insert into %s (create_at,update_at) values (?,?)", d.tableName)
		res, err := d.db.Exec(sqlStr, now, now)
		if err != nil {
			d.log.Error(zerror.NewByErr(err))
			return 0, err
		}
		return res.LastInsertId()
	}

}
func (d *urlMapData) Update(e UrlMapEntity) error {
	sqlStr := fmt.Sprintf("update %s set short_key=?,original_url=?,update_at=? where id=?", d.tableName)
	_, err := d.db.Exec(sqlStr, e.ShortKey, e.OriginalUrl, e.UpdatedAt, e.ID)
	if err != nil {
		d.log.Error(zerror.NewByErr(err))
	}
	return err
}
func (d *urlMapData) GetByID(id int64) (UrlMapEntity, error) {
	sqlStr := fmt.Sprintf("select original_url from %s where id=?", d.tableName)
	row := d.db.QueryRow(sqlStr, id)
	entity := UrlMapEntity{}
	var originalUrl sql.NullString
	err := row.Scan(&originalUrl)
	if err != nil && err != sql.ErrNoRows {
		d.log.Error(zerror.NewByErr(err))
		return entity, err
	}
	if originalUrl.Valid {
		entity.OriginalUrl = originalUrl.String
	}
	return entity, nil
}
func (d *urlMapData) GetByOriginal(originalUrl string) (UrlMapEntity, error) {
	sqlStr := fmt.Sprintf("select id, short_key from %s where original_url=?", d.tableName)
	row := d.db.QueryRow(sqlStr, originalUrl)
	entity := UrlMapEntity{}
	var shortKey sql.NullString
	err := row.Scan(&entity.ID, &shortKey)
	if err != nil && err != sql.ErrNoRows {
		d.log.Error(zerror.NewByErr(err))
		return entity, err
	}
	if shortKey.Valid {
		entity.ShortKey = shortKey.String
	}
	return entity, nil
}
func (d *urlMapData) IncrementTimes(id int64, incrementTimes int, now int64) error {
	sqlStr := fmt.Sprintf("update %s set times = times+?,update_at=? where id=?", d.tableName)
	_, err := d.db.Exec(sqlStr, incrementTimes, now, id)
	if err != nil {
		d.log.Error(zerror.NewByErr(err))
		return err
	}
	return nil

}
