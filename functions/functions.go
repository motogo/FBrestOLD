package _functions

import (
		
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"encoding/json"
	_struct "fbrest/Base/struct"		
	_sessions "fbrest/Base/sessions"
	//_httpstuff "fbrest/Base/httpstuff"
	"net/http"
	"net/url"
	"html/template"	
	"path"
	"fbrest/Base/config"
	"time"
)


type Profile struct {
	Appname    string
	Version string
	Copyright string
	Key string
	Duration time.Duration
  }

func RestponWithText(w http.ResponseWriter, code int) {
	

	profile := Profile{config.AppName,  config.Version,config.Copyright, "-MNhE7Yf50sz6U9Hgqae", _sessions.MaxDuration}
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
	

	profile := Profile{config.AppName,  config.Version,config.Copyright, "-MNhE7Yf50sz6U9Hgqae", _sessions.MaxDuration}
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

func MakeSQL(entitiesData _struct.GetTABLEAttributes) (cmd string) {
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

func GetSQLParamsFromURL(r *http.Request , entitiesData *_struct.SQLAttributes) {
	const funcstr = "func GetSQLParamsFromURL"

	curlparamst, okt := r.URL.Query()["ftext"]
	if okt && len(curlparamst[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Text
		var par = strings.SplitN(r.RequestURI,"?ftext=",2)
		if(len(par) > 0) {
			GetSQLParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}

	curlparamsj, okj := r.URL.Query()["fjson"]
	if okj && len(curlparamsj[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Json
		var par = strings.SplitN(r.RequestURI,"?fjson=",2)
		if(len(par) > 0) {
			GetSQLParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}

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
	}
	if okinfo && len(infoparams[0]) > 0 {
				
		var info string =  string(infoparams[0])				
		log.WithFields(log.Fields{"Info": info,	}).Debug(funcstr+"->set key")
		entitiesData.Info = info
	}

	return
}

func GetSessionParamsFromURL(r *http.Request , entitiesData *_struct.DatabaseAttributes) {
	
	const funcstr = "func GetSessionParamsFromURL"
	curlparamst, okt := r.URL.Query()["ftext"]
	if okt && len(curlparamst[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Text
		var par = strings.SplitN(r.RequestURI,"?ftext=",2)
		if(len(par) > 0) {
			GetSessionParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}

	curlparamsj, okj := r.URL.Query()["fjson"]
	if okj && len(curlparamsj[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Json
		var par = strings.SplitN(r.RequestURI,"?fjson=",2)
		if(len(par) > 0) {
			GetSessionParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}
	 
	databaseparams, databaseok := r.URL.Query()["Database"]
	locationparams, locationok := r.URL.Query()["Location"]
	portparams, portok := r.URL.Query()["Port"]
	userparams, userok := r.URL.Query()["User"]
	passwordparams, passwordok := r.URL.Query()["Password"]
	
	
	
	if databaseok && len(databaseparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(databaseparams),	}).Debug(funcstr)
		urlparam := databaseparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Database = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Database = urlparam
		}
	}

	if locationok && len(locationparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(locationparams),	}).Debug(funcstr)
		urlparam := locationparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Location = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Location = urlparam
		}
	}

	if portok && len(portparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(portparams),	}).Debug(funcstr)
		urlparam := portparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Port,_ = strconv.Atoi(urlparam[1:len(urlparam)-1])
		} else {
			entitiesData.Port,_ = strconv.Atoi(urlparam)
		}
	}

	if userok && len(userparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(userparams),	}).Debug(funcstr)
		urlparam := userparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.User = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.User = urlparam
		}
	}

	if passwordok && len(passwordparams[0]) > 0 {
		log.WithFields(log.Fields{"URL params length": len(passwordparams),	}).Debug(funcstr)
		urlparam := passwordparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Password = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Password = urlparam
		}
	}

	return
}

