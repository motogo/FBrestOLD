package _permissions

import (  
	"encoding/xml"
	"io/ioutil"
	"sync"
	log "github.com/sirupsen/logrus"
)

const PermissionKeyStr = "permission key"


type PermissionType string

const(
	All PermissionType = "All"
	None  = "None"
	Read  = "Read"
	ReadWrite  = "ReadWrite"
)

type prepository struct {
	permissions map[string]Permission
	mu    sync.RWMutex
}

type Permission struct {
	Key string `xml:",key"`
	Type PermissionType `xml:",type"`
}

type Permissions struct {
    Permission []Permission `xml:"Permission"`
}

func (r *prepository) Add(pkey string, ptype PermissionType) (ky Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var data Permission
	data.Key  = pkey
	data.Type = ptype
	
	r.permissions[pkey] = data
	log.WithFields(log.Fields{"Key": pkey,	}).Debug("func Add ")	
	return data
}

func (r *prepository) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()	
	delete(r.permissions, token)
	log.WithFields(log.Fields{PermissionKeyStr: token,	}).Debug("func Delete "+PermissionKeyStr)	
}

func (r *prepository) Get(token string) (Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.permissions[token]
	if !ok {
		
		var err error
		log.WithFields(log.Fields{PermissionKeyStr+" not found": token,	}).Debug("func Get "+PermissionKeyStr)	
		return item,err
	}
	log.WithFields(log.Fields{PermissionKeyStr+" found": token,	}).Debug("func Get "+PermissionKeyStr)	
	return item, nil
}




var (
	pr *prepository
)

func Repository() *prepository {
	if pr == nil {
		pr = &prepository {
			permissions: make(map[string]Permission),
		}
	}
	return pr
}

func GetPermissionFromRepository(permission string) (perm Permission) {
	
	log.WithFields(log.Fields{PermissionKeyStr: permission,	}).Debug("func GetPermissionFromRepository")	
	var rep = Repository() 
	var result,_ = rep.Get(permission)	 
	result.Key = permission
	if(len(result.Type) < 1) {
		result.Type = None
	}
	return result 
}

func ReadPermissions() {
data, err := ioutil.ReadFile("permissions.xml")
    if err != nil {		
		log.WithFields(log.Fields{"File reading error": err,	}).Error("func ReadPermissions")	
        return
    }
	xdata := &Permissions{}
	xml.Unmarshal(data,&xdata)
	var rep = Repository() 
	for _, xd := range xdata.Permission {
		var itm = rep.Add(xd.Key,xd.Type)		
		log.WithFields(log.Fields{"Added permission": itm,	}).Debug("func ReadPermissions")	
	}  
}

func WritePermissions() {
	var data Permissions
	var dt Permission
	dt.Key = "123"
	dt.Type = All
	data.Permission  = append(data.Permission,dt)

	dt.Key = "456"
	dt.Type = None
	data.Permission  = append(data.Permission,dt)

	dt.Key = "789"
	dt.Type = Read
	data.Permission  = append(data.Permission,dt)
	
	file, _ := xml.MarshalIndent(data, "", " ")
 
	_ = ioutil.WriteFile("permissions.xml", file, 0644)
	
}