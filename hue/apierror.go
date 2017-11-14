package hue

import "fmt"

type apierror struct {
	Error apierrordetail `json:"error"`
}

type apierrordetail struct {
	Type        int    `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

func newApiErrorUnauthorizedUser(address string) apierror {
	return apierror{
		apierrordetail{
			Type:        1,
			Address:     address,
			Description: "unauthorized user",
		},
	}
}

func newApiErrorBodyContainsInvalidJson(address string, err error) apierror {
	return apierror{
		apierrordetail{
			Type:        2,
			Address:     address,
			Description: fmt.Sprintf("Body Contains Invalid Json (%v)", err),
		},
	}
}

func newApiErrorResourceNotAvailable(address string) apierror {
	return apierror{
		apierrordetail{
			Type:        3,
			Address:     address,
			Description: "resource not available",
		},
	}
}

func newApiErrorMethodNotAvailable(address string) apierror {
	return apierror{
		apierrordetail{
			Type:        4,
			Address:     address,
			Description: "method not available for resource",
		},
	}
}

// The emulator uses this error (not part of Hue spec)
func newApiErrorNotImplemented(address string) apierror {
	return apierror{
		apierrordetail{
			Type:        50,
			Address:     address,
			Description: "Not Implemented",
		},
	}
}

func newApiErrorLinkButtonNotPressed(address string) apierror {
	return apierror{
		apierrordetail{
			Type:        101,
			Address:     address,
			Description: "link button not pressed",
		},
	}
}
