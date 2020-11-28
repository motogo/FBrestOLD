package apis

import (
	"fbrest/Dokumento/config"
	"fbrest/Dokumento/models"
	_struct "fbrest/Dokumento/struct"
	"encoding/json"
	"net/http"
)

func GetAllLocations(response http.ResponseWriter, r *http.Request) {

	db, err := config.Koneksi()
	var Response _struct.ResponseData
	if err != nil {
		Response.Status = http.StatusInternalServerError
		Response.Message = err.Error()
		Response.Data = nil
		restponWithJson(response, http.StatusInternalServerError, Response)
	} else {
		_models := models.ModelGetData{DB:db}
		IsiData, err2 := _models.GetDataTableStandort()
		if err2 != nil {
			Response.Status = http.StatusInternalServerError
			Response.Message = err2.Error()
			Response.Data = nil
			restponWithJson(response, http.StatusInternalServerError, Response)

		} else {
			Response.Status = http.StatusOK
			Response.Message = "Sukses"
			Response.Data = IsiData
			restponWithJson(response, http.StatusOK, Response)

		}
	}

}

func restponWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}