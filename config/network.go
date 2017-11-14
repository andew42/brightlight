package config

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/andew42/brightlight/hue"
	"net"
	"strconv"
	"strings"
	"sync"
)

var networkInterfaceInfo hue.NetworkInterfaceInfo
var networkInterfaceInfoError error
var networkInterfaceInfoOnce sync.Once

// Return an arbitrary network interface on which to set up server
func GetNetworkInterfaceInfo() (hue.NetworkInterfaceInfo, error) {

	// Discover interface info once and cache the results
	networkInterfaceInfoOnce.Do(func() {

		// Determine a 'good' ip address to use on this machine
		networkInterfaceInfo.Ip, networkInterfaceInfoError = getLocalIP()
		if networkInterfaceInfoError != nil {
			return
		}

		// Use port 80 because Hue can only be discovered on port 80
		// Note: port 80 forces us to run as root on linux and OSX
		networkInterfaceInfo.Port = ":80"

		// Get interface and address structure associated with IP address
		ourInterface, ourAddress, err := getInterfaceAndAddress(networkInterfaceInfo.Ip)
		ip := net.ParseIP(networkInterfaceInfo.Ip)

		// Extract network mask. The address should be of the form 8.8.8.8/24
		am := strings.Split(ourAddress.String(), "/")
		mask := 0
		if len(am) == 2 {
			if mask, err = strconv.Atoi(am[1]); err != nil {
				mask = 0
			}
		}
		if mask == 0 {
			// Use the default mask
			log.Warn("Failed to retrieve network mask, using default mask")
			networkInterfaceInfo.Mask = ip.DefaultMask().String()
		} else {
			networkInterfaceInfo.Mask = net.IP(net.CIDRMask(mask, 32)).String()
		}

		// Get the MAC address
		networkInterfaceInfo.Mac = ourInterface.HardwareAddr.String()

		// TODO: How should the gateway be retrieved? (for now we guess)
		gateway := ip.Mask(net.CIDRMask(mask, 32))
		gateway[3] |= 1
		networkInterfaceInfo.Gateway = gateway.String()
	})

	return networkInterfaceInfo, networkInterfaceInfoError
}

// A problem with this is that it fails if there is no router (e.g. no hot spot)
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

// Given a string IP address returns it's interface and address struct
func getInterfaceAndAddress(ip string) (net.Interface, net.Addr, error) {

	// Get all the machines interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, nil, err
	}

	// Foreach interface
	for _, i := range interfaces {

		// Get the interface's addresses
		interfaceAddresses, err := i.Addrs()
		if err != nil {
			continue
		}

		// Search interface address for our target address
		var foundAddress net.Addr = nil
		for _, addr := range interfaceAddresses {
			if strings.HasPrefix(addr.String(), ip) {
				foundAddress = addr
				break
			}
		}

		// If our address isn't found try next interface
		if foundAddress == nil {
			continue
		}

		// Here we have found the interface and address
		return i, foundAddress, nil
	}
	return net.Interface{}, nil, errors.New("Interface not found for " + ip)
}
