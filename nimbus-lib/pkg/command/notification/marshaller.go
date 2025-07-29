package notification

import (
	"encoding/json"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/nimbus-lib/pkg/message"
)

func MarshalEnvelopeFromCommand(cmd NotificationCommand) ([]byte, error) {
	var empty []byte

	payload, err := json.Marshal(cmd)
	if err != nil {
		return empty, fmt.Errorf("marshal command: %w", err)
	}

	env := message.Envelope{
		Type:    cmd.Type(),
		Payload: payload,
	}

	body, err := json.Marshal(env)
	if err != nil {
		return empty, fmt.Errorf("marshal envelope: %w", err)
	}

	return body, nil
}
