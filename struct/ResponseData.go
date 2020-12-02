package _struct


type ResponseData struct {
    Status int64 `json:”status”`
    Message string `json:”message”`
    Data interface{} `json:”data”`
}

type OutputData struct {
 ID string `json:"id"`
 BEZ string `json:"bez"`
 GUELTIG int64 `json:"gueltig"`
}

