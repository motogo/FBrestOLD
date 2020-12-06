package main

import (
	"github.com/gorilla/mux"
	"fbrest/Base/apis"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"fbrest/Base/config"
	_permissions "fbrest/Base/permissions"
)

func main(){
	customFormatter := new(log.TextFormatter)
	customFormatter.ForceColors = true
	customFormatter.TimestampFormat = "2006-01-02T15:04:05Z07:00"
	customFormatter.FullTimestamp = true	
	log.SetFormatter(customFormatter)
	argsWithProg := os.Args[1:]
	if(len(argsWithProg) < 0) {
	    log.Info(argsWithProg)
	}
	var port int = 4488
	var verbose int = 1
	var err error = nil

	for i, _ := range argsWithProg {
		var arg = string(argsWithProg[i])
	
		if(strings.EqualFold(arg, "-P")||strings.EqualFold(arg, "--Port")) {					
			port , err = strconv.Atoi(string(argsWithProg[i+1]))
		}
		if(strings.EqualFold(arg, "-V")||strings.EqualFold(arg, "--Verbose")) {					
			verbose , err = strconv.Atoi(string(argsWithProg[i+1]))
		}
		if(strings.EqualFold(arg, "-H")||strings.EqualFold(arg, "--Help")) {				
			log.Info(config.AppName+" " + config.Copyright);
			log.Info("PARAMS");
			log.Info("-h, --Help          -> This Help");
			log.Info("-p, --Port <port>   -> Port on wich the Server is listening for REST commands");
			log.Info("-v, --Verbose <num> -> 0 = Log's only errors");
			log.Info("                    -> 1 = Log's infos,errors (default)");
			log.Info("                    -> 2 = Debugmode log's everything :-)");
		}
	}
    if(err != nil){
		port = 4488
		verbose = 3
	}
	
	if verbose > 1 {
		log.Info("Logging set to level debug");
	} else if verbose == 1 {
		log.Info("Logging set to level info");
	} else {
		log.Info("Logging set to level error");
	}

	if verbose > 1 {
		log.SetLevel(log.DebugLevel)
	} else if verbose == 1 {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.ErrorLevel)	
	}

	//_permissions.WritePermissions()
	
	_permissions.ReadPermissions()

	router := mux.NewRouter()

	router.HandleFunc("/rest/get/{token}/{table}",apis.GetTableData).Methods("GET")	
	router.HandleFunc("/db/sql/get/{token}",apis.GetSQLData).Methods("GET")
	router.HandleFunc("/db/test/{token}",apis.TestDBOpenClose).Methods("GET")
	router.HandleFunc("/api/sql/test",apis.TestSQLResponse).Methods("GET")	
	router.HandleFunc("/api/token/get",apis.GetSessionKey).Methods("GET")
	router.HandleFunc("/api/token/delete/{token}",apis.DeleteSessionKey).Methods("GET")
	router.HandleFunc("/api/token/set/{token}",apis.SetSessionKey).Methods("GET")
	router.HandleFunc("/api/help",apis.GetHelp).Methods("GET")
	
	log.Info(config.AppName+" " + config.Copyright);
	log.Info(" ")
	log.Info("Version:"+config.Version);
	log.Info("Server listening for Firebird REST at port "+strconv.Itoa(port)+" ...")
	err = http.ListenAndServe(":"+strconv.Itoa(port),router)
	
	if err != nil {		
		log.WithFields(log.Fields{"Error": err.Error(),	}).Error("func main()")
	}
}