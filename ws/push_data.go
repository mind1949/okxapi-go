package ws

import (
	"encoding/json"
)

// PushData Push data parameters
type PushData[Arg Channel, Data any] struct {
	// Successfully subscribed channel
	Arg Arg `json:"arg,omitempty"`
	// Subscribed data
	Data []Data `json:"data,omitempty"`

	rawMessage []byte
}

// GetRawMessage get raw websocket message
func (d PushData[A, D]) GetRawMessage() []byte {
	return d.rawMessage
}

// MapPushData convert input to output
func MapPushData[Arg Channel, Data any](input PushData[ChannelRawMessage, json.RawMessage]) (output PushData[Arg, Data], err error) {
	output = PushData[Arg, Data]{
		rawMessage: input.GetRawMessage(),
	}

	err = json.Unmarshal(input.Arg.RawMessage, &output.Arg)
	for _, e := range input.Data {
		var data Data
		err = json.Unmarshal(e, &data)
		if err != nil {
			return output, err
		}
		output.Data = append(output.Data, data)
	}
	return output, nil
}
