package ws

import "fmt"

// Request
type Request[T any] struct {
	// Unique identifier of the messageProvided by client.
	// It will be returned in response message for identifying the corresponding request.
	// A combination of case-sensitive alphanumerics, all numbers, or all letters of up to 32 characters.
	Id string `json:"id,omitempty"`
	// Operation
	Op string `json:"op,omitempty"`
	// Request Parameters
	Args []T `json:"args,omitempty"`
}

type validater interface {
	validate() error
}

func validate[T any](vs ...T) error {
	if len(vs) == 0 {
		return fmt.Errorf("params cann't be empty")
	}
	for i := range vs {
		v, ok := any(vs[i]).(validater)
		if !ok {
			continue
		}
		err := v.validate()
		if err != nil {
			return fmt.Errorf("args[%d] is invalid, err: %v", i, err)
		}
	}
	return nil
}
