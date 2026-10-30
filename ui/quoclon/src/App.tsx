import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { getDefaultGameState, type GameState } from './game/GameState';
import { gameTimedOutId, requestClickStartGame, requestToggleReady, startConnection  } from './server/server';
import toast, { Toaster } from 'react-hot-toast';
import { useNavigate, useParams } from 'react-router';
import { gracefullyCloseConnection } from './server/server-conn';
import { DefaultLobby, isGameOver, type Lobby } from './game/lobby/lobby';
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

  const {user, setUser} = useAuth()
  
  useEffect(() => {
    if (id !== undefined) {
      if (user) {
        startConnection(id, user.ppid, user.name, setGameState, setUser, setLobby, setMatchConfiguration, redirectToHome);
      }
    }

    return () => {
      gracefullyCloseConnection("going away");
    }
  }, [user?.ppid])

  
  useEffect(() => {
      // Is there a winner (or a draw)?
      if (isGameOver(lobby)) {
          if (lobby.winnerPlayerId !== undefined) {
            const winnerPlayerName = gameState.players.get(lobby.winnerPlayerId)?.name || "Unknown";
            setWinnerPlayerName(winnerPlayerName);
          }
          onGameOver();
      }
      return () => {
        // set again to default, covering the case of leaving the game even if it wasn't game over
        if (user) {
          setUser({...user, inGame: false})
        }
      }
  }, [lobby.winnerPlayerId, lobby.draw])

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
        lobby.gameStarted ?
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
