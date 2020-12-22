package _functions

import (
		
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	//"io/ioutil"
	"encoding/json"
	_struct "fbrest/FBxRESTBase/struct"		
	_sessions "fbrest/FBxRESTBase/sessions"
	_httpstuff "fbrest/FBxRESTBase/httpstuff"
	"net/http"
	//"net/url"
	"html/template"	
	"path"
	"fbrest/FBxRESTBase/config"
	
)

func RestponWithText(w http.ResponseWriter, code int) {
	
	profile := _struct.Profile{config.AppName,  config.Version,config.Copyright, "-MNhE7Yf50sz6U9Hgqae", _sessions.MaxDuration}
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
  
	if err := tmpl.Execute(w, profile); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RestponInfoBusyText(w http.ResponseWriter, code int) {
	
	profile := _struct.Profile{config.AppName,  config.Version,config.Copyright, "-MNhE7Yf50sz6U9Hgqae", _sessions.MaxDuration}
	fp := path.Join("templates", "busy.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
  
	if err := tmpl.Execute(w, profile); err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MakeFieldValue(field string, value string) string {
	var filter string
	if(strings.HasPrefix(value,"'")) {
		if(value == "'*'") {
			filter = ""
		} else if (strings.Contains(value,"*")) {
			value = strings.Replace(value,"*","%",-1)		
			filter = field+" like ("+value+")"
		} else {
			filter = field+" = "+value	
		}
	} else {
		if(HasOperator(value)){
			filter = field+" "+value
		} else {
			filter = field+" = "+value	
		}
	}
	return filter
}

func MakeSelectSQL(entitiesData _struct.GetTABLEAttributes) (cmd string) {
	var sql string
	if(len(entitiesData.Fields) < 1) {
		entitiesData.Fields = "*"
	}

	var limitstr string
	if(entitiesData.First > 0) {
		limitstr = " FIRST "+strconv.Itoa(entitiesData.First)	
	}

	if(entitiesData.Skip > 0) {
		limitstr = limitstr + " SKIP "+strconv.Itoa(entitiesData.Skip)	
	}

	sql = "SELECT"+limitstr+" "+entitiesData.Fields+" FROM " + entitiesData.Table
	if(len(entitiesData.Filter) > 0) {
		sql = sql + " WHERE " + entitiesData.Filter
	}
	if(len(entitiesData.GroupBy) > 0) {
		sql = sql + " GROUP BY " + entitiesData.GroupBy
	}
	if(len(entitiesData.OrderBy) > 0) {
		sql = sql + " ORDER BY " + entitiesData.OrderBy
	}
	cmd = sql
	return cmd
}

func MakeUpdateTableSQL(entitiesData _struct.FIELDVALUEAttributes) (cmd string) {
	
	var cmdHead  = "UPDATE " + entitiesData.Table + " SET " 
	for _, fv := range entitiesData.FieldValue {
		if(len(cmd)>0) { 
			cmd = cmd + " , "
		}
		cmd = cmd + fv
	}
	cmd = cmdHead + cmd + " WHERE " + entitiesData.Filter
	return
}

func MakeDeleteTableSQL(entitiesData _struct.FIELDVALUEAttributes) (cmd string) {
	
	cmd = "DROP TABLE " + entitiesData.Table
	return
}

func MakeDeleteTableFieldSQL(entitiesData _struct.FIELDVALUEAttributes) (cmd string) {
	
	cmd = "ALTER TABLE " + entitiesData.Table + " DROP " + entitiesData.FieldValue[0]
	return
}

func OutParameters(entitiesData _struct.SQLAttributes) {
	const funcstr = "func OutParameters"
	log.Debug(funcstr)
	
	log.WithFields(log.Fields{"Given command     ": entitiesData.Cmd,}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill,}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info,}).Info("")
}

func OutTableParameters(entitiesData _struct.GetTABLEAttributes) {
	const funcstr = "func OutTableParameters"
	log.Debug(funcstr)
	
	log.WithFields(log.Fields{"Given fields      ": entitiesData.Fields,}).Info("")
	log.WithFields(log.Fields{"Given order by    ": entitiesData.OrderBy,}).Info("")
	log.WithFields(log.Fields{"Given group by    ": entitiesData.GroupBy,}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill,}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info,}).Info("")	
}

