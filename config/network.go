package config

import (
	"net"
	"strings"
)

// A problem with this is that it fails if there is no router (e.g. no hot spot)
// http://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func GetLocalIP() (string, error) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Determine address that would be used for connection but
	// doesn't actually make any connection (as it's UDP)
	localAddress := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddress, ":")
	return localAddress[0:idx], nil
}
