package utils

import (
	"net"
	"strings"
)

func GetIPFromAddr(addr net.Addr) string {
	a := addr.String()
	if idx := strings.Index(a, ":"); idx != -1 {
		return a[:idx]
	}

	return "localhost"
}