func GetSQLParamsFromBODY(r *http.Request , entitiesData *_struct.SQLAttributes) (bool) {
	const funcstr = "func GetSQLParamsFromBODY"
	log.Debug(funcstr) 
	var xdata _struct.GetUrlSQLAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error(funcstr)
		return false
	}
	entitiesData.Cmd = ReplaceAritmetics(xdata.Cmd)	
	entitiesData.Info = xdata.Info
	return true
}

func GetTableParamsFromBODY(r *http.Request , entitiesData *_struct.GetTABLEAttributes) (ok bool) {
	const funcstr = "func GetTableParamsFromBODY"
	log.Debug(funcstr) 
	var xdata _struct.GetUrlTABLEAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error(funcstr)
		return false
	}
	entitiesData.Fields = strings.Join(xdata.Fields,",")
	entitiesData.Filter = xdata.Filter
	entitiesData.GroupBy = strings.Join(xdata.GroupBy,",")
	entitiesData.OrderBy = strings.Join(xdata.OrderBy,",")
	entitiesData.Skip = xdata.Skip
	entitiesData.First = xdata.First	
	return true
}
func GetFieldParamsFromBODY(r *http.Request , entitiesData *_struct.GetTABLEAttributes) (ok bool) {
	const funcstr = "func GetFieldParamsFromBODY"
	log.Debug(funcstr) 
	var xdata _struct.GetUrlTABLEAttributes
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error(funcstr)
		return false
	}
	entitiesData.Fields = strings.Join(xdata.Fields,",")
	entitiesData.Filter = xdata.Filter
	entitiesData.GroupBy = strings.Join(xdata.GroupBy,",")
	entitiesData.OrderBy = strings.Join(xdata.OrderBy,",")
	entitiesData.Skip = xdata.Skip
	entitiesData.First = xdata.First	
	return true
}
func GetSQLParamsFromURL(r *http.Request , entitiesData *_struct.SQLAttributes) (ok bool) {
	const funcstr = "func GetSQLParamsFromURL"
	var u = r.URL
	
	if(strings.HasPrefix(u.RawQuery,_struct.FormatJson)) {			
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
			
		if(len(par) > 0) {
			xdata := &_struct.GetUrlSQLAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if(err != nil) {
				return false
			}

			log.Debug(xdata)
			entitiesData.Cmd = xdata.Cmd
			entitiesData.Info = xdata.Info			
			return true
		}
	}
	var okret bool
	urlparams, ok := r.URL.Query()["cmd"]
	infoparams, okinfo := r.URL.Query()["info"]
	log.WithFields(log.Fields{"URL params length": len(urlparams),	}).Debug(funcstr)

    if ok && len(urlparams[0]) > 0 {
		
		var cmd string =  string(urlparams[0])
		if(strings.HasPrefix(cmd,"'")||strings.HasPrefix(cmd,"-")) {
				
			entitiesData.Sepfill = cmd[:1]
            cmd = cmd[1:]
			cmd = strings.ReplaceAll(cmd, entitiesData.Sepfill, " ")	
			log.WithFields(log.Fields{"Sepfill": entitiesData.Sepfill,	}).Debug(funcstr+"->set key")				
		}
			    
		log.WithFields(log.Fields{"Cmd": cmd,	}).Debug(funcstr+"->set key")
		entitiesData.Cmd = cmd
		okret = true;
	}
	if okinfo && len(infoparams[0]) > 0 {
				
		var info string =  string(infoparams[0])				
		log.WithFields(log.Fields{"Info": info,	}).Debug(funcstr+"->set key")
		entitiesData.Info = info
		okret = true
	}

	return okret
}

