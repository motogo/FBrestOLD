package apis

import (
	"fbrest/FBxRESTBase/config"
	"fbrest/FBxRESTBase/models"	
	"strconv"	
	"strings"
	_struct "fbrest/FBxRESTBase/struct"
	_functions "fbrest/FBxRESTBase/functions"
	_sessions "fbrest/FBxRESTBase/sessions"
	_permissions "fbrest/FBxRESTBase/permissions"
	_httpstuff "fbrest/FBxRESTBase/httpstuff"
	_apperrors "fbrest/FBxRESTBase/apperrors"
	"net/http"
	log "github.com/sirupsen/logrus"	
	bguuid "github.com/kjk/betterguid"
)

func TestDBOpenClose(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Info("func TestDBOpenClose")
	var key = _functions.GetLeftPathSliceFromURL(r,0)	
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var err = config.TestConnLocation(kv.Value)
	var Response _struct.ResponseData
		
	if err != nil {	
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = kv.Token			
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("Database not available")
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		Response.Status  = http.StatusInternalServerError
		Response.Message = "Database open/close successfully"
		Response.Data    = kv.Token
		log.Info("Database opend/closed successfully, ping sent")
			
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	}
}



func TestSQLResponse(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func TestSQLResponse")
	
	var Response _struct.ResponseData
	var entitiesData _struct.SQLAttributes
	
	if(! _functions.GetSQLParamsFromURL(r , &entitiesData)) {
		_functions.GetSQLParamsFromBODY(r , &entitiesData)	
	}
	_functions.OutParameters(entitiesData) 	
	
	Response.Status  = http.StatusOK
	Response.Message = "Test SQL response"
	Response.Data    = entitiesData.Cmd
	_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	
}

func TestTABLEResponse(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func TestSQLResponse")
	
	var Response _struct.ResponseData
	var entitiesData _struct.GetTABLEAttributes
	
	_functions.GetTableParamsFromURL(r , &entitiesData)				
	//_functions.GetParamsFromBODY(r , &entitiesData)	
	//_functions.OutParameters(entitiesData) 	
	
	Response.Status  = http.StatusOK
	Response.Message = "Test table response"
	Response.Data    = entitiesData
	_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
}

func GetSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var dbData _struct.DatabaseAttributes

	dbData.Location = "localhost"	
	dbData.Port = 3050
	dbData.Password = "masterkey"
	dbData.User = "SYSDBA"
	
	var Response _struct.ResponseData
		
	if(!_functions.GetSessionParamsFromURL(r , &dbData)) {
		_functions.GetSessionParamsFromBODY(r , &dbData)
	}

	id := bguuid.New()				
		
	var rep = _sessions.Repository() 
	var perm,errperm = _permissions.GetPermissionFromRepository(dbData.User,dbData.Password)
	if(errperm != nil) {
		Response.Status  = http.StatusNotFound
		Response.Message = errperm.Error()
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		return
	}
	var connstr = string(perm.DBUser+":"+perm.DBPassword+"@"+dbData.Location+":"+strconv.Itoa(dbData.Port)+"/"+dbData.Database)		

	var itm = rep.Add(string(id),perm.Type,connstr)
	Response.Status  = http.StatusOK
	Response.Message = "Created token for permission:"+strconv.Itoa(int(perm.Type)) +", duration:"+ itm.Duration.String()
	Response.Data    =  string(id)
	_httpstuff.RestponWithJson(response, http.StatusOK, Response)
	
}
func GetHelp(response http.ResponseWriter, r *http.Request) {

	_functions.RestponWithText(response, http.StatusOK)
}
func DeleteSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var Response _struct.ResponseData								
	var rep = _sessions.Repository() 
	rep.Delete(key)
	Response.Status  = http.StatusOK
	Response.Message = "Deleted " + _sessions.SessionTokenStr
	Response.Data    =  key
	_httpstuff.RestponWithJson(response, http.StatusOK, Response)			
}

func SetSessionKeyGET(response http.ResponseWriter, r *http.Request) {
	r.Method = "POST"
	SetSessionKey(response, r)
}

