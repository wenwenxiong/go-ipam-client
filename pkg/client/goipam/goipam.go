package goipam

import (
	"context"
	"github.com/bufbuild/connect-go"
	v1 "github.com/metal-stack/go-ipam/api/v1"
	"github.com/metal-stack/go-ipam/api/v1/apiv1connect"
	"net/http"
)

type Client struct {
	client apiv1connect.IpamServiceClient
}

func NewGoipamClient(endponit string) Client {
	var r Client
	r.client = apiv1connect.NewIpamServiceClient(
		http.DefaultClient,
		endponit,
		connect.WithGRPC(),
	)
	return r
}

func (c Client) GetPrefix(cidr string) *v1.GetPrefixResponse {
	result, err := c.client.GetPrefix(context.Background(), connect.NewRequest(&v1.GetPrefixRequest{Cidr: cidr}))
	if err != nil {
		panic(err)
	}
	return result.Msg
}

func (c Client) CreatePrefix(cidr string) *v1.CreatePrefixResponse {
	result, err := c.client.CreatePrefix(context.Background(), connect.NewRequest(&v1.CreatePrefixRequest{Cidr: cidr}))
	if err != nil {
		panic(err)
	}
	return result.Msg
}

func (c Client) DeletePrefix(cidr string) *v1.DeletePrefixResponse {
	result, err := c.client.DeletePrefix(context.Background(), connect.NewRequest(&v1.DeletePrefixRequest{Cidr: cidr}))
	if err != nil {
		panic(err)
	}
	return result.Msg
}

func (c Client) AcquireIP(cidr string, ip string) *v1.AcquireIPResponse {
	result, err := c.client.AcquireIP(context.Background(), connect.NewRequest(&v1.AcquireIPRequest{PrefixCidr: cidr, Ip: &ip}))
	if err != nil {
		panic(err)
	}
	return result.Msg
}

func (c Client) ReleaseIP(cidr string, ip string) *v1.ReleaseIPResponse {
	result, err := c.client.ReleaseIP(context.Background(), connect.NewRequest(&v1.ReleaseIPRequest{PrefixCidr: cidr, Ip: ip}))
	if err != nil {
		panic(err)
	}
	return result.Msg
}
