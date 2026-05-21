package gamemanager

import (
	"errors"
	"fmt"
	"slices"
)

type IOConnection interface {
	Close() error
	SendJSON(message any) error
	SendRaw(message string) error
}

type IOManager interface {
	AddConnection(id string, conn IOConnection) error
	GetConnection(id string) IOConnection
	RemoveConnection(id string) error

	Broadcast(message string) error
	BroadcastJSON(message any) error
	BroadcastJSONExcept(message any, ids []string) error

	DisconnectAll() error
}

type ConnectionsManager struct {
	Connections map[string]*IOConnection
}

func NewConnectionsManager() *ConnectionsManager {
	return &ConnectionsManager{Connections: make(map[string]*IOConnection)}
}

func (cm *ConnectionsManager) AddConnection(id string, conn IOConnection) error {
	if cm.Connections[id] != nil {
		return errors.New("connection already exists")
	}
	cm.Connections[id] = &conn
	return nil
}

func (cm *ConnectionsManager) GetConnection(id string) IOConnection {
	return *cm.Connections[id]
}

func (cm *ConnectionsManager) RemoveConnection(id string) error {
	if cm.Connections[id] == nil {
		return errors.New("connection does not exist")
	}
	delete(cm.Connections, id)
	return nil
}

func (cm *ConnectionsManager) Broadcast(message string) error {
	for id, conn := range cm.Connections {
		err := (*conn).SendRaw(message)
		if err != nil {
			errMessage := fmt.Sprintf("error sending message to id %s: %s, error: %s", id, message, err)
			return errors.New(errMessage)
		}
	}
	return nil
}

func (cm *ConnectionsManager) BroadcastJSON(message any) error {
	for id, conn := range cm.Connections {
		err := (*conn).SendJSON(message)
		if err != nil {
			errMessage := fmt.Sprintf("error sending message to id %s: %s, error: %s", id, message, err)
			return errors.New(errMessage)
		}
	}
	return nil
}

func (cm *ConnectionsManager) BroadcastJSONExcept(message any, ids []string) error {
	for id, conn := range cm.Connections {
		if !slices.Contains(ids, id) {
			err := (*conn).SendJSON(message)
			if err != nil {
				errMessage := fmt.Sprintf("error sending message to id %s: %s, error: %s", id, message, err)
				return errors.New(errMessage)
			}
		}
	}
	return nil
}

func (cm *ConnectionsManager) DisconnectAll() error {
	for id, conn := range cm.Connections {
		err := (*conn).Close()
		if err != nil {
			errMessage := fmt.Sprintf("error disconnecting id %s, error: %s", id, err)
			return errors.New(errMessage)
		}
	}
	return nil
}
