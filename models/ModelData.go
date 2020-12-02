package models

import (
	
	"database/sql"
	//"encoding/hex"
	"null"
	"encoding/json"
	//guuid "github.com/google/uuid"
	//"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)



type ModelGetData struct {
	DB *sql.DB
}

type Db_init struct {
	DB *sql.DB
}


func (model ModelGetData) GetSQLData(cmd string) (getStruct [] string, err error) {
    log.WithFields(log.Fields{"CMD": cmd,	}).Debug("func GetSQLData")
	
	row, err := model.DB.Query(cmd)
	if err != nil {				
		log.WithFields(log.Fields{"Error": err.Error(),	}).Error("func GetSQLData, query command")
		return  nil,err
	} else {
		
		colNames, err := row.Columns()
		
		if err != nil {
				log.WithFields(log.Fields{"Error": err.Error(),	}).Error("func GetSQLData, get colnames")
			} else {
		
		}
		
		readCols := make([]interface{}, len(colNames))
		writeCols := make([]null.String, len(colNames))
		var _isiStruct [] string

		for i, _ := range writeCols {
			readCols[i] = &writeCols[i]
		}
		
		for row.Next() {
			err := row.Scan(readCols...)			
			if err != nil {
				log.WithFields(log.Fields{"Error": err.Error(),	}).Error("func GetSQLData, scan next row")
			} else {						
				pagesJson, err := json.Marshal(writeCols)
				if err != nil {
					log.WithFields(log.Fields{"Error": err.Error(),	}).Error("func GetSQLData, marshal to JSON")
				}
				_isiStruct = append(_isiStruct,string(pagesJson))
			}
		}
		return _isiStruct, nil
	}
}








