package data

import (
	"database/sql"
	"shorturl/pkg/constants"
	"shorturl/pkg/log"
)

type IUrlMapDataFactory interface {
	NewUrlMapData(isPublic bool) IUrlMapData
}

type UrlMapDataFactory struct {
	log       log.ILogger
	db        *sql.DB
	tableName string
}

func NewUrlMapDataFactory(log log.ILogger, db *sql.DB) IUrlMapDataFactory {
	return &UrlMapDataFactory{
		log: log,
		db:  db,
	}
}

func (f *UrlMapDataFactory) NewUrlMapData(isPublic bool) IUrlMapData {
	tableName := constants.TABLE_URL_MAP
	if !isPublic {
		tableName = constants.TABLE_URL_MAP_USER
	}
	return newUrlMapData(f.log, f.db, tableName)
}
