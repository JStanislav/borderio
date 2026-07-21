import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { allPlayersReady, getDefaultGameState, type GameState } from './game/GameState';
import { gameTimedOutId, startConnection  } from './server/server';
import toast, { Toaster } from 'react-hot-toast';
import { useNavigate, useParams, useSearchParams } from 'react-router';
import { send, gracefullyCloseConnection } from './server/server-conn';
import { DefaultLobby, type Lobby } from './game/lobby/lobby';
import type { MatchConfiguration } from './game/MatchConfiguration';
import { Lobby as LobbyComponent } from './components/lobby/Lobby.tsx';
import { useAuth } from './contexts/auth-provider.tsx';


export const LobbyContext = createContext<Lobby>(DefaultLobby);

function App() {
  const [gameState, setGameState] = useState<GameState>(getDefaultGameState())
  const [lobby, setLobby] = useState<Lobby>(DefaultLobby)
  const [matchConfiguration, setMatchConfiguration] = useState<MatchConfiguration>({ playerAmount: 2 })
  const navigate = useNavigate();

  const { id } = useParams()
  const [searchParams] = useSearchParams()

  const {user, setUser} = useAuth()
  
  if (!user) {
    return null;
  }

  useEffect(() => {
    if (id !== undefined && searchParams.get("action") !== null) {      
      const action = searchParams.get("action") as "create" | "join"

      startConnection(id, action, user.ppid, setGameState, setUser, setLobby, setMatchConfiguration, redirectToHome);
    }

    return () => {
      gracefullyCloseConnection("going away");
    }
  }, [])

  const redirectToHome = () => {
    toast.dismiss(gameTimedOutId)
    navigate("/");
  };

  const toggleReady = () => {
    const player = lobby.players.find(p => p.id === user.id);
    if (player === undefined) {
      console.error("Player not found in lobby");
      return;
    }
    const type = "playerReady";
    const data = {playerId: player.id, ppid: user.ppid, ready: !player.ready };
    send(type, data);
  }

  const onClickStartGame = () => {
    const type = "startGame";
    const data = {ppid: user.ppid};
    send(type, data);
  }


  return (
    <LobbyContext value={lobby}>
      {
        allPlayersReady(gameState) ?
        <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>
          <GameFrame gameState={gameState}/>
        </div>
        :
        <LobbyComponent players={lobby.players} matchConfiguration={matchConfiguration} actions={{toggleReady, onPlayerClickStartGame: onClickStartGame}} />
      }
      <button onClick={redirectToHome}>Leave Game</button>
      <Toaster />
    </LobbyContext>
  )
}

export default App
