package lobby

type Lobbyist struct {
	ID        int
	Name      string
	Ready     bool
	Host      bool
	Connected bool
	Spectator bool
}

func (l *Lobbyist) ToggleReady() {
	l.Ready = !l.Ready
}

// Generate a new Lobbyist with the given parameters
func NewLobbyist(id int, name string) *Lobbyist {
	return &Lobbyist{
		ID:        id,
		Name:      name,
		Ready:     false,
		Host:      false,
		Connected: false,
		Spectator: false,
	}
}
