package hue

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"text/template"
)

// Serves /description.xml for Hue discovery
// TODO: Point to correct icons
func GetUpnpDescriptionHandler(httpAddress string) (func(http.ResponseWriter, *http.Request), error) {

	// https://developers.meethue.com/documentation/hue-bridge-discovery
	// https://developers.meethue.com/documentation/bridge-v2
	const descriptionTemplateText = `<?xml version="1.0"?>
		<root xmlns="urn:schemas-upnp-org:device-1-0">
		   <specVersion>
			  <major>1</major>
			  <minor>0</minor>
		   </specVersion>
		   <URLBase>http://{{.}}/</URLBase>
		   <device>
			  <deviceType>urn:schemas-upnp-org:device:Basic:1</deviceType>
			  <friendlyName>Philips hue ({{.}})</friendlyName>
			  <manufacturer>Royal Philips Electronics</manufacturer>
			  <manufacturerURL>http://www.philips.com</manufacturerURL>
			  <modelDescription>Philips hue Personal Wireless Lighting</modelDescription>
			  <modelName>Philips hue bridge 2015</modelName>
			  <modelNumber>BSB002</modelNumber>
			  <modelURL>http://www.meethue.com</modelURL>
			  <serialNumber>001788102201</serialNumber>
			  <UDN>uuid:2f402f80-da50-11e1-9b23-001788102201</UDN>
			  <presentationURL>index.html</presentationURL>
			  <iconList>
				 <icon>
					<mimetype>image/png</mimetype>
					<height>48</height>
					<width>48</width>
					<depth>24</depth>
					<url>hue_logo_0.png</url>
				</icon>
				<icon>
				   <mimetype>image/png</mimetype>
				   <height>120</height>
				   <width>120</width>
				   <depth>24</depth>
				   <url>hue_logo_3.png</url>
				</icon>
			  </iconList>
		   </device>
		</root>`

	// Create template
	var err error
	descriptionTemplate, err := template.New("").Parse(descriptionTemplateText)
	if err != nil {
		return nil, err
	}

	// Expand template with address
	b := &bytes.Buffer{}
	err = descriptionTemplate.Execute(b, httpAddress)
	if err != nil {
		return nil, err
	}

	// The expanded response text
	descriptionBody := b.Bytes()

	// Return response handler function
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		l, err := w.Write(descriptionBody)
		if err != nil || l != len(descriptionBody) {
			log.WithField("Error", err).
				Error("Hue API UPNP failed to serve description.xml")
		}
	}, nil
}
