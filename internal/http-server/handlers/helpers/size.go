package helpers

import (
	"encoding/json"
	"fmt"
	resp "github.com/AlexandrLitkevich/qwery/internal/lib/api/response"
	"log"
	"net/http"
)

//Протестить хандлер

func SizeHandlerHelpers(w http.ResponseWriter, r *http.Request) {
	var req any
	err := json.NewDecoder(r.Body).Decode(&req)
	var res resp.Response
	if err != nil {
		errFmt := fmt.Sprintf("failed to decode request %v", err)
		res = resp.Response{Error: errFmt, Status: "error"}
		w.WriteHeader(http.StatusBadRequest)

	} else {
		res = resp.Response{Error: "successful request", Status: "success"}
		w.WriteHeader(http.StatusOK)
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
