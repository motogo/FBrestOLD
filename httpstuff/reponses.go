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
	}	
	return uee
}