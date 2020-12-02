package _sessions

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const MaxDuration = 30*60*1E9  //ns
const SessionKeyStr = "session key"

type Items struct {
	Key string  `json:"Key"` 
	Value string  `json:"Value"` 
	Start time.Time  `json:"Time"` 
	Duration time.Duration  `json:"Duration"` 
	Valid bool `json:"Valid"` 
   }

type repository struct {
	items map[string]Items
	mu    sync.RWMutex
}

func (r *repository) Add(key, conn string) (ky Items) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var data Items
	data.Key = key
	data.Value = conn
	data.Start = time.Now()
	data.Duration = MaxDuration
	r.items[key] = data
	log.WithFields(log.Fields{SessionKeyStr: key,	}).Debug("func Add "+SessionKeyStr)	
	return data
}

func (r *repository) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()	
	delete(r.items, key)
	log.WithFields(log.Fields{SessionKeyStr: key,	}).Debug("func Delete "+SessionKeyStr)	
}

func (r *repository) Get(key string) (Items, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[key]
	if !ok {
		
		var err error
		log.WithFields(log.Fields{SessionKeyStr+" not found": key,	}).Debug("func Get "+SessionKeyStr)	
		return item,err
	}
	log.WithFields(log.Fields{SessionKeyStr+" found": key,	}).Debug("func Get "+SessionKeyStr)	
	return item, nil
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



