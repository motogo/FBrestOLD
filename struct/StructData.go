package _struct

//import "time"




type SQLAttributes struct {
 
 Cmd string  `json:"CMD"`
 Sepfill string  `json:"sepfill"`
 Info string `json:"INFO"`
 Debug string `json:"Debug"`
}

type GetTABLEAttributes struct {	
	Table string `json:"table"` 
	Function string `json:"Function"` 
	FunctionParams string `json:"FunctionParams"` 
	Fields string `json:"fields"` 
	Filter string `json:"filter"` 
	OrderBy string `json:"orderby"` 
	GroupBy string `json:"groupby"` 	
	Sepfill string  `json:"sepfill"`
	Info string `json:"info"`
	Debug string `json:"debug"`
   }

   type DatabaseAttributes struct {
	
	Location string  `json:"location"`
	Database string  `json:"database"`
	Port int  `json:"port"`	
	Password string `json:"password"`
	User string `json:"user"`
	
   }

   type DBFilterAttributes struct {
	Field string `json:"Field"`
	Value string `json:"Value"`
   }
   type SelectDBAttributes struct {
	   Fields []string `json:"Fields"`
	   Filter []DBFilterAttributes `json:"Filter"`
   }

   type ParamFormatType string

const(
    Text ParamFormatType = "Text"
    Json = "Json"
)

   