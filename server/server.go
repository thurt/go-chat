// server
package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/thurt/go-chat/proto"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	//"github.com/golang/protobuf/proto"
)

const port = 10000

type chatServer struct {
	streams []pb.Chat_ConnectServer
}

func (s *chatServer) Connect(stream pb.Chat_ConnectServer) error {
	s.streams = append(s.streams, stream)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		go func() {
			for i, stream_client := range s.streams {
				if err := stream_client.Send(msg); err != nil {
					s.streams = append(s.streams[:i], s.streams[i+1:]...)
				}
			}
		}()
	}
}

func newServer() *chatServer {
	s := new(chatServer)
	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
