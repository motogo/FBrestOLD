package _sessions

import (
	log "github.com/sirupsen/logrus"
	_httpstuff "fbrest/Base/httpstuff"
	_permissions "fbrest/Base/permissions"
	"sync"
	"time"
	"strconv"
	"net/http"
	_struct "fbrest/Base/struct"		
)




const MaxDuration = 30*60*1E9  //ns
const SessionTokenStr = "session token"

type Items struct {
	Token string  `json:"Token"` 
	Value string  `json:"Value"` 
	Permission _permissions.PermissionType `json:"Permission"`
	Start time.Time  `json:"Time"` 
	Duration time.Duration  `json:"Duration"` 
	Valid bool `json:"Valid"` 
   }

type repository struct {
	items map[string]Items
	mu    sync.RWMutex
}

func (r *repository) Add(token string, permission _permissions.PermissionType, conn string) (ky Items) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var data Items
	data.Token = token
	data.Value = conn
	data.Start = time.Now()
	data.Duration = MaxDuration
	data.Permission = permission
	r.items[token] = data
	log.WithFields(log.Fields{SessionTokenStr: token,	}).Debug("func Add "+SessionTokenStr)	
	return data
}

func (r *repository) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()	
	delete(r.items, token)
	log.WithFields(log.Fields{SessionTokenStr: token,	}).Debug("func Delete "+SessionTokenStr)	
}

func (r *repository) Get(token string) (Items, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[token]
	if !ok {
		
		var err error
		log.WithFields(log.Fields{SessionTokenStr+" not found": token,	}).Debug("func Get "+SessionTokenStr)	
		return item,err
	}
	log.WithFields(log.Fields{SessionTokenStr+" found": token,	}).Debug("func Get "+SessionTokenStr)	
	return item, nil
}

func GetTokenDataFromRepository(token string) (kval Items) {
	
	log.WithFields(log.Fields{SessionTokenStr: token,	}).Debug("func GetSessionKeyFromRepository")	
	var rep = Repository() 
	var result,_ = rep.Get(token)	 
	result.Token = token
	return result 
}

func TokenValid(response http.ResponseWriter, key string) (kv Items) {
	var Response _struct.ResponseData
	kv  = GetTokenDataFromRepository(key)	
	if(len(kv.Value) < 1) {
		Response.Status = http.StatusForbidden
		Response.Message = "No valid database connection found by "+SessionTokenStr+" "+kv.Token
		Response.Data = nil
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false
		return kv
	}
	var duration time.Duration = kv.Duration
	var end = time.Now();
	difference := end.Sub(kv.Start)
	if(difference > duration) {
		//Zeit für Key abgelaufen
		Response.Status = http.StatusForbidden
		Response.Message = SessionTokenStr+" "+kv.Token + " has expired after "+ strconv.Itoa(MaxDuration/1E9) +" seconds"
		Response.Data = nil
		var rep = Repository() 
	    rep.Delete(key)
		_httpstuff.RestponWithJson(response, http.StatusInternalServerError, Response)
		kv.Valid = false
		return kv
	}
	log.WithFields(log.Fields{"Remaining "+SessionTokenStr+" duration (s)": (duration-difference),	}).Debug("func KeyValid")
	kv.Valid = true
	return kv
}



	
	var (
		r *repository
	)
	
	func Repository() *repository {
		if r == nil {
			r = &repository {
				items: make(map[string]Items),
			}
		}
		return r
	}



