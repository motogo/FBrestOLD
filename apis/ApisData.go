package apis

import (
	"fbrest/Base/config"
	"fbrest/Base/models"	
	"strconv"	
	_struct "fbrest/Base/struct"
	_functions "fbrest/Base/functions"
	_sessions "fbrest/Base/sessions"
	"net/http"
	log "github.com/sirupsen/logrus"	
	bguuid "github.com/kjk/betterguid"
)

func TestDBOpenClose(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Info("func TestDBOpenClose")
	var key = _functions.GetPathSliceFromURL(r,0)	
	var kv  = _functions.KeyValid(response, key)
	if(!kv.Valid) {
		return
	}
	var err = config.TestConnLocation(kv.Value)
	var Response _struct.ResponseData
		
	if err != nil {
			
		Response.Status = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data = kv.Key			
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("Database not available")
		_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		Response.Status = http.StatusInternalServerError
		Response.Message = "Database open/close successfully"
		Response.Data = kv.Key
		log.Info("Database opend/closed successfully, ping sent")
			
		_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	}
}



func TestSQLResponse(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func TestSQLResponse")
	
	var Response _struct.ResponseData
	var entitiesData _struct.SQLAttributes
	
	_functions.GetSQLParamsFromURL(r , &entitiesData)				
	_functions.GetParamsFromBODY(r , &entitiesData)	
	_functions.OutParameters(entitiesData) 	
	
	Response.Status = http.StatusOK
	Response.Message = "Test SQL response"
	Response.Data = entitiesData.Cmd
	_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	
}

func TestTABLEResponse(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func TestSQLResponse")
	
	var Response _struct.ResponseData
	var entitiesData _struct.GetTABLEAttributes
	
	_functions.GetTableParamsFromURL(r , &entitiesData)				
	//_functions.GetParamsFromBODY(r , &entitiesData)	
	//_functions.OutParameters(entitiesData) 	
	
	Response.Status = http.StatusOK
	Response.Message = "Test table response"
	Response.Data = entitiesData
	_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	
}

func GetSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var dbData _struct.DatabaseAttributes

	dbData.Location = "localhost"	
	dbData.Port = 3050
	dbData.Password = "masterkey"
	dbData.User = "SYSDBA"

	
	var Response _struct.ResponseData
		
	_functions.GetSessionParamsFromURL(r , &dbData)		


	id := bguuid.New()				
	var connstr = string(dbData.User+":"+dbData.Password+"@"+dbData.Location+":"+strconv.Itoa(dbData.Port)+"/"+dbData.Database)			
	var rep = _sessions.Repository() 
	var itm = rep.Add(string(id),connstr)
	Response.Status = http.StatusOK
	Response.Message = "Created UUID, duration "+ itm.Duration.String()
	Response.Data =  string(id)
	_functions.RestponWithJson(response, http.StatusOK, Response)
	
}
func GetHelp(response http.ResponseWriter, r *http.Request) {

	_functions.RestponWithText(response, http.StatusOK)
	
}
func DeleteSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var key = _functions.GetPathSliceFromURL(r,0)
	var Response _struct.ResponseData								
	var rep = _sessions.Repository() 
	rep.Delete(key)
	Response.Status = http.StatusOK
	Response.Message = "Deleted " + _sessions.SessionKeyStr
	Response.Data =  key
	_functions.RestponWithJson(response, http.StatusOK, Response)
			
}

func SetSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var dbData _struct.DatabaseAttributes

	dbData.Location = "localhost"	
	dbData.Port = 3050
	dbData.Password = "masterkey"
	dbData.User = "SYSDBA"

	var key = _functions.GetPathSliceFromURL(r,0)
	_functions.GetSessionParamsFromURL(r , &dbData)		
	var Response _struct.ResponseData								
	var rep = _sessions.Repository()
	
	var cmd string = config.MakeConnectionStringFromStruct(dbData)
	rep.Add(key,cmd)
	Response.Status = http.StatusOK
	Response.Message = "Set " + _sessions.SessionKeyStr
	Response.Data =  key
	_functions.RestponWithJson(response, http.StatusOK, Response)
			
}

// GetSQLData returns result rows in json format from an given sql statement.
// The attribute can be given by body or url of statement.
// Following attributes are possible:
//    Location -> database location such as ip, webaddress, default localhost
//    Database -> database file path on database server
//    Port     -> communicating port of database, default 3050
//    User     -> user to logon database, default 'SYSDBA' as it's default user in previous Firebird versions
//    Password -> password to logon database, default 'masterkey' as it's default password in previous Firebird versions
//    Info     -> info string wich can be used to check weather response of REST api is working
//    Cmd      -> SQL command
func GetSQLData(response http.ResponseWriter, r *http.Request) {


	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func GetSQLData")
	
	var key = _functions.GetPathSliceFromURL(r,0)
	var kv  = _functions.KeyValid(response, key)
	if(!kv.Valid) {
		return
	}
	
	var Response _struct.ResponseData
	var entitiesData _struct.SQLAttributes
	
	
	_functions.GetSQLParamsFromURL(r , &entitiesData)				
	_functions.GetParamsFromBODY(r , &entitiesData)	
	_functions.OutParameters(entitiesData) 	
	db, err := config.ConnLocationWithSession(kv)	
	
	if err != nil {
		Response.Status = http.StatusInternalServerError
		Response.Message = entitiesData.Info
		Response.Data = err.Error()
		_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		IsiData, err2 := _models.GetSQLData(entitiesData.Cmd)
		if err2 != nil {
			Response.Status = http.StatusInternalServerError
			Response.Message = entitiesData.Info
			Response.Data = err2.Error()
			_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status = http.StatusOK
			Response.Message = entitiesData.Info
			Response.Data = &IsiData
			_functions.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}



func GetTableData(response http.ResponseWriter, r *http.Request) {


	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func GetTableData")
	
	var key = _functions.GetPathSliceFromURL(r,1)
	var kv  = _functions.KeyValid(response, key)
	if(!kv.Valid) {
		return
	}
	var entitiesData _struct.GetTABLEAttributes
	var table = _functions.GetPathSliceFromURL(r,0)
	entitiesData.Table = table
	var Response _struct.ResponseData

	_functions.GetTableParamsFromURL(r , &entitiesData)		
	_functions.GetTableParamsFromBODY(r , &entitiesData)	
	_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status = http.StatusInternalServerError
		Response.Message = entitiesData.Info
		Response.Data = "No Table given"
		_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data = nil
		_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeSQL(entitiesData)
		
		IsiData, err2 := _models.GetSQLData(cmd)
		if err2 != nil {
			Response.Status = http.StatusInternalServerError
			Response.Message = entitiesData.Info
			Response.Data = err2.Error()
			_functions.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status = http.StatusOK
			Response.Message = entitiesData.Info
			Response.Data = &IsiData
			_functions.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}




