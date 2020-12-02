package _functions

import (
		
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"encoding/json"
	_struct "fbrest/Base/struct"		
	_sessions "fbrest/Base/sessions"
	"net/http"
	"html/template"	
	"path"
	"fbrest/Base/config"
	"time"
)


func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}



func RestponWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type Profile struct {
	Appname    string
	Version string
	Copyright string
	Duration time.Duration
  }

func RestponWithText(w http.ResponseWriter, code int) {
	

	profile := Profile{config.AppName,  config.Version,config.Copyright, _sessions.MaxDuration}
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


func MakeSQL(entitiesData _struct.GetTABLEAttributes) (cmd string) {
	var sql string
	if(len(entitiesData.Fields) < 1) {
		entitiesData.Fields = "*"
	}
	sql = "SELECT "+entitiesData.Fields+" FROM " + entitiesData.Table
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

func OutParameters(entitiesData _struct.SQLAttributes) {
	log.Debug("func OutParameters")
	
	log.WithFields(log.Fields{"Given command     ": entitiesData.Cmd,}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill,}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info,}).Info("")
}

func OutTableParameters(entitiesData _struct.GetTABLEAttributes) {
	log.Debug("func OutTableParameters")
	
	log.WithFields(log.Fields{"Given fields      ": entitiesData.Fields,}).Info("")
	log.WithFields(log.Fields{"Given order by    ": entitiesData.OrderBy,}).Info("")
	log.WithFields(log.Fields{"Given group by    ": entitiesData.GroupBy,}).Info("")
	log.WithFields(log.Fields{"Given sepfill     ": entitiesData.Sepfill,}).Info("")
	log.WithFields(log.Fields{"Given info        ": entitiesData.Info,}).Info("")	
}

func GetParamsFromBODY(r *http.Request , entitiesData *_struct.SQLAttributes) {
	
	log.Debug("func GetParamsFromBODY") 
	err2 := json.NewDecoder(r.Body).Decode(&entitiesData)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error("func GetParamsFromBODY")
	}
	
}

func GetTableParamsFromBODY(r *http.Request , entitiesData *_struct.GetTABLEAttributes) {
	log.Debug("func GetTableParamsFromBODY") 
	err2 := json.NewDecoder(r.Body).Decode(&entitiesData)
	if err2 != nil {			
		log.WithFields(log.Fields{"Decode params to JSON": err2.Error(),	}).Error("func GetTableParamsFromBODY")
	}
}

func GetParamsFromURL(r *http.Request , entitiesData *_struct.SQLAttributes) {
	
	urlparams, ok := r.URL.Query()["q"]
	log.WithFields(log.Fields{"URL params length": len(urlparams),	}).Debug("func GetParamsFromURL")



    if ok && len(urlparams[0]) > 0 {
		
		urlparam := urlparams[0]
		s :=  strings.Split(urlparam,",")
		for _, pars := range s {
			
			keyval :=  strings.SplitN(pars,":",2)
			
			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug("func GetParamsFromURL->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug("func GetParamsFromURL->found value")
			
			if(strings.EqualFold(string(keyval[0]), string("CMD"))) {
				var cmd string =  string(keyval[1])
				if(strings.HasPrefix(cmd,"'")||strings.HasPrefix(cmd,"-")) {
				
					entitiesData.Sepfill = cmd[:1]
                    cmd = cmd[1:]
					cmd = strings.ReplaceAll(cmd, entitiesData.Sepfill, " ")	
					log.WithFields(log.Fields{"Sepfill": entitiesData.Sepfill,	}).Debug("func GetParamsFromURL->set key")				
			    }
			    
				log.WithFields(log.Fields{"Cmd": cmd,	}).Debug("func GetParamsFromURL->set key")
				entitiesData.Cmd = cmd
			}
			
			if(strings.EqualFold(string(keyval[0]), string("INFO"))) {								
				log.WithFields(log.Fields{"Info": string(keyval[1]),	}).Debug("func GetParamsFromURL->set key")
				entitiesData.Info = string(keyval[1])
			}
		}		
	} 
	return
}

func GetSessionParamsFromURL(r *http.Request , entitiesData *_struct.DatabaseAttributes) {
	
	urlparams, ok := r.URL.Query()["h"]
	if ok && len(urlparams[0]) > 0 {
		GetSessionParamsFromURL2(r , entitiesData)
		return
	}

	databaseparams, databaseok := r.URL.Query()["Database"]
	locationparams, locationok := r.URL.Query()["Location"]
	portparams, portok := r.URL.Query()["Port"]
	userparams, userok := r.URL.Query()["User"]
	passwordparams, passwordok := r.URL.Query()["Password"]
	
	log.WithFields(log.Fields{"URL params length": len(urlparams),	}).Debug("func GetSessionParamsFromURL")
	
	if databaseok && len(databaseparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(databaseparams),	}).Debug("func GetSessionParamsFromURL")
		urlparam := databaseparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Database = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Database = urlparam
		}
	}

	if locationok && len(locationparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(locationparams),	}).Debug("func GetSessionParamsFromURL")
		urlparam := locationparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Location = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Location = urlparam
		}
	}

	if portok && len(portparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(portparams),	}).Debug("func GetSessionParamsFromURL")
		urlparam := portparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Port,_ = strconv.Atoi(urlparam[1:len(urlparam)-1])
		} else {
			entitiesData.Port,_ = strconv.Atoi(urlparam)
		}
	}

	if userok && len(userparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(userparams),	}).Debug("func GetSessionParamsFromURL")
		urlparam := userparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.User = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.User = urlparam
		}
	}

	if passwordok && len(passwordparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(passwordparams),	}).Debug("func GetSessionParamsFromURL")
		urlparam := passwordparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Password = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Password = urlparam
		}
	}
    
	return
}