func GetSessionParamsFromBODY(r *http.Request , entitiesData *_struct.DatabaseAttributes) (bool) {
	const funcstr = "func GetSessionParamsFromBODY"
	log.Debug(funcstr) 
	var xdata _struct.DatabaseAttributes
	
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error(funcstr)
		return false
	}
	
	entitiesData.Database = xdata.Database
	entitiesData.Location = xdata.Location
	entitiesData.Password = xdata.Password
	entitiesData.User     = xdata.User
	entitiesData.Port     = xdata.Port
	
	return true
}

func GetSessionParamsFromURL(r *http.Request , entitiesData *_struct.DatabaseAttributes) (bool) {
	
	const funcstr = "func GetSessionParamsFromURL"

	var u = r.URL
	if(strings.HasPrefix(u.RawQuery,_struct.FormatJson)) {			
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)		
	//	 par = "{\"location\":\"localhost\",\"database\":\"D:/Data/Dokuments/DOKUMENTS30.FDB\",\"port\":3050,\"password\":\"su\",\"user\":\"superuser\"}"	
		if(len(par) > 0) {
			xdata := &_struct.DatabaseAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if(err != nil) {
				return false
			}

			log.Info(xdata)
			entitiesData.Database = xdata.Database
			entitiesData.Port = xdata.Port
			entitiesData.Password = xdata.Password
			entitiesData.User = xdata.User		
			entitiesData.Location = xdata.Location	
			return true
		}
	}

	var okret bool
	databaseparams, databaseok := u.Query()["Database"]
	locationparams, locationok := u.Query()["Location"]
	portparams, portok := u.Query()["Port"]
	userparams, userok := u.Query()["User"]
	passwordparams, passwordok := u.Query()["Password"]
	
	if databaseok && len(databaseparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(databaseparams),	}).Debug(funcstr)
		urlparam := databaseparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Database = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Database = urlparam
		}
		okret = true
	}

	if locationok && len(locationparams[0]) > 0 {
		log.WithFields(log.Fields{"URL location length": len(locationparams),	}).Debug(funcstr)
		urlparam := locationparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Location = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Location = urlparam
		}
		okret = true
	}

	if portok && len(portparams[0]) > 0 {
		log.WithFields(log.Fields{"URL port length": len(portparams),	}).Debug(funcstr)
		urlparam := portparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Port,_ = strconv.Atoi(urlparam[1:len(urlparam)-1])
		} else {
			entitiesData.Port,_ = strconv.Atoi(urlparam)
		}
		okret = true
	}

	if userok && len(userparams[0]) > 0 {
		log.WithFields(log.Fields{"URL user length": len(userparams),	}).Debug(funcstr)
		urlparam := userparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.User = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.User = urlparam
		}
		okret = true
	}

	if passwordok && len(passwordparams[0]) > 0 {
		log.WithFields(log.Fields{"URL password length": len(passwordparams),	}).Debug(funcstr)
		urlparam := passwordparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Password = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Password = urlparam
		}
		okret = true
	}

	return okret
}

