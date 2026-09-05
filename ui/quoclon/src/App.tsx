import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { allPlayersReady, getDefaultGameState, getPlayerById, type GameState } from './game/GameState';
import { gameTimedOutId, requestClickStartGame, requestToggleReady, startConnection  } from './server/server';
import toast, { Toaster } from 'react-hot-toast';
import { useNavigate, useParams, useSearchParams } from 'react-router';
import { gracefullyCloseConnection } from './server/server-conn';
import { DefaultLobby, type Lobby } from './game/lobby/lobby';
import type { MatchConfiguration } from './game/MatchConfiguration';
import { Lobby as LobbyComponent } from './components/lobby/Lobby.tsx';
import { useAuth } from './contexts/auth-provider.tsx';
import { GameOverDialog } from './components/game/GameOverDialog.tsx';


export const LobbyContext = createContext<Lobby>(DefaultLobby);

function App() {
  const [gameState, setGameState] = useState<GameState>(getDefaultGameState())
  const [lobby, setLobby] = useState<Lobby>(DefaultLobby)
  const [matchConfiguration, setMatchConfiguration] = useState<MatchConfiguration>({ playerAmount: 2 })
  const [winnerPlayerName, setWinnerPlayerName] = useState("Unknown");
  
  const navigate = useNavigate();

  const { id } = useParams()
  const [searchParams] = useSearchParams()

  const {user, setUser} = useAuth()
  
  useEffect(() => {
    if (id !== undefined && searchParams.get("action") !== null) {      
      const action = searchParams.get("action") as "create" | "join"

      if (user) {
        startConnection(id, action, user.ppid, user.name, setGameState, setUser, setLobby, setMatchConfiguration, redirectToHome);
      }
    }

    return () => {
      gracefullyCloseConnection("going away");
    }
  }, [user?.ppid])

  
  useEffect(() => {
      // Is there a winner?
      if (lobby.winnerPlayerId !== undefined) {
          const winnerPlayerName = getPlayerById(gameState, lobby.winnerPlayerId)?.name || "Unknown";
          setWinnerPlayerName(winnerPlayerName);
          onGameOver();
      }
      return () => {
        // set again to default, covering the case of leaving the game even if it wasn't game over
        if (user) {
          setUser({...user, inGame: false})
        }
      }
  }, [lobby.winnerPlayerId])

  const redirectToHome = () => {
    toast.dismiss(gameTimedOutId)
    navigate("/");
  };

  if (!user) {
    return null
  }

  const onGameOver = () => {
    setUser({...user, inGame: false})
  }

  return (
    <LobbyContext value={lobby}>
      {
        allPlayersReady(gameState) ?
        <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>
          <GameFrame gameState={gameState}/>
        </div>
        :
        <LobbyComponent players={lobby.players} matchConfiguration={matchConfiguration} actions={{toggleReady: requestToggleReady, onPlayerClickStartGame: requestClickStartGame}} />
      }
      <button onClick={redirectToHome}>Leave Game</button>
      <GameOverDialog winnerPlayerName={winnerPlayerName}/>
      <Toaster />
    </LobbyContext>
  )
}

export default App
