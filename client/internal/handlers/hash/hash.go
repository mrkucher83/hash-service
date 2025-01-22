package hash

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"github.com/mrkucher83/hash-service/server/pkg/pb"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"time"
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
	grpcServerAddress := os.Getenv("GRPC_SERVER")
	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.ErrorCtx(r.Context(), "failed to connect to grpc server: %v", err)
		return
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			logger.Error("failed to close grpc connection: %v", err)
		}
	}(conn)

	grpcClient := pb.NewStringHashServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := grpcClient.HashStrings(ctx)
	if err != nil {
		logger.ErrorCtx(r.Context(), "failed to create grpc stream: %v", err)
		http.Error(w, "server error occurred", http.StatusInternalServerError)
		return
	}

	request := &pb.StringArrayRequest{Values: req.Params}
	if err = stream.Send(request); err != nil {
		logger.ErrorCtx(r.Context(), "failed to send a request to grpc server: %v", err)
		return
	}

	if err = stream.CloseSend(); err != nil {
		logger.ErrorCtx(r.Context(), "failed to close grpc stream: %v", err)
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
			logger.ErrorCtx(r.Context(), "failed to receive a response from grpc server: %v", err)
			http.Error(w, "server error occurred", http.StatusInternalServerError)
			return
		}
		hashArrays = append(hashArrays, resp.Hashes...)
	}
	fmt.Println(hashArrays)
}

func GetHashes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get hashes page"))
}
