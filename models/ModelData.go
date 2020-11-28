package models

import (
	_struct "fbrest/Dokumento/struct"
	"database/sql"
)

type ModelGetData struct {
	DB *sql.DB
}

func (model ModelGetData) GetDataTableStandort() (getStruct []_struct.StructData, err error) {
	row, err := model.DB.Query("select ID, BEZ, SCHLUESSEL FROM TSTANDORT")
	if err != nil {
		return  nil,err
	} else {
		var _isiStruct []_struct.StructData
		var data _struct.StructData
		for row.Next() {
			err2 := row.Scan(
				&data.ID,
				&data.BEZ,
				&data.SCHLUESSEL)
			if err2 != nil {
				return nil, err2
			} else {
				_data := _struct.StructData{
					ID:  data.ID,
					BEZ: data.BEZ,
					SCHLUESSEL:  data.SCHLUESSEL,
				}
				_isiStruct = append(_isiStruct, _data)

			}
		}

		return _isiStruct, nil
	}
}