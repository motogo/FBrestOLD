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
	
)

func main(){
	customFormatter := new(log.TextFormatter)
	customFormatter.ForceColors = true
	customFormatter.TimestampFormat = "2020-01-02 15:04:05"
	customFormatter.FullTimestamp = true	
	log.SetFormatter(customFormatter)
	argsWithProg := os.Args[1:]
	if(len(argsWithProg) < 0){
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

	router := mux.NewRouter()

	router.HandleFunc("/rest/get/{key}/{table}",apis.GetTableData).Methods("GET")
	//router.HandleFunc("/rest/{key}/{table}/{field}",apis.GetTableData).Methods("GET")
	router.HandleFunc("/db/sql/get/{key}",apis.GetSQLData).Methods("GET")
	router.HandleFunc("/db/test/{key}",apis.TestDBOpenClose).Methods("GET")
	router.HandleFunc("/api/test/{key}",apis.TestResponse).Methods("GET")
	router.HandleFunc("/api/getkey",apis.GetSessionKey).Methods("GET")
	router.HandleFunc("/api/deletekey/{key}",apis.DeleteSessionKey).Methods("GET")
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