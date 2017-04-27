package config

import (
	"net"
	"strings"
)

// Get preferred outbound ip address of this machine TODO: Multiple network adapters
// The problem with this is that it fails if there is no router (e.g. no hot spot)
// http://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func getLocalIP() (string, error) {

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

// Use port 80 because Hue can only be discovered on port 80
// Note: port 80 forces us to run as root on linux and OSX
func GetServerPort() string {
	return ":80"
}

// e.g. 172.20.10.2:80
func GetServerAddressAndPort() (string, error) {
	address, err := getLocalIP()
	if err != nil {
		return "", err
	}
	return address + GetServerPort(), nil
}
