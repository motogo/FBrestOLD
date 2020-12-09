package _struct

import "time"

type FieldValueAttributes struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type GetUrlSQLAttributes struct {
 
	Cmd string  `json:"cmd"`
	Info string `json:"info"`	
   }

type SQLAttributes struct {
 
 Cmd string  `json:"cmd"`
 Sepfill string  `json:"sepfill"`
 Info string `json:"info"`
 Debug string `json:"Debug"`
}

type FIELDVALUEAttributes struct {	
	Table string `json:"table"` 
	Filter string `json:"filter"` 
	FieldValue []string `json:"fieldvalue"` 
	Debug string `json:"debug"`
   }

type GetUrlTABLEAttributes struct {	
	Table string `json:"table"` 		
	Fields []string `json:"fields"` 
	Filter string `json:"filter"` 
	OrderBy []string `json:"order"` 
	GroupBy []string `json:"group"`	
	First int `json:"first"` 	
	Skip int `json:"skip"` 	
   }

  
type GetUrlPayloadAttributes struct {	
	
	Payload []string `json:"payload"` 
	Filter string `json:"filter"` 	
   }

type GetTABLEAttributes struct {	
	Table string `json:"table"` 
	Function string `json:"function"` 
	FunctionParams string `json:"functionparams"` 
	Fields string `json:"fields"` 
	Filter string `json:"filter"` 
	OrderBy string `json:"orderby"` 
	GroupBy string `json:"groupby"` 	
	First int `json:"first"` 	
	Skip int `json:"skip"` 	
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


   type SelectDBAttributes struct {
	   Fields []string `json:"Fields"`
	   Filter []FieldValueAttributes `json:"Filter"`
   }

type ParamFormatType string
type SepFillType string

const(
    Text ParamFormatType = "Text"
    Json = "Json"
)

const(
    SepFillApos SepFillType = "'"
    SepFillMinus = "-"
)

type Profile struct {
	Appname    string
	Version string
	Copyright string
	Key string
	Duration time.Duration
  }

   