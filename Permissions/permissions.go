package _permissions

import (  
	"encoding/xml"
	"io/ioutil"
	"sync"
	//"os"
	_apperrors "fbrest/Base/apperrors"
	log "github.com/sirupsen/logrus"
)

const PermissionKeyStr = "user permission key"

type PermissionType int

const(
	All PermissionType = 9
	None  = 0
	Read  = 1
	ReadWrite  = 2
)

type repository struct {
	permissions map[string]Permission
	mu    sync.RWMutex
}

//UserKey and UserPassword for AppNeeds and Session Token geb´neration
//DBUser and DBPassword for datanabase login
type Permission struct {
	UserKey string `xml:"userkey"`
	UserPassword string `xml:"userpassword"`
	DBUser string `xml:"dbuser"`
	DBPassword string `xml:"dbpassword"`
	Type PermissionType `xml:"type"`
}

type Permissions struct {
    Permission []Permission `xml:"Permission"`
}

func (r *repository) Add(userkey string,userpassword string, dbuser string, dbuserpassword string, ptype PermissionType) (ky Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var data Permission
	data.UserKey  = userkey
	data.UserPassword  = userpassword
	data.Type = ptype
	data.DBUser = dbuser
	data.DBPassword = dbuserpassword
	r.permissions[userkey] = data
	log.WithFields(log.Fields{"Key": userkey,	}).Debug("func Add ")	
	return data
}

func (r *repository) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()	
	delete(r.permissions, token)
	log.WithFields(log.Fields{PermissionKeyStr: token,	}).Debug("func Delete "+PermissionKeyStr)	
}

func (r *repository) Get(token string) (item Permission, err error ) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.permissions[token]
	if !ok {
		err = _apperrors.ErrPermissionKeyNotFound
		log.WithFields(log.Fields{PermissionKeyStr+" not found": token,	}).Debug("func Get "+PermissionKeyStr)	
		return item, err
	}
	log.WithFields(log.Fields{PermissionKeyStr+" found": token,	}).Debug("func Get "+PermissionKeyStr)	
	return item, nil
}


func (r *repository) Exists(token string) (bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.permissions[token]
	if !ok {
		return false
	}
	return true
}

var (
	pr *repository
)

func Repository() *repository {
	if pr == nil {
		pr = &repository {
			permissions: make(map[string]Permission),
		}
	}
	return pr
}


//Gets Permissiom from a Userkey and its password
//If permissions.DBUser for this key is empty than will be set to userkey, means userkey == DBLoginUser
//If permissions.DBPassword for this key is empty than will be set to userpassword
//if permissions.Type is empty it will be set to none
func GetPermissionFromRepository(userkey string, userpassword string) (perm Permission, err error ) {
	
	log.WithFields(log.Fields{PermissionKeyStr: userkey,	}).Debug("func GetPermissionFromRepository")	
	var rep = Repository() 
	var result,err1 = rep.Get(userkey)	 
	if(err1 != nil) {
		return result,err1
	}

	if(result.UserPassword != userpassword) {
		err1 = _apperrors.ErrUserOrPasswordWrong
		return result,err1
	}

	//if(n(result.Type) < 1) {
	//	result.Type = None
	//}

	if(len(result.DBUser) < 1) {
		result.DBUser = result.UserKey
	}

	if(len(result.DBPassword) < 1) {
		result.DBPassword = result.UserPassword
	}
	return result,err
}

func ReadPermissions(pfile string) {
	data, err := ioutil.ReadFile(pfile)
    if err != nil {		
		log.WithFields(log.Fields{"File reading error": err,	}).Error("func ReadPermissions")	
        return
    }
	xdata := &Permissions{}
	xml.Unmarshal(data,&xdata)
	var rep = Repository() 
	for _, xd := range xdata.Permission {
		if(len(xd.UserKey) > 0) {				
			if(!rep.Exists(xd.UserKey)) {
				var itm = rep.Add(xd.UserKey, xd.UserPassword, xd.DBUser, xd.DBPassword, xd.Type)		
				log.WithFields(log.Fields{"Added permission": "User key:"+itm.UserKey+" DBUser:"+itm.DBUser+" Permission:"+string(itm.Type),	}).Debug("func ReadPermissions")	
			}
		}
	}  
}

func WritePermissions(pfile string) {
	var data Permissions
	var dt Permission
	dt.UserKey = "superuser"
	dt.UserPassword = "su"
	dt.DBUser = "SYSDBA"
	dt.DBPassword = "masterkey"
	dt.Type = All
	data.Permission  = append(data.Permission,dt)

	dt.UserKey = "456"
	dt.UserPassword = ""
	dt.DBUser = ""
	dt.DBPassword = ""
	dt.Type = None
	data.Permission  = append(data.Permission,dt)

	dt.UserKey = "789"
	dt.UserPassword = ""
	dt.DBUser = ""
	dt.DBPassword = ""
	dt.Type = Read
	data.Permission  = append(data.Permission,dt)
	
	file, _ := xml.MarshalIndent(data, "", " ")
 
	_ = ioutil.WriteFile(pfile, file, 0644)
	
}