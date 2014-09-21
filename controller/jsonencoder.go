package controller

import "strconv"

// Convert led to JSON colour array
func (led Rgb) MarshalJSON() ([]byte, error) {
	rc := make([]byte, 0, 16)
	rc = strconv.AppendInt(rc, int64(led.Red)<<16+int64(led.Green)<<8+int64(led.Blue), 10)
	return rc, nil
}

