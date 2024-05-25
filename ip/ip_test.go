package ip

import (
	"context"
	"fmt"
	"go-micro.dev/v4/metadata"
	"testing"
)

func TestOutboundIP(t *testing.T) {
	err := DetectOutboundIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(OutboundIP())
}

func TestGetIP(t *testing.T) {
	reqIp := "192.168.196.19"
	ctx := context.WithValue(context.Background(), Ip, reqIp)

	ip := GetIPByCtx(ctx)
	fmt.Printf("ip: %v\n", ip)

	ctxMetadata := metadata.NewContext(context.Background(), map[string]string{
		Ip: reqIp,
	})

	ip = GetIPByCtx(ctxMetadata)
	fmt.Printf("ip: %v\n", ip)
}

func TestGetOutBoundIP(t *testing.T) {
	got, err := GetOutBoundIP()
	if err != nil {
		t.Errorf("GetOutBoundIP() error = %v", err)
		return
	}
	fmt.Printf("got: %v\n", got)
}

func TestGetPublicIP(t *testing.T) {
	got, err := GetPublicIP()
	if err != nil {
		t.Errorf("GetPublicIP() error = %v", err)
		return
	}
	fmt.Printf("got: %v\n", got)
}
