package ws

import "encoding/json"

var _ error = responseErr[Channel]{}

type responseErr[T Channel] struct {
	Response[T]
}

func (r responseErr[T]) Error() string {
	if !r.isError() {
		return ""
	}
	b, _ := json.Marshal(r)
	return string(b)
}

// Response Operation result
type Response[T Channel] struct {
	// Operation
	// login
	// subscribe
	// unsubscribe
	// error
	Event string `json:"event"`

	Arg T `json:"arg"`
	// Error code
	Code string `json:"code"`
	// Error message
	Msg string `json:"msg"`
	// WebSocket connection ID
	ConnId string `json:"connId"`

	rawMessage []byte
}

func (r Response[T]) GetError() error {
	if r.isError() {
		return responseErr[T]{r}
	}
	return nil
}

// isError
func (r Response[T]) isError() bool {
	switch r.Code {
	case "0", "":
		return false
	default:
		return true
	}
}

// GetRawMessage get raw websocket message
func (r Response[T]) GetRawMessage() []byte {
	return r.rawMessage
}

// MapResponse map input to output
func MapResponse[T Channel](input Response[ChannelRawMessage]) (output Response[T], err error) {
	output = Response[T]{
		Event:  input.Event,
		Code:   input.Code,
		Msg:    input.Msg,
		ConnId: input.ConnId,
	}
	return output, json.Unmarshal(input.Arg.RawMessage, &output.Arg)
}