func GetSessionParamsFromURL2(r *http.Request , entitiesData *_struct.DatabaseAttributes) {
	
	urlparams, ok := r.URL.Query()["h"]
	
	
	
	log.WithFields(log.Fields{"URL params length": len(urlparams),	}).Debug("func GetSessionParamsFromURL")
	

    if ok && len(urlparams[0]) > 0 {
		
		urlparam := urlparams[0]
		s :=  strings.Split(urlparam,",")
		for _, pars := range s {
			
			keyval :=  strings.SplitN(pars,":",2)
			
			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug("func GetSessionParamsFromURL->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug("func GetSessionParamsFromURL->found value")
			
			if(strings.EqualFold(string(keyval[0]), string("LOCATION"))) {								
				log.WithFields(log.Fields{"Location": string(keyval[1]),	}).Debug("func GetSessionParamsFromURL->set key")
				entitiesData.Location = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("DATABASE"))) {								
				log.WithFields(log.Fields{"Database": string(keyval[1]),	}).Debug("func GetSessionParamsFromURL->set key")
				entitiesData.Database = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("PORT"))) {								
				log.WithFields(log.Fields{"Port": string(keyval[1]),	}).Debug("func GetSessionParamsFromURL->set key")
				entitiesData.Port,_ = strconv.Atoi(string(keyval[1]))
			}
			if(strings.EqualFold(string(keyval[0]), string("USER"))) {								
				log.WithFields(log.Fields{"User": string(keyval[1]),	}).Debug("func GetSessionParamsFromURL->set key")
				entitiesData.User = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("PASSWORD"))) {								
				log.WithFields(log.Fields{"Password": "****",	}).Debug("func GetSessionParamsFromURL->set key")
				entitiesData.Password = string(keyval[1])
			}
		}		
	} 
	return
}

func KeyValid(response http.ResponseWriter, key string) (kv _sessions.Items) {
	var Response _struct.ResponseData
	kv  = GetSessionKeyFromRepository(key)	
	if(len(kv.Value) < 1) {
		Response.Status = http.StatusForbidden
		Response.Message = "No valid database connection found by "+_sessions.SessionKeyStr+" "+kv.Key
		Response.Data = nil
		RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false
		return kv
	}
	var duration time.Duration = kv.Duration
	var end = time.Now();
	difference := end.Sub(kv.Start)
	if(difference > duration) {
		//Zeit für Key abgelaufen
		Response.Status = http.StatusForbidden
		Response.Message = _sessions.SessionKeyStr+" "+kv.Key + " has expired after "+ strconv.Itoa(_sessions.MaxDuration) +" seconds"
		Response.Data = nil
		var rep = _sessions.Repository() 
	    rep.Delete(key)
		RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false
		return kv
	}
	log.WithFields(log.Fields{"Remaining "+_sessions.SessionKeyStr+" duration (s)": duration,	}).Debug("func KeyValid")
	kv.Valid = true
	return kv
}

