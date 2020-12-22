package _functions

import(
	"strings"
)

func ReplaceAritmetics(str string) (string) {
	var s string = strings.Replace(str," .gt. "," > ",-1)
	s = strings.Replace(s," .lt. "," < ",-1)
	s = strings.Replace(s," .gte. "," >= ",-1)
	s = strings.Replace(s," .lte. "," <= ",-1)
	return s
}

func HasOperator(value string) bool {
	
	if(strings.HasPrefix(value,">")){
		 return true
	}
	if(strings.HasPrefix(value,"<")){
		return true
	}
	if(strings.HasPrefix(value,"=")){
		 return true
	}
	return false
}