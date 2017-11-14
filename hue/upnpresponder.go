package hue

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"net"
	"text/template"
)

type upnpResponse struct {
	bytes.Buffer
}

// https://developers.meethue.com/documentation/changes-bridge-discovery
func newUpnpResponse(httpAddress string) (upnpResponse, error) {

	const templateText = "HTTP/1.1 200 OK\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"EXT:\r\n" +
		"CACHE-CONTROL: max-age=100\r\n" +
		"LOCATION: http://{{.}}/description.xml\r\n" +
		"SERVER: FreeRTOS/7.4.2 UPnP/1.0 IpBridge/1.10.0\r\n" +
		"hue-bridgeid: 001788FFFE09A206\r\n" +
		"ST: upnp:rootdevice\r\n" +
		"USN: uuid:2f402f80-da50-11e1-9b23-00178809a206::upnp:rootdevice\r\n\r\n"

	// Create the template object
	parsedTemplate, err := template.New("").Parse(templateText)
	if err != nil {
		log.WithField("error", err).Error("Hue UPNP Template creation failed")
		return upnpResponse{}, err
	}

	// Expand the template with our address
	var rc upnpResponse
	err = parsedTemplate.Execute(&rc, httpAddress)
	if err != nil {
		log.WithField("error", err).Error("Hue UPNP Template expansion failed")
		return upnpResponse{}, err
	}

	// Return the expanded template
	return rc, nil
}

// Send response to a UPNP M-SEARCH request
func (r upnpResponse) respond(c *net.UDPConn, address *net.UDPAddr) {

	l, err := c.WriteToUDP(r.Bytes(), address)
	if err != nil {
		log.WithFields(log.Fields{
			"Address": address.String(),
			"Error":   err}).
			Error("Hue API UPNP WriteToUDP failed")
		return
	}

	_ = l
	// TODO: Excessive Logging
	//log.WithFields(
	//	log.Fields{"address": address.String(), "length": l, "body": string(r.Bytes())}).
	//	Info("Hue UPNP Response Sent")
}

// Start listening to UPNP discovery address and respond to requests
func upnpResponder(httpAddress string) {

	// Initialise a responder
	response, err := newUpnpResponse(httpAddress)
	if err != nil {
		return
	}

	upnpDiscoveryAddress := &net.UDPAddr{
		IP:   net.IPv4(239, 255, 255, 250),
		Port: 1900}

	// Listen to the UPNP discovery address
	c, err := net.ListenMulticastUDP("udp4", nil, upnpDiscoveryAddress)
	if err != nil {
		log.WithField("Error", err).
			Error("Hue API UPNP ListenMulticastUDP failed")
		return
	}

	// Foreach received discovery packet
	for {
		// Read a packet
		b := make([]byte, 1024)
		n, src, err := c.ReadFromUDP(b)
		if err != nil {
			log.WithField("Error", err).
				Error("Hue API UPNP ReadFromUDP failed")
			continue
		}

		_ = n
		// TODO: Excessive Logging
		//log.WithFields(
		//	log.Fields{"source": src.String(), "body": string(b[:n])}).
		//	Info("Hue UPNP discovery packet received")

		// Respond to request on it's own go routine
		go response.respond(c, src)
	}
}

func StartUpnpResponder(addressAndPort string) {

	go upnpResponder(addressAndPort)
}
