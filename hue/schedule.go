package hue

// https://developers.meethue.com/documentation/schedules-api-0
type schedule struct {
	name        string
	description string
	command     command
	time        string
	// 1.2.1 localtime string
	// 1.2.1 status string
	// 1.3 autodelete 1.3
	// 1.12 recycle bool
}

type command struct {
	address string
	method  string
	body    string
}
