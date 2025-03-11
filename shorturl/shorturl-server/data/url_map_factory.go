package data

import (
	"database/sql"
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
	tableName := "url_map"
	if !isPublic {
		tableName = "url_map_user"
	}
	return newUrlMapData(f.log, f.db, tableName)
}
