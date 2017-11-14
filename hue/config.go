package hue

import (
	"encoding/json"
)

// HTTP Handler for /api/config
func processConfigRequest(c *cmdContext) {

	if c.Method == "GET" {
		if len(c.Resource) == 0 {
			respondWithJsonEncodedObject(c.W, &c.FullState.Config)
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	} else if c.Method == "PUT" {
		if len(c.Resource) == 0 {
			setConfig(c)
			return
		}
		reportError(c.W, newApiErrorResourceNotAvailable(c.ErrorAddress))
		return
	}
	reportError(c.W, newApiErrorMethodNotAvailable(c.ErrorAddress))
}

// https://developers.meethue.com/documentation/configuration-api
type config struct {
	Name             string                    `json:"name"`
	SwUpdate         interface{}               `json:"swupdate,omitempty"`
	WhiteList        map[string]whitelistEntry `json:"whitelist"`
	ApiVersion       string                    `json:"apiversion"`
	SwVersion        string                    `json:"swversion"`
	ProxyAddress     string                    `json:"proxyaddress"`
	ProxyPort        uint16                    `json:"proxyport"`
	LinkButton       bool                      `json:"linkbutton"`
	IpAddress        string                    `json:"ipaddress"`
	Mac              string                    `json:"mac"`
	NetMask          string                    `json:"netmask"`
	Gateway          string                    `json:"gateway"`
	Dhcp             bool                      `json:"dhcp"`
	PortalServices   bool                      `json:"portalservices"`
	Utc              string                    `json:"UTC"`
	LocalTime        string                    `json:"localtime"`
	TimeZone         string                    `json:"timezone"`
	ZigbeeChannel    uint16                    `json:"zigbeechannel"`
	ModelId          string                    `json:"modelid"`
	BridgeId         string                    `json:"bridgeid"`
	FactoryNew       bool                      `json:"factorynew"`
	ReplacesBridgeId interface{}               `json:"replacesbridgeid"`
	DatastoreVersion string                    `json:"datastoreversion"`
	StarterKitId     string                    `json:"starterkitid"`
}

type whitelistEntry struct {
	LastUseDate string `json:"last use date"`
	CreateDate  string `json:"create date"`
	Name        string `json:"name"`
}

func newConfig() *config {

	return &config{
		Name:             "Philips hue",
		SwVersion:        "9999999999", // TODO
		ApiVersion:       "1.20.0",     // TODO
		ProxyAddress:     "none",
		ProxyPort:        0,
		LinkButton:       true,
		PortalServices:   false,
		TimeZone:         "Europe/London",
		ZigbeeChannel:    25,
		ModelId:          "BSB002",
		BridgeId:         "001788FFFE2B9495", //TODO
		FactoryNew:       false,
		DatastoreVersion: "61", // TODO
		WhiteList:        make(map[string]whitelistEntry),
	}
}

func (c *config) setNetworkInfo(nii NetworkInterfaceInfo) {
	c.IpAddress = nii.Ip
	c.Mac = nii.Mac
	c.NetMask = nii.Mask
	c.Gateway = nii.Gateway
	c.Dhcp = true // TODO
}

// Body={"UTC":"2017-06-16T06:44:37"} Header=map[Accept:[*/*] Content-Type:[application/json] Content-Length:[29]]
// Method=PUT Url=/api/4d65822107fcfd5278629a0f5f3f164f/config
func setConfig(c *cmdContext) {

	// Decode body as a config
	var bodyConfig config
	if err := json.Unmarshal(c.Body, &bodyConfig); err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	// Decode the body as a map so we can determine what was set
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(c.Body, &bodyMap); err != nil {
		reportError(c.W, newApiErrorBodyContainsInvalidJson(c.ErrorAddress, err))
		return
	}

	cf := c.FullState.Config
	response := apiResponseList{}

	// Update config fields from body
	if _, ok := bodyMap["UTC"]; ok {
		cf.Utc = bodyConfig.Utc
		response = response.AppendSuccessResponse(c.ResourceUrl+"/UTC", cf.Utc)
	}

	// Return response
	respondWithJsonEncodedObject(c.W, response)
}
