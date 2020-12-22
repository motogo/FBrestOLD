package _apperrors
 
import "errors"
 
// ------------------------------------------------
 
var ErrPermissionKeyNotFound = errors.New("Permission key not found")
var ErrUserOrPasswordWrong = errors.New("User or password wrong !!!")
var ErrPermissionItemNotFound = errors.New("Permission item not found")
var ErrNoPermission = errors.New("No permission")
var DataNil = "No data"
var GotDatas = "Got datas"
var NoTableGiven = "No table given"
var NoFieldsGiven = "No fields given"


 
/*
 
type NarrowError struct {
	Message string
}
 
func NewNarrowError(m string) NarrowError {
	return NarrowError{Message:m}
}
 
func (n NarrowError) Error() string {
	return n.Message
}
 

 
type SmallError struct {
	Message string
}
 
func NewSmallError(m string) SmallError {
	return SmallError{Message:m}
}
 
func (s SmallError) Error() string {
	return s.Message
}

*/