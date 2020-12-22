package _httpstuff

import (
		
	"net/url"
	"net/http"
	"encoding/json"
	"strings"
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

func UnEscape(par string) (string) {
	if(strings.HasPrefix(par,"%22")) { 
		par = par[3:] 
	}
	if(strings.HasSuffix(par,"%22")) { 
		par = par[:len(par)-3] 
	}
	
	var uee,erruee = url.QueryUnescape(par)
	if(erruee != nil) {
	  uee = strings.Replace(par, "%22", "\"", -1)
	  uee = strings.Replace(uee, "%20", " ", -1)
	  uee = strings.Replace(uee, "%27", "'", -1)
	  uee = strings.Replace(uee, "%C3%BC", "ü",-1)
	  uee = strings.Replace(uee, "%3E", ">",-1)
	  uee = strings.Replace(uee, "%3C", "<",-1)
	  uee = strings.Replace(uee, "%24", "$",-1)
	  uee = strings.Replace(uee, "%26", "&",-1)
	  uee = strings.Replace(uee, "%60", "`",-1)
	  uee = strings.Replace(uee, "%3A", ":",-1)
	  uee = strings.Replace(uee, "%5B", "[",-1)
	  uee = strings.Replace(uee, "%5D", "]",-1)
	  uee = strings.Replace(uee, "%7B", "{",-1)
	  uee = strings.Replace(uee, "%7D", "}",-1)
	  uee = strings.Replace(uee, "%2B", "+",-1)
	  uee = strings.Replace(uee, "%23", "#",-1)
	  uee = strings.Replace(uee, "%40", "@",-1)	  
	  uee = strings.Replace(uee, "%2F", "/",-1)
	  uee = strings.Replace(uee, "%3B", ";",-1)
	  uee = strings.Replace(uee, "%3D", "=",-1)
	  uee = strings.Replace(uee, "%3F", "?",-1)
	  uee = strings.Replace(uee, "%5C", "\\",-1)
	  uee = strings.Replace(uee, "%5E", "^",-1)
	  uee = strings.Replace(uee, "%7C", "|",-1)
	  uee = strings.Replace(uee, "%7E", "~",-1)
	  uee = strings.Replace(uee, "%27", "´",-1)
	  uee = strings.Replace(uee, "%2C", ",",-1)

	  uee = strings.Replace(uee, "%25", "%",-1)
	}	
	return uee
}

