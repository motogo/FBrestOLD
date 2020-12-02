//get AVG weight
func (model Db_init) Avg_weight() (OutputData []entities.AvgWeight, err error) {
	var ResultData []entities.AvgWeight
	var CData entities.AvgWeight
	queryExecute, errQueryExecute := model.DB.Query(`SELECT GUELTIG as gueltig_max,GUELTIG as gueltig_min, GUELTIG as gueltig_avg FROM TSTANDORT stand where stand.gueltig = "1"`)
	if errQueryExecute != nil {
		log.Println(errQueryExecute.Error())
	}

	for queryExecute.Next() {
		checkData := queryExecute.Scan(&CData.Max, &CData.Min, &CData.Perbedaan)
		if checkData != nil {
			log.Println(checkData.Error())
		}
		_Data := entities.AvgWeight{
			Max:       CData.Max,
			Min:       CData.Min,
			Perbedaan: CData.Perbedaan,
		}
		ResultData = append(ResultData, _Data)

	}
	return ResultData, nil
}

//get all data weight
func (model Db_init) Vweight() (OutputData []entities.Weight, err error) {
	var ResultData []entities.Weight
	var CData entities.Weight

	queryExecute, errQueryExecute := model.DB.Query(
		`select b.id,
  					  b.tanggal,
					  b.berat_max,
					  b.berat_min,
					  (b.berat_max-b.berat_min) as perbedaan
					  from berat b WHERE b.status_data ="AKTIF"`)

	if errQueryExecute != nil {
		log.Println(errQueryExecute.Error())
	}

	for queryExecute.Next() {
		checkData := queryExecute.Scan(&CData.WeightID, &CData.WeightDate, &CData.WeightMax, &CData.WeightMin, &CData.WeightDifference)
		if checkData != nil {
			log.Println(checkData.Error())
		}
		_Data := entities.Weight{
			WeightID:         CData.WeightID,
			WeightDate:       CData.WeightDate,
			WeightMax:        CData.WeightMax,
			WeightMin:        CData.WeightMin,
			WeightDifference: CData.WeightDifference,
		}
		ResultData = append(ResultData, _Data)

	}
	return ResultData, nil
}

//get AVG By ID weight
func (model Db_init) Avg_weightByID(id string) (OutputData []entities.AvgWeight, err error) {
	var ResultData []entities.AvgWeight
	var CData entities.AvgWeight
	queryExecute, errQueryExecute := model.DB.Query(`SELECT AVG(b.berat_max) as beratmax,AVG(b.berat_min) as beratmin, AVG(b.berat_max-b.berat_min) as perbedaann FROM berat b where b.status_data="AKTIF" AND id = ?`, id)
	if errQueryExecute != nil {
		log.Println(errQueryExecute.Error())
	}

	for queryExecute.Next() {
		checkData := queryExecute.Scan(&CData.Max, &CData.Min, &CData.Perbedaan)
		if checkData != nil {
			log.Println(checkData.Error())
		}
		_Data := entities.AvgWeight{
			Max:       CData.Max,
			Min:       CData.Min,
			Perbedaan: CData.Perbedaan,
		}
		ResultData = append(ResultData, _Data)

	}
	return ResultData, nil
}

//get all By ID data weight
func (model Db_init) VweightByID(id string) (OutputData []entities.Weight, err error) {
	var ResultData []entities.Weight
	var CData entities.Weight

	queryExecute, errQueryExecute := model.DB.Query(
		`select b.id,
  					  b.tanggal,
					  b.berat_max,
					  b.berat_min,
					  (b.berat_max-b.berat_min) as perbedaan
					  from berat b WHERE b.status_data ="AKTIF" AND id = ?`, id)

	if errQueryExecute != nil {
		log.Println(errQueryExecute.Error())
	}

	for queryExecute.Next() {
		checkData := queryExecute.Scan(&CData.WeightID, &CData.WeightDate, &CData.WeightMax, &CData.WeightMin, &CData.WeightDifference)
		if checkData != nil {
			log.Println(checkData.Error())
		}
		_Data := entities.Weight{
			WeightID:         CData.WeightID,
			WeightDate:       CData.WeightDate,
			WeightMax:        CData.WeightMax,
			WeightMin:        CData.WeightMin,
			WeightDifference: CData.WeightDifference,
		}
		ResultData = append(ResultData, _Data)

	}
	return ResultData, nil
}

//insert data weight
func (model Db_init) InsertWeight(OutputData *entities.Weight) (err error) {
	var dateData = OutputData.WeightDate
	var WeightMax = OutputData.WeightMax
	var WeightMin = OutputData.WeightMin
	var Count_data int

	CheckDataDate, errDataDate := config.Connection()
	CheckDataDate.QueryRow(`select COUNT(*) from berat where status_data ="AKTIF" AND tanggal =?`, dateData).Scan(&Count_data)
	if errDataDate != nil {
		log.Println(errDataDate)
	}
	if Count_data == 0 {
		insertData, errInsertData := model.DB.Exec(`insert into berat(tanggal,berat_max,berat_min,status_data) VALUES(?,?,?,'AKTIF') `, dateData, WeightMax, WeightMin)
		if errInsertData != nil {
			return errors.New("Server Disconnect")
		} else {
			insertData.LastInsertId()
			return nil
		}

	}
	return nil
}

//update data weight
func (model Db_init) UpdateWeight(OutputData *entities.Weight) (err error) {
	var idWeight = OutputData.WeightID
	var dateData = OutputData.WeightDate
	var WeightMax = OutputData.WeightMax
	var WeightMin = OutputData.WeightMin

		insertData, errInsertData := model.DB.Exec(`update berat set tanggal = ?, berat_max = ?,berat_min =?,status_data ='AKTIF' where id = ? `, dateData, WeightMax, WeightMin, idWeight)
		if errInsertData != nil {
			return errors.New("Server Disconnect")
		} else {
			insertData.RowsAffected()
			return nil
		}


	return nil
}

//delete data weight
func (model Db_init) DeleteWeight(id string) string {
	var responseData string
	insertData, errInsertData := model.DB.Exec(`update berat set status_data = 'NONAKTIF' where id = ? `, id)
	if errInsertData != nil {
		responseData = "gagal"
	} else {
		insertData.RowsAffected()
		responseData = "sukses"
	}
	return responseData
}