// filepath: d:\Projects\MakeshopPayment\master-console\backend\internal\domain\object\audit_log\ip_address.go
package object

import (
	"net"
	"strings"
)

type IPAddress string

func (ip IPAddress) IsValid() bool {
	return net.ParseIP(string(ip)) != nil
}

func (ip IPAddress) IsIPv4() bool {
	parsedIP := net.ParseIP(string(ip))
	return parsedIP != nil && parsedIP.To4() != nil
}

func (ip IPAddress) IsIPv6() bool {
	parsedIP := net.ParseIP(string(ip))
	return parsedIP != nil && parsedIP.To4() == nil
}

func (ip IPAddress) IsPrivate() bool {
	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}

	if parsedIP.To4() != nil {
		// 10.0.0.0/8
		if parsedIP[0] == 10 {
			return true
		}
		// 172.16.0.0/12
		if parsedIP[0] == 172 && parsedIP[1] >= 16 && parsedIP[1] <= 31 {
			return true
		}
		// 192.168.0.0/16
		if parsedIP[0] == 192 && parsedIP[1] == 168 {
			return true
		}
		// 127.0.0.0/8 (localhost)
		if parsedIP[0] == 127 {
			return true
		}
	}

	// IPv6 local addresses
	if strings.HasPrefix(string(ip), "::1") || strings.HasPrefix(string(ip), "fc00:") || strings.HasPrefix(string(ip), "fd00:") {
		return true
	}

	return false
}

func (ip IPAddress) IsLoopback() bool {
	parsedIP := net.ParseIP(string(ip))
	if parsedIP == nil {
		return false
	}
	return parsedIP.IsLoopback()
}

func (ip IPAddress) String() string {
	return string(ip)
}
