package message

import (
	"encoding/json"

	"github.com/benu-cloud/benu-livestreaming-gst/pkg/pkgerrors"
)

func Unmarshal(bytes []byte) (GenericPayload, error) {
	genericMessage := &GenericMessage{}
	if err := json.Unmarshal(bytes, genericMessage); err != nil {
		return nil, pkgerrors.NewUnmarshalError(err)
	}
	switch genericMessage.Type {
	case SessionDescriptionMessage:
		payload := &SessionDescriptionPayload{}
		if err := json.Unmarshal(genericMessage.Payload, payload); err != nil {
			return nil, pkgerrors.NewUnmarshalError(err)
		}
		return payload, nil
	case IceCandidateMessage:
		payload := &IceCandidatePayload{}
		if err := json.Unmarshal(genericMessage.Payload, payload); err != nil {
			return nil, pkgerrors.NewUnmarshalError(err)
		}
		return payload, nil
	}
	return nil, pkgerrors.NewUnsupportedMessageTypeError(string(genericMessage.Type))
}

func Marshal(payload GenericPayload) ([]byte, error) {
	GenericMessage := &GenericMessage{}
	switch pay := payload.(type) {
	case SessionDescriptionPayload:
		GenericMessage.Type = SessionDescriptionMessage
		payloadBytes, err := json.Marshal(pay)
		if err != nil {
			return nil, pkgerrors.NewMarshalError(err)
		}
		GenericMessage.Payload = payloadBytes
		payloadBytesFinal, err := json.Marshal(GenericMessage)
		if err != nil {
			return nil, pkgerrors.NewMarshalError(err)
		}
		return payloadBytesFinal, nil
	case IceCandidatePayload:
		GenericMessage.Type = IceCandidateMessage
		payloadBytes, err := json.Marshal(pay)
		if err != nil {
			return nil, pkgerrors.NewMarshalError(err)
		}
		GenericMessage.Payload = payloadBytes
		payloadBytesFinal, err := json.Marshal(GenericMessage)
		if err != nil {
			return nil, pkgerrors.NewMarshalError(err)
		}
		return payloadBytesFinal, nil
	}
	return nil, pkgerrors.NewUnsupportedMessageTypeError("")
}