func GetSessionParamsFromString(params string , paramtype _struct.ParamFormatType, entitiesData *_struct.DatabaseAttributes) {
	
	
	var psplit  = "&"
	var csplit  = "="
	if(paramtype == _struct.Text) {
		psplit = "&"
		csplit = "="
	} else if(paramtype == _struct.Json) {
		psplit = ","
		csplit = ":"
	}

	var par = strings.Split(params,psplit)

	log.WithFields(log.Fields{"URL params length": len(par),	}).Debug("func GetSessionParamsFromString")
	
    if len(par) > 0 {
		for _, pars := range par {

			keyval :=  strings.SplitN(pars,csplit,2)
			if(paramtype == _struct.Json) {
				var st = keyval[0][:4]
				if(st == "{%22") {
					keyval[0] = keyval[0][4:]
				}	

				st = keyval[0][len(keyval[0])-3:]	
				if(st == "%22") {
					keyval[0] = keyval[0][:len(keyval[0])-3]
				}

				st = keyval[1][:3]
				if(st == "%22") {
					keyval[1] = keyval[1][3:]
				}	

				st = keyval[1][len(keyval[1])-4:]	
				if(st == "%22}") {
					keyval[1] = keyval[1][:len(keyval[1])-4]
				}
			}
			// "x=y&k=n"
			if(paramtype == _struct.Text) {
				var st = keyval[0][:3]
				if(st == "%22") {
					keyval[0] = keyval[0][3:]
				}	

				st = keyval[1][len(keyval[1])-3:]	
				if(st == "%22") {
					keyval[1] = keyval[1][:len(keyval[1])-3]
				}
			}
			
			
			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug("func GetSessionParamsFromString->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug("func GetSessionParamsFromString->found value")
			
			if(strings.EqualFold(string(keyval[0]), string("LOCATION"))) {								
				log.WithFields(log.Fields{"Location": string(keyval[1]),	}).Debug("func GetSessionParamsFromString->set key")
				entitiesData.Location = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("DATABASE"))) {								
				log.WithFields(log.Fields{"Database": string(keyval[1]),	}).Debug("func GetSessionParamsFromString->set key")
				entitiesData.Database = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("PORT"))) {								
				log.WithFields(log.Fields{"Port": string(keyval[1]),	}).Debug("func GetSessionParamsFromString->set key")
				entitiesData.Port,_ = strconv.Atoi(string(keyval[1]))
			}
			if(strings.EqualFold(string(keyval[0]), string("USER"))) {								
				log.WithFields(log.Fields{"User": string(keyval[1]),	}).Debug("func GetSessionParamsFromString->set key")
				entitiesData.User = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("PASSWORD"))) {								
				log.WithFields(log.Fields{"Password": "****",	}).Debug("func GetSessionParamsFromString->set key")
				entitiesData.Password = string(keyval[1])
			}
		}		
	} 
	return
}

func GetTABLEParamsFromString(params string , paramtype _struct.ParamFormatType, entitiesData *_struct.GetTABLEAttributes) {
	
	const funcstr = "func GetTABLEParamsFromString"
	var psplit  = "&"
	var csplit  = "="
	if(paramtype == _struct.Text) {
		psplit = "&"
		csplit = "="
	} else if(paramtype == _struct.Json) {
		psplit = ","
		csplit = ":"
	}

	var par = strings.Split(params,psplit)

	log.WithFields(log.Fields{"URL params length": len(par),	}).Debug(funcstr)
	
    if len(par) > 0 {
		for _, pars := range par {

			keyval :=  strings.SplitN(pars,csplit,2)
			if(paramtype == _struct.Json) {
				var st = keyval[0][:4]
				if(st == "{%22") {
					keyval[0] = keyval[0][4:]
				}	

				st = keyval[0][len(keyval[0])-3:]	
				if(st == "%22") {
					keyval[0] = keyval[0][:len(keyval[0])-3]
				}

				st = keyval[1][:3]
				if(st == "%22") {
					keyval[1] = keyval[1][3:]
				}	

				st = keyval[1][len(keyval[1])-4:]	
				if(st == "%22}") {
					keyval[1] = keyval[1][:len(keyval[1])-4]
				}
			}
			// "x=y&k=n"
			if(paramtype == _struct.Text) {
				var st = keyval[0][:3]
				if(st == "%22") {
					keyval[0] = keyval[0][3:]
				}	

				st = keyval[1][len(keyval[1])-3:]	
				if(st == "%22") {
					keyval[1] = keyval[1][:len(keyval[1])-3]
				}
			}
			
			
			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug(funcstr+"->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug(funcstr+"->found value")
			
			if(strings.EqualFold(string(keyval[0]), string("FIELDS"))) {								
				log.WithFields(log.Fields{"Fields": string(keyval[1]),	}).Debug(funcstr+"->set Fields")
				entitiesData.Fields = string(keyval[1])
				if(len(entitiesData.Fields) < 1) {
					entitiesData.Fields = "*"
				}
			}
			if(strings.EqualFold(string(keyval[0]), string("ORDER"))) {								
				log.WithFields(log.Fields{"Order by": string(keyval[1]),	}).Debug(funcstr+"->set Order by")
				entitiesData.OrderBy = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("GROUP"))) {								
				log.WithFields(log.Fields{"Group by": string(keyval[1]),	}).Debug(funcstr+"->set Group by")
				entitiesData.GroupBy = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("FILTER"))) {								
				log.WithFields(log.Fields{"Filter": string(keyval[1]),	}).Debug(funcstr+"->set Filter")
				entitiesData.Filter = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("INFO"))) {								
				log.WithFields(log.Fields{"Info:": string(keyval[1]),	}).Debug(funcstr+"->set Info")
				entitiesData.Info = string(keyval[1])
			}
			if(strings.EqualFold(string(keyval[0]), string("LIMIT"))) {								
				log.WithFields(log.Fields{"Limit:": string(keyval[1]),	}).Debug(funcstr+"->set Limit")
				var lm = strings.Split(string(keyval[1]),",")
				if(len(lm) == 1) {
					entitiesData.First,_ = strconv.Atoi(lm[0])
					entitiesData.Skip = 0
				} else if (len(lm) == 2) {
					entitiesData.First,_ = strconv.Atoi(lm[0])
					entitiesData.Skip,_ = strconv.Atoi(lm[1])		
				}
			}
		}		
	} 
	return
}

