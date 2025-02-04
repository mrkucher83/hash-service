package hash

import (
	"context"
	"encoding/json"
	"github.com/mrkucher83/hash-service/client/internal/godb"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"github.com/mrkucher83/hash-service/server/pkg/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"time"
)

func (hr *Repo) CreateHashes(w http.ResponseWriter, r *http.Request) {
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

	// sending data via grpc streaming
	hashArrays, err := sendToGRPC(r.Context(), req.Params)
	if err != nil {
		http.Error(w, "server error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// storing hashes into DB
	response, err := saveHashesToDB(r.Context(), hr.storage, hashArrays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode response: %v", err)
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendToGRPC(rCtx context.Context, params []string) ([]string, error) {
	// sending data via grpc streaming
	grpcServerAddress := os.Getenv("GRPC_SERVER")
	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.ErrorCtx(rCtx, "failed to connect to grpc server: %v", err)
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			logger.Error("failed to close grpc connection: %v", err)
		}
	}(conn)

	grpcClient := pb.NewStringHashServiceClient(conn)
	ctx, cancel := context.WithTimeout(rCtx, time.Second*10)
	defer cancel()

	stream, err := grpcClient.HashStrings(ctx)
	if err != nil {
		logger.ErrorCtx(rCtx, "failed to create grpc stream: %v", err)
		return nil, err
	}

	request := &pb.StringArrayRequest{Values: params}
	if err = stream.Send(request); err != nil {
		logger.ErrorCtx(rCtx, "failed to send a request to grpc server: %v", err)
		return nil, err
	}

	if err = stream.CloseSend(); err != nil {
		logger.ErrorCtx(rCtx, "failed to close grpc stream: %v", err)
	}

	// receiving a response from grpc streaming
	var hashArrays []string
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Info("response received successfully from grpc stream")
				break
			}
			logger.ErrorCtx(rCtx, "failed to receive a response from grpc server: %v", err)
			return nil, err
		}
		hashArrays = append(hashArrays, resp.Hashes...)
	}

	return hashArrays, nil
}

func saveHashesToDB(rCtx context.Context, storage *godb.Instance, data []string) ([]godb.Resp, error) {
	var records []godb.Resp
	for _, hash := range data {
		rec := &godb.Record{Text: hash}

		res, err := storage.AddRecord(rCtx, rec)
		if err != nil {
			return nil, err
		}

		records = append(records, *res)
	}

	return records, nil
}
