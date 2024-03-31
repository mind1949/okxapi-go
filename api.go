package okxapigo

import (
	"fmt"
)

// Api okx api
type Api struct {
	// API Key
	Key string
	// API Key password
	Passphrase string
	// API Secret key
	Secretkey string
}

// Validate validate api
func (a Api) Validate() error {
	if a.Key == "" {
		return fmt.Errorf("api key cann't be empty")
	}
	if a.Passphrase == "" {
		return fmt.Errorf("passphrase cann't be empty")
	}
	if a.Secretkey == "" {
		return fmt.Errorf("secretkey cann't be empty")
	}
	return nil
}