func GetFIELDPayloadFromString2(params string ,  entitiesData *_struct.FIELDVALUEAttributes) {
	
	// payload=(id:1, username: 'admin', email: 'email@example.org')

	const funcstr = "func GetFIELDPayloadFromString"
	var psplit  = "&"
	var csplit  = "="
	var par = strings.Split(params,psplit)

	log.WithFields(log.Fields{"URL params length": len(par),	}).Debug(funcstr)
	
    if len(par) > 0 {
		for _, pars := range par {
			params = _httpstuff.UnEscape(pars)
			keyval :=  strings.SplitN(params,csplit,2)
						
			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug(funcstr+"->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug(funcstr+"->found value")
			
			if(strings.EqualFold(string(keyval[0]), string("FIELDS"))) {								
				log.WithFields(log.Fields{"Fields": string(keyval[1]),	}).Debug(funcstr+"->set Fields")
				
			}			
		}		
	} 
	return
}

func GetSQLParamsFromString2(params string ,entitiesData *_struct.SQLAttributes) {
	
	const funcstr = "func GetSQLParamsFromString"
	var csplit  = "="
	var par = params

	log.WithFields(log.Fields{"URL params length": len(par),	}).Debug(funcstr)
    params = _httpstuff.UnEscape(par)
	log.WithFields(log.Fields{"SQL": params,	}).Debug(funcstr+"->set key")

	
	keyval :=  strings.SplitN(params,csplit,2)
	
	log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug(funcstr+"->found key")
	log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug(funcstr+"->found value")
	
	if(strings.EqualFold(string(keyval[0]), string("CMD"))) {
		var cmd string =  string(keyval[1])
		if(strings.HasPrefix(cmd,"'")||strings.HasPrefix(cmd,"-")) {
		
			entitiesData.Sepfill = cmd[:1]
			cmd = cmd[1:]
			cmd = strings.ReplaceAll(cmd, entitiesData.Sepfill, " ")	
			log.WithFields(log.Fields{"Sepfill": entitiesData.Sepfill,	}).Debug(funcstr+"->set key")				
		}
		
		log.WithFields(log.Fields{"Cmd": cmd,	}).Debug(funcstr+"->set key")
		entitiesData.Cmd = cmd
	}
	
	if(strings.EqualFold(string(keyval[0]), string("INFO"))) {								
		log.WithFields(log.Fields{"Info": string(keyval[1]),	}).Debug(funcstr+"->set key")
		entitiesData.Info = string(keyval[1])
	}							
	return
}

//Returns the last-nLeft slice from URL
//e.g. when nLeft == 0 returns the last slice
func GetRightPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	
	urlstr := string(r.URL.String())
	var keyval =  strings.SplitN(urlstr,"?",2)
	
	urlstr = keyval[0]
	t2 :=  strings.Split(urlstr,"/")
	key = t2[len(t2)-1-nLeft]
	log.WithFields(log.Fields{_sessions.SessionTokenStr: key,	}).Debug("func GetPathSliceFromURL")	
	return key
}

func GetLeftPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	const funcstr = "func GetLeftPathSliceFromURL"
	urlstr := string(r.URL.String())
	var keyval =  strings.SplitN(urlstr,"?",2)
	
	urlstr = keyval[0]
	t2 :=  strings.Split(urlstr,"/")
	key = t2[nLeft+1]
	log.WithFields(log.Fields{_sessions.SessionTokenStr: key,	}).Debug(funcstr)	
	return key
}

