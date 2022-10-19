package rpc

import (
	"context"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
)

type (
	RPCClient struct {
		httpClient	http.Client
		rpcAddr		string
	}

	rpcError struct {
		Code	int64		`json:"code"`
		Message	string		`json:"message"`
		Details	[]string	`json:"details"`
	}
)

func NewRPCClient(rpcAddr string) *RPCClient {
	c := &RPCClient{
		httpClient:	http.Client{},
		rpcAddr:	rpcAddr,
	}

	return c
}

func (c *RPCClient) rpcRequestGet(ctx context.Context, option string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.rpcAddr+option, nil)
	if err != nil {
		klog.Errorf("failed to http request: %w", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		klog.Errorf("failed to get http response: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("failed to read body: %w", err)
		return nil, err
	}

	return body, nil
}
