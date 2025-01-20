package handlers

import (
	"errors"
	"github.com/mrkucher83/hash-service/client/pkg/logger"
	"github.com/mrkucher83/hash-service/server/pkg/cryptographer"
	"github.com/mrkucher83/hash-service/server/pkg/pb"
	"io"
)

type Server struct {
	pb.UnimplementedStringHashServiceServer
}

func (s *Server) HashStrings(stream pb.StringHashService_HashStringsServer) error {
	for {
		// получение сообщения от клиента
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}

			logger.Error("failed unexpectedly while reading from stream: %v", err)
			return err
		}

		// хэшируем строки запроса
		hashes := cryptographer.Hash(req.Values)

		// отправляем ответ
		resp := &pb.HashArrayResponse{
			Hashes: hashes,
		}
		if err = stream.Send(resp); err != nil {
			logger.Error("failed to send response to the client: %v", err)
		}
		return err
	}
}
