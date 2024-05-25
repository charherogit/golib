package ip

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"go-micro.dev/v4/metadata"
	"golib/config"
	"io"
	"net"
	"net/http"
	"strings"
)

const (
	Ip = "ip"
)

var outboundIP string

func OutboundIP() string {
	return outboundIP
}

func DetectOutboundIP() error {
	ip, err := GetOutBoundIP()
	if err != nil {
		return err
	}
	outboundIP = ip
	return nil
}

func GetOutBoundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", fmt.Errorf("get outbound ip error: %w", err)
	}
	defer conn.Close()
	if ua, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return ua.IP.String(), nil
	} else {
		return strings.Split(conn.LocalAddr().String(), ":")[0], nil
	}
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		return "", fmt.Errorf("failed to get public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析 JSON 响应，获取 IP 地址
	ip := gjson.GetBytes(body, "origin").String()
	return ip, nil
}

func GetMachineName() string {
	ipAddr, _ := GetOutBoundIP()
	publicIp, _ := GetPublicIP()
	return config.C.ServerName + "(" + ipAddr + "|" + publicIp + ")"
}

func GetIP(conn net.Conn) string {
	addr := conn.RemoteAddr().String()
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		fmt.Printf("addr err split: %s", split)
		return addr
	}
	return split[0]
}

func GetIPByCtx(ctx context.Context) string {
	if ip, ok := ctx.Value(Ip).(string); ok {
		return ip
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "unknown"
	}

	get, b := md.Get(Ip)
	if !b {
		return "unknown"
	}
	return get
}