//Returns the last-nLeft slice from URL
//e.g. when nLeft == 0 returns the last slice
func GetPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	
	urlstr := string(r.URL.String())
	var keyval =  strings.SplitN(urlstr,"?",2)
	
	urlstr = keyval[0]
	t2 :=  strings.Split(urlstr,"/")
	key = t2[len(t2)-1-nLeft]
	log.WithFields(log.Fields{_sessions.SessionKeyStr: key,	}).Debug("func GetSessionKeyFromURL")	
	return key
}

func GetSessionKeyFromRepository(key string) (kval _sessions.Items) {
	
	log.WithFields(log.Fields{_sessions.SessionKeyStr: key,	}).Debug("func GetSessionKeyFromRepository")	
	var rep = _sessions.Repository() 
	var result,_ = rep.Get(key)	 
	result.Key = key
	return result 
}

func GetTableParamsFromURL(r *http.Request , entitiesData *_struct.GetTABLEAttributes) {
	//	SELECT * FROM employees where last_name LIKE 'G%' ORDER BY emp_no;		  
	//  http://localhost/api/v2/_table/employees?q&filter=(last_name like 'G%')&order=emp_no)
	//
	//	SELECT id,bez FROM employees where last_name LIKE 'G%' ORDER BY emp_no;		  
	//  http://localhost/api/v2/_table/employees?q&fields=(id,bez)&filter=(last_name%20like%20G%25)&order=emp_no	
	
	urlhintparams, okhint := r.URL.Query()["h"]
	//urlparams, okquery := r.URL.Query()["q"]
	fieldparams, okfields := r.URL.Query()["fields"]
	orderparams, okorder := r.URL.Query()["order"]
	filterparams, okfilter := r.URL.Query()["filter"]
	groupparams, okgroup := r.URL.Query()["group"]

	log.WithFields(log.Fields{"urlhint": urlhintparams,	}).Info("func GetTableData")
	//log.WithFields(log.Fields{"urlquery": urlparams,	}).Info("func GetTableData")
	log.WithFields(log.Fields{"fieldparams": fieldparams,	}).Info("func GetTableData")
	log.WithFields(log.Fields{"filterparams": filterparams,	}).Info("func GetTableData")
	log.WithFields(log.Fields{"orderparams": orderparams,	}).Info("func GetTableData")
	log.WithFields(log.Fields{"groupparams": groupparams,	}).Info("func GetTableData")



	//log.WithFields(log.Fields{"URL params length": len(urlparams),	}).Debug("func GetTableDataParamsFromURL")

	if okhint && len(urlhintparams[0]) > 0 {
		urlparam :=urlhintparams[0]
		s :=  strings.Split(urlparam,",")
	
		for _, pars := range s {
			
			keyval :=  strings.SplitN(pars,":",2)

			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug("func GetTableParamsFromURL->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug("func GetTableParamsFromURL->found value")
			if(strings.EqualFold(string(keyval[0]), string("MIN"))) {				
				log.WithFields(log.Fields{"Function "+keyval[0]: "Params "+string(keyval[1]),	}).Debug("func GetTableDataParamsFromURL->set key")
				entitiesData.Function = string(keyval[0])
				entitiesData.FunctionParams = string(keyval[1])
			}
		}
	}
	
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
	} 

	if okfilter && len(filterparams[0]) > 0 {
		
		urlparam := filterparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Filter = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Filter = urlparam
		}
	} 

	if okorder && len(orderparams[0]) > 0 {
		
		urlparam := orderparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.OrderBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.OrderBy = urlparam
		}
	} 

	if okgroup && len(groupparams[0]) > 0 {
		
		urlparam := groupparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.GroupBy = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.GroupBy = urlparam
		}
	}

	return
}