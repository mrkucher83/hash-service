package hash

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/mrkucher83/hash-service/client/internal/godb"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (hr *Repo) GetHashes(w http.ResponseWriter, r *http.Request) {
	idsParams := chi.URLParam(r, "ids")
	ids := strings.Split(idsParams, ",")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var response []godb.Resp
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("failed to convert string id to integer %v", err)
		}
		hash, err := hr.storage.GetRecordById(ctx, idInt)
		if err != nil {
			logger.Error("failed to get hash by id from DB: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response = append(response, *hash)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode response: %v", err)
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
