package websocket

import "fmt"

func ValidateParams(ppid, name, gameHash string) error {
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