func GetTableParamsFromURL(r *http.Request , entitiesData *_struct.GetTABLEAttributes) (bool) {
	
	//  http://localhost:4488/{{.Key}}/rest/get/TSTANDORT?fjson={"table": "TSTANDORT","fields": ["ID","BEZ","GUELTIG"],"filter":"ID=1 AND BEZ like 'x%'","order": ["BEZ ASC","ID DESC"],"groupby": ["ID","BEZ"],"first": 0}
	
	const funcstr = "func GetTableParamsFromURL"

	var u = r.URL
	if(strings.HasPrefix(u.RawQuery,_struct.FormatJson)) {			
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
			

		if(len(par) > 0) {
			xdata := &_struct.GetUrlTABLEAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if(err != nil) {
				return false
			}

			log.Info(xdata)
			entitiesData.Fields = strings.Join(xdata.Fields,",")
			entitiesData.Filter = xdata.Filter
			entitiesData.GroupBy = strings.Join(xdata.GroupBy,",")
			entitiesData.OrderBy = strings.Join(xdata.OrderBy,",")
			entitiesData.Skip = xdata.Skip
			entitiesData.First = xdata.First			
			return true
		}
	}
	
	var okret bool
	fieldparams, okfields := u.Query()[_struct.Fields]
	orderparams, okorder := u.Query()[_struct.Order]
	filterparams, okfilter := u.Query()[_struct.Filter]
	groupparams, okgroup := u.Query()[_struct.Group]
	infoparams, okinfo := u.Query()[_struct.Info]
	limitparams, oklimit := u.Query()[_struct.Limit]
	
	log.WithFields(log.Fields{"fieldparams": fieldparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"filterparams": filterparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"orderparams": orderparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"groupparams": groupparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"infoparams": infoparams,}).Debug(funcstr)

		
    if okfields && len(fieldparams[0]) > 0 {
		
		urlparam := fieldparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Fields = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Fields = urlparam
		}
		if(len(entitiesData.Fields) < 1) {
			entitiesData.Fields = "*"
		}
		okret = true
	} 

	if okfilter && len(filterparams[0]) > 0 {
		
		urlparam := filterparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Filter = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Filter = urlparam
		}
		okret = true
	} 

	if okorder && len(orderparams[0]) > 0 {
		
		urlparam := orderparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.OrderBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.OrderBy = urlparam
		}
		okret = true
	} 

	if okgroup && len(groupparams[0]) > 0 {
		
		urlparam := groupparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.GroupBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.GroupBy = urlparam
		}
		okret = true
	}

	entitiesData.First = -1
	entitiesData.Skip  = -1
	if oklimit && len(limitparams[0]) > 0 {
		
		urlparam := limitparams[0]
		var limit string
		if(urlparam[:1] == "(") {
			limit = urlparam[1:len(urlparam)-1]
		} else {
			limit = urlparam
		}
		var lm = strings.Split(limit,",")
		if(len(lm) == 1) {
			entitiesData.First,_ = strconv.Atoi(lm[0])
			entitiesData.Skip = 0
		} else if (len(lm) == 2) {
			entitiesData.First,_ = strconv.Atoi(lm[0])
			entitiesData.Skip,_ = strconv.Atoi(lm[1])		
		}
		okret = true
	}

	if okinfo && len(infoparams[0]) > 0 {
		
		urlparam := infoparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Info = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Info = urlparam
		}
		okret = true
	}

	return okret
}

func GetFieldParamsFromURL(r *http.Request , entitiesData *_struct.GetTABLEAttributes) (bool) {
	
	//  http://localhost:4488/{{.Key}}/rest/get/TSTANDORT?fjson={"table": "TSTANDORT","fields": ["ID","BEZ","GUELTIG"],"filter":"ID=1 AND BEZ like 'x%'","order": ["BEZ ASC","ID DESC"],"groupby": ["ID","BEZ"],"first": 0}
	
	const funcstr = "func GetFieldParamsFromURL"

	var u = r.URL
	if(strings.HasPrefix(u.RawQuery,_struct.FormatJson)) {			
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
			

		if(len(par) > 0) {
			xdata := &_struct.GetUrlTABLEAttributes{}

			err := json.Unmarshal([]byte(par), &xdata)

			if(err != nil) {
				return false
			}

			log.Info(xdata)
			entitiesData.Fields = strings.Join(xdata.Fields,",")
			entitiesData.Filter = xdata.Filter
			entitiesData.GroupBy = strings.Join(xdata.GroupBy,",")
			entitiesData.OrderBy = strings.Join(xdata.OrderBy,",")
			entitiesData.Skip = xdata.Skip
			entitiesData.First = xdata.First			
			return true
		}
	}
	
	var okret bool
	
	orderparams, okorder := u.Query()[_struct.Order]
	filterparams, okfilter := u.Query()[_struct.Filter]
	groupparams, okgroup := u.Query()[_struct.Group]
	infoparams, okinfo := u.Query()[_struct.Info]
	limitparams, oklimit := u.Query()[_struct.Limit]
	
	log.WithFields(log.Fields{"filterparams": filterparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"orderparams": orderparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"groupparams": groupparams,}).Debug(funcstr)
	log.WithFields(log.Fields{"infoparams": infoparams,}).Debug(funcstr)

	if okfilter && len(filterparams[0]) > 0 {
		urlparam := filterparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Filter = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Filter = urlparam
		}
		okret = true
	} 

	if okorder && len(orderparams[0]) > 0 {
		urlparam := orderparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.OrderBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.OrderBy = urlparam
		}
		okret = true
	} 

	if okgroup && len(groupparams[0]) > 0 {	
		urlparam := groupparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.GroupBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.GroupBy = urlparam
		}
		okret = true
	}

	entitiesData.First = -1
	entitiesData.Skip  = -1
	if oklimit && len(limitparams[0]) > 0 {
		
		urlparam := limitparams[0]
		var limit string
		if(urlparam[:1] == "(") {
			limit = urlparam[1:len(urlparam)-1]
		} else {
			limit = urlparam
		}
		var lm = strings.Split(limit,",")
		if(len(lm) == 1) {
			entitiesData.First,_ = strconv.Atoi(lm[0])
			entitiesData.Skip = 0
		} else if (len(lm) == 2) {
			entitiesData.First,_ = strconv.Atoi(lm[0])
			entitiesData.Skip,_ = strconv.Atoi(lm[1])		
		}
		okret = true
	}

	if okinfo && len(infoparams[0]) > 0 {
		
		urlparam := infoparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Info = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Info = urlparam
		}
		okret = true
	}

	return okret
}

