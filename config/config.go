package config

import ("database/sql"
		"strconv"
		_sessions "fbrest/Base/sessions"
		//"null"
		_"github.com/nakagami/firebirdsql"
		log "github.com/sirupsen/logrus"
)



func ConnLocation(port int,location string, filename string, user string, password string) (db *sql.DB, err error) {
   
    if(port < 1){
	    port = 3050
	}
	if(len(location)) < 1{
	    location = "localhost"
	}

	if(len(user) < 1) {
		user = "SYSDBA"
	}

	if(len(password) < 1) {
		password = "masterkey"
	}

    var connstr = string(user+":"+password+"@"+location+":"+strconv.Itoa(port)+"/"+filename)	
	return ConnLocationWithString(connstr)
}

func ConnLocationWithString(connectionstring string) (db *sql.DB, err error) {
   
	log.WithFields(log.Fields{"Open database:": connectionstring,	}).Info("func ConnLocationWithString")    
	db, err = sql.Open("firebirdsql", connectionstring) 
	if err != nil {		
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("func ConnLocationWithString")
	}
	return
}
func ConnLocationWithSession(kv _sessions.Items) (db *sql.DB, err error) {
   
	log.WithFields(log.Fields{"Open database:": kv.Key,	}).Info("func ConnLocationWithSession")    
	db, err = sql.Open("firebirdsql", kv.Value) 
	if err != nil {		
		log.WithFields(log.Fields{"Open database, Error": err.Error(),	}).Error("func ConnLocationWithSession")
	}
	return
}



func TestConnLocation(connstr string) (err error) {
    
	log.WithFields(log.Fields{"Open database, connection:": connstr,	}).Info("func TestConnLocation")    
	
	var db *sql.DB
	db, err = sql.Open("firebirdsql", connstr) 

	if err != nil {
		return
	}

	err = db.Ping(); 
	
	return
}




