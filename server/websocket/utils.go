package websocket

import "fmt"

func ValidateParams(action, ppid, name, gameHash string) error {
	if !(action == "create" || action == "join" || action == "spectate") {
		if action == "" {
			return fmt.Errorf("missing action")
		}

		return fmt.Errorf("invalid action: %s", action)
	}

	if ppid == "" {
		return fmt.Errorf("missing ppid")
	}

	if name == "" {
		return fmt.Errorf("missing name")
	}

	if gameHash == "" {
		return fmt.Errorf("missing game hash")
	}

	return nil
}
