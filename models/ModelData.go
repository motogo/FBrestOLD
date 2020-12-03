package models

import (
	
	"database/sql"
	"null"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)



type ModelGetData struct {
	DB *sql.DB
}

type Db_init struct {
	DB *sql.DB
}


func (model ModelGetData) GetSQLData(cmd string) (getStruct [] string, err error) {
	const funcstr = "func GetSQLData"
    log.WithFields(log.Fields{"Getting data": cmd,	}).Debug(funcstr)
	
	row, err := model.DB.Query(cmd)
	if err != nil {				
		log.WithFields(log.Fields{"Error": err.Error(),	}).Error(funcstr+"->query command")
		return  nil,err
	} else {
		
		colNames, err := row.Columns()
		
		if err != nil {
				log.WithFields(log.Fields{"Error": err.Error(),	}).Error(funcstr+"->get colnames")
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
				log.WithFields(log.Fields{"Error": err.Error(),	}).Error(funcstr+"->scan next row")
			} else {						
				pagesJson, err := json.Marshal(writeCols)
				if err != nil {
					log.WithFields(log.Fields{"Error": err.Error(),	}).Error(funcstr+"->marshal to JSON")
				}
				_isiStruct = append(_isiStruct,string(pagesJson))
			}
		}
		log.WithFields(log.Fields{"Getting data done": cmd,	}).Debug(funcstr)
		return _isiStruct, nil
	}
}