func SetSessionKey(response http.ResponseWriter, r *http.Request) {
	
	var dbData _struct.DatabaseAttributes

	dbData.Location = "localhost"	
	dbData.Port = 3050
	dbData.Password = "masterkey"
	dbData.User     = "guest"

	var key = _functions.GetLeftPathSliceFromURL(r,0)
	_functions.GetSessionParamsFromURL(r , &dbData)		
	var Response _struct.ResponseData								
	var rep = _sessions.Repository()
	
	var cmd string = config.MakeConnectionStringFromStruct(dbData)
	var perm,errperm = _permissions.GetPermissionFromRepository(dbData.User,dbData.Password )
	if(errperm != nil) {
		Response.Status  = http.StatusOK
		Response.Message = errperm.Error()
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		return
	}
	rep.Add(key,perm.Type,cmd)
	
	Response.Status  = http.StatusOK
	Response.Message = "Set " + _sessions.SessionTokenStr
	Response.Data    =  key
	_httpstuff.RestponWithJson(response, http.StatusOK, Response)
			
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
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	
	var Response _struct.ResponseData
	var entitiesData _struct.SQLAttributes
		
	if(! _functions.GetSQLParamsFromURL(r , &entitiesData)) {
		_functions.GetSQLParamsFromBODY(r , &entitiesData)	
	}
	_functions.OutParameters(entitiesData) 	
	db, err := config.ConnLocationWithSession(kv)	
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = entitiesData.Info
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		IsiData, err2 := _models.GetSQLData(entitiesData.Cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = entitiesData.Info
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = "Got datas"
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}

// localhost:4488/token/rest/get/TSTANDORT?fields=(id,bez)&order=(bez asc, id desc)&filter=(BEZ like 'T%')&group=(id,bez)</li>
// localhost:4488/token/rest/get/TSTANDORT?ftext="fields=(id,bez)&order=(bez asc, id desc)&filter=(BEZ like 'T%')&group=(id,bez)"</li>
// localhost:4488/token/rest/get/TSTANDORT?fjson="{"fields":["id","bez"]},{"order":["bez asc","id desc"]},{"filter":"BEZ like 'T%'"},{"group":["id","bez"]}"</li>
func GetTableData(response http.ResponseWriter, r *http.Request) {


	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func GetTableData")
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var Response _struct.ResponseData
	if(kv.Permission < _permissions.Read) {
		Response.Status = http.StatusNotFound
		Response.Message = _apperrors.ErrNoPermission.Error() + " (read)"
		Response.Data = key
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}
	
	var entitiesData _struct.GetTABLEAttributes
	var table = _functions.GetRightPathSliceFromURL(r,0)
	entitiesData.Table = table
	

	if(! _functions.GetTableParamsFromURL(r , &entitiesData)		) {
		_functions.GetTableParamsFromBODY(r , &entitiesData)	
	}
	_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = "No Table given"
		Response.Data    = entitiesData.Info
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = entitiesData.Info
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeSelectSQL(entitiesData)
		
		IsiData, err2 := _models.GetSQLData(cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = entitiesData.Info
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = "Got datas"
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}

func GetTableFieldData(response http.ResponseWriter, r *http.Request) {


	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func GetTableData")
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var Response _struct.ResponseData
	if(kv.Permission < _permissions.Read) {
		Response.Status = http.StatusNotFound
		Response.Message = _apperrors.ErrNoPermission.Error() + " (read)"
		Response.Data = key
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}
	
	var entitiesData _struct.GetTABLEAttributes
	var fieldvalue = _functions.GetRightPathSliceFromURL(r,0)
	var field = _functions.GetRightPathSliceFromURL(r,1)
	var table = _functions.GetRightPathSliceFromURL(r,2)
	fieldvalue = _httpstuff.UnEscape(fieldvalue)
	
	entitiesData.Table = table
	entitiesData.Fields = field	
	entitiesData.Filter = _functions.MakeFieldValue(field, fieldvalue)

	if(! _functions.GetFieldParamsFromURL(r , &entitiesData)) {
		_functions.GetFieldParamsFromBODY(r , &entitiesData)	
	}
	_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = "No Table given"
		Response.Data    = entitiesData.Info
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = entitiesData.Info
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeSelectSQL(entitiesData)
		
		IsiData, err2 := _models.GetSQLData(cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = entitiesData.Info
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = "Got datas"
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}

func IsBrowser(ua string) (bool) {
	if(strings.Contains(ua, "curl") || strings.Contains(ua, "Postman")) {
		return false
	}
	return true
}

func UpdateTableDataGET(response http.ResponseWriter, r *http.Request) {
	var ua = r.Header.Get("User-Agent")
	if(!IsBrowser(ua)) {
		var Response _struct.ResponseData	
		Response.Status  = http.StatusNotFound
		Response.Message = "GET not allowed, use PUT instead"
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}
    log.Info(ua)
	r.Method = "PUT"
	UpdateTableData(response,r)
}

func UpdateTableData(response http.ResponseWriter, r *http.Request) {

	// http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?payload=(username='admin', email='email@example.org')&filter=(id=1)

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func GetTableData")
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var Response _struct.ResponseData
	if(kv.Permission < _permissions.Read) {
		Response.Status  = http.StatusNotFound
		Response.Message = _apperrors.ErrNoPermission.Error() + " (read)"
		Response.Data    = "Token:"+key
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	var entitiesData _struct.FIELDVALUEAttributes
	var table = _functions.GetRightPathSliceFromURL(r,0)
	entitiesData.Table = table
	

	if(! _functions.GetFIELDPayloadFromURL(r , &entitiesData)) {
	    _functions.GetFIELDPayloadFromBODY(r , &entitiesData)	
	}
	//_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = "Update failure"
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeUpdateTableSQL(entitiesData)
		
		IsiData, err2 := _models.GetSQLData(cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = cmd
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = "Got datas"
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}


func DeleteTableDataGET(response http.ResponseWriter, r *http.Request) {
	var ua = r.Header.Get("User-Agent")
	if(!IsBrowser(ua)) {
		var Response _struct.ResponseData	
		Response.Status  = http.StatusNotFound
		Response.Message = "GET not allowed, use POST instead"
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}
	r.Method = "POST"
	DeleteTableData(response, r)
}
func DeleteTableData(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func DeleteTableData")
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var Response _struct.ResponseData
	if(kv.Permission < _permissions.Read) {
		Response.Status  = http.StatusNotFound
		Response.Message = _apperrors.ErrNoPermission.Error() + " (read)"
		Response.Data    = key
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	var entitiesData _struct.FIELDVALUEAttributes
	var table = _functions.GetRightPathSliceFromURL(r,0)
	
	entitiesData.Table = table
	
	//_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = _apperrors.NoTableGiven
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeDeleteTableSQL(entitiesData)
		
		IsiData, err2 := _models.RunSQLData(cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = cmd
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = _apperrors.GotDatas
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}

func DeleteTableFieldDataGET(response http.ResponseWriter, r *http.Request) {
	r.Method = "POST"
	DeleteTableFieldData(response, r)
}

func DeleteTableFieldData(response http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.URL,	}).Debug("func DeleteTableFieldData")
	
	var key = _functions.GetLeftPathSliceFromURL(r,0)
	var kv  = _sessions.TokenValid(response, key)
	if(!kv.Valid) {
		return
	}
	var Response _struct.ResponseData
	if(kv.Permission < _permissions.Read) {
		Response.Status  = http.StatusNotFound
		Response.Message = _apperrors.ErrNoPermission.Error() + " (read)"
		Response.Data    = key
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	var entitiesData _struct.FIELDVALUEAttributes
	var table = _functions.GetRightPathSliceFromURL(r,1)
	var field = _functions.GetRightPathSliceFromURL(r,0)
	entitiesData.Table = table	
	entitiesData.FieldValue = append(entitiesData.FieldValue,field);
	
	//_functions.OutTableParameters(entitiesData) 	
	
	if(len(entitiesData.Table) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = _apperrors.NoTableGiven
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	if(len(entitiesData.FieldValue) < 1) {
		Response.Status  = http.StatusInternalServerError
		Response.Message = _apperrors.NoFieldsGiven
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		return
	}

	db, err := config.ConnLocationWithSession(kv)
	
	if err != nil {
		Response.Status  = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data    = _apperrors.DataNil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		var cmd string = _functions.MakeDeleteTableFieldSQL(entitiesData)
		
		IsiData, err2 := _models.RunSQLData(cmd)
		if err2 != nil {
			Response.Status  = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data    = _apperrors.DataNil
			_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		} else {
			Response.Status  = http.StatusOK
			Response.Message = _apperrors.GotDatas
			Response.Data    = &IsiData
			_httpstuff.RestponWithJson(response, http.StatusOK, Response)
		}
	}
}




