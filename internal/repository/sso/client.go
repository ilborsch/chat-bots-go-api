package sso

import (
	"github.com/ilborsch/sso-proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
)

//go:generate mockgen -destination=../../mocks/mock_auth_client.go -package=mocks github.com/ilborsch/sso-proto/gen/go/sso AuthClient

type ClientRepository struct {
	client sso.AuthClient
}

func NewRepository(host string, port int) *ClientRepository {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("could not dial grpc sso client on " + addr)
	}
	client := sso.NewAuthClient(conn)
	return &ClientRepository{
		client: client,
	}
}