func GetSQLParamsFromString(params string , paramtype _struct.ParamFormatType, entitiesData *_struct.SQLAttributes) {
	
	const funcstr = "func GetSQLParamsFromString"

	var csplit  = "="
	if(paramtype == _struct.Text) {
		
		csplit = "="
	} else if(paramtype == _struct.Json) {
		
		csplit = ":"
	}

	var par = params

	log.WithFields(log.Fields{"URL params length": len(par),	}).Debug(funcstr)
	
	if(paramtype == _struct.Json) {
		var st = params[:4]
		if(st == "{%22") {
			params = params[4:]
		}	

		st = params[len(params)-4:]	
		if(st == "%22}") {
			params = params[:len(params)-4]
		}
	}
	// "x=y&k=n"
	if(paramtype == _struct.Text) {
		var st = params[:3]
		if(st == "%22") {
			params = params[3:]
		}	

		st = params[len(params)-3:]	
		if(st == "%22") {
			params = params[:len(params)-3]
		}
	}
																				
	log.WithFields(log.Fields{"SQL cmd": params,	}).Debug(funcstr+"->set key")

	urlparam,_ := url.QueryUnescape(params)

	keyval :=  strings.SplitN(urlparam,csplit,2)
	
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
func GetPathSliceFromURL(r *http.Request, nLeft int) (key string) {
	
	urlstr := string(r.URL.String())
	var keyval =  strings.SplitN(urlstr,"?",2)
	
	urlstr = keyval[0]
	t2 :=  strings.Split(urlstr,"/")
	key = t2[len(t2)-1-nLeft]
	log.WithFields(log.Fields{_sessions.SessionTokenStr: key,	}).Debug("func GetPathSliceFromURL")	
	return key
}

func GetTableParamsFromURL(r *http.Request , entitiesData *_struct.GetTABLEAttributes) {
	//	SELECT * FROM employees where last_name LIKE 'G%' ORDER BY emp_no;		  
	//  http://localhost/api/v2/_table/employees?q&filter=(last_name like 'G%')&order=emp_no)
	//
	//	SELECT id,bez FROM employees where last_name LIKE 'G%' ORDER BY emp_no;		  
	//  http://localhost/api/v2/_table/employees?q&fields=(id,bez)&filter=(last_name%20like%20G%25)&order=emp_no
	
	const funcstr = "func GetTableParamsFromURL"
	curlparamst, okt := r.URL.Query()["ftext"]
	if okt && len(curlparamst[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Text
		var par = strings.SplitN(r.RequestURI,"?ftext=",2)
		if(len(par) > 0) {
			GetTABLEParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}

	curlparamsj, okj := r.URL.Query()["fjson"]
	if okj && len(curlparamsj[0]) > 0 {
		var paramtype _struct.ParamFormatType = _struct.Json
		var par = strings.SplitN(r.RequestURI,"?fjson=",2)
		if(len(par) > 0) {
			GetTABLEParamsFromString(par[1] , paramtype, entitiesData)
			return
		}
	}

	urlhintparams, okhint := r.URL.Query()["h"]	
	fieldparams, okfields := r.URL.Query()["fields"]
	orderparams, okorder := r.URL.Query()["order"]
	filterparams, okfilter := r.URL.Query()["filter"]
	groupparams, okgroup := r.URL.Query()["group"]
	infoparams, okinfo := r.URL.Query()["info"]
	limitparams, oklimit := r.URL.Query()["limit"]

	log.WithFields(log.Fields{"urlhint": urlhintparams,	}).Debug(funcstr)	
	log.WithFields(log.Fields{"fieldparams": fieldparams,	}).Debug(funcstr)
	log.WithFields(log.Fields{"filterparams": filterparams,	}).Debug(funcstr)
	log.WithFields(log.Fields{"orderparams": orderparams,	}).Debug(funcstr)
	log.WithFields(log.Fields{"groupparams": groupparams,	}).Debug(funcstr)
	log.WithFields(log.Fields{"infoparams": infoparams,	}).Debug(funcstr)


	if okhint && len(urlhintparams[0]) > 0 {
		urlparam :=urlhintparams[0]
		s :=  strings.Split(urlparam,",")
	
		for _, pars := range s {
			
			keyval :=  strings.SplitN(pars,":",2)

			log.WithFields(log.Fields{"Key": string(keyval[0]),	}).Debug(funcstr+"->found key")
			log.WithFields(log.Fields{"Val": string(keyval[1]),	}).Debug(funcstr+"->found value")
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
	}

	if okinfo && len(infoparams[0]) > 0 {
		
		urlparam := infoparams[0]
		if(urlparam[:1] == "(") {
			entitiesData.Info = urlparam[1:len(urlparam)-1]
		} else {
			entitiesData.Info = urlparam
		}
	}

	return
}