func GetFIELDPayloadFromBODY(r *http.Request , entitiesData *_struct.FIELDVALUEAttributes) (bool) {
	log.Debug("func GetFIELDPayloadFromBODY") 
	var xdata _struct.GetUrlPayloadAttributes
	//body, err2 := ioutil.ReadAll(r.Body)
    //var str string = string(body)
	//log.Info(str)
	err2 := json.NewDecoder(r.Body).Decode(&xdata)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error("func GetFIELDPayloadFromBODY")
		return false
	}
	
	entitiesData.FieldValue = xdata.Payload
	xdata.Filter = ReplaceAritmetics(xdata.Filter)
	
	entitiesData.Filter = xdata.Filter
	return true
}

func GetFIELDPayloadFromURL(r *http.Request , entitiesData *_struct.FIELDVALUEAttributes) (bool){
	//  http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?payload=(bez='Nürnberg2')&filter=(bez='Nürnberg1')  
	//  http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?ftext="payload=(bez='Nürnberg2')&filter=(bez='Nürnberg1')"
	// http://localhost:4488/{{.Key}}/rest/put/TSTANDORT?fjson={"payload":["ID='123'","BEZ='test'","GUELTIG=1"], "filter": "ID=1 AND BEZ like 'x%'"}

	const funcstr = "func GetFIELDPayloadFromURL"
	var u = r.URL
	
	if(strings.HasPrefix(u.RawQuery,_struct.FormatJson)) {			
		var par = u.RawQuery[len(_struct.FormatJson)+1:]
		par = _httpstuff.UnEscape(par)
		if(len(par) > 0) {
			xdata := &_struct.GetUrlPayloadAttributes{}
			err := json.Unmarshal([]byte(par), &xdata)
			if(err != nil) {
				return false
			}

			log.Info(xdata)
			for _, vals := range xdata.Payload {
				entitiesData.FieldValue = append(entitiesData.FieldValue,vals)
			}
			entitiesData.Filter = xdata.Filter
			return true
		}
	}
	var ok bool
	payloadparams, okpayload := r.URL.Query()[_struct.Payload]
	filterparams, okfilter := r.URL.Query()[_struct.Filter]

	log.WithFields(log.Fields{"payloadparams": payloadparams,}).Debug(funcstr)
			
    if okpayload && len(payloadparams[0]) > 0 {
		var pars = payloadparams[0]
		log.Info(pars)
		
		var st = pars[:1]
		if(st == "(") {
			pars = pars[1:]
		}	

		st = pars[len(pars)-1:]	
		if(st == ")") {
			pars = pars[:len(pars)-1]
		}

		keyval :=  strings.SplitN(pars,",",2)

		for _, pars := range keyval {
			entitiesData.FieldValue = append(entitiesData.FieldValue,pars)	
		}
		ok = true
	} 

	if okfilter && len(filterparams[0]) > 0 {	
		entitiesData.Filter = filterparams[0]
		ok = true
	}
	
	return ok
}