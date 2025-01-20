package hash

import (
	"encoding/json"
	"fmt"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"io"
	"net/http"
)

type ReqBody struct {
	Params []string
}

func CreateHashes(w http.ResponseWriter, r *http.Request) {
	// reading request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.ErrorCtx(r.Context(), "failed to get data form request %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = r.Body.Close(); err != nil {
		logger.Warn("failed to close request body: ", err)
	}

	// parsing body params to struct
	var req ReqBody
	if err = json.Unmarshal(data, &req); err != nil {
		logger.Error("failed to parse request params %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Params) == 0 {
		http.Error(w, "received empty params", http.StatusBadRequest)
		return
	}
	fmt.Println(req.Params)

	// sending data via grpc streaming
}

func GetHashes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get hashes page"))
}
