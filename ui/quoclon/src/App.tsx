import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { allPlayersReady, getDefaultGameState, getPlayerById, type GameState } from './game/GameState';
import { gameTimedOutId, startConnection  } from './server/server';
import toast, { Toaster } from 'react-hot-toast';
import { useNavigate, useParams, useSearchParams } from 'react-router';
import { send, gracefullyCloseConnection } from './server/server-conn';
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
  }, [lobby.winnerPlayerId])

  const redirectToHome = () => {
    toast.dismiss(gameTimedOutId)
    navigate("/");
  };

  if (!user) {
    return null
  }

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
        <LobbyComponent players={lobby.players} matchConfiguration={matchConfiguration} actions={{toggleReady, onPlayerClickStartGame: onClickStartGame}} />
      }
      <button onClick={redirectToHome}>Leave Game</button>
      <GameOverDialog winnerPlayerName={winnerPlayerName}/>
      <Toaster />
    </LobbyContext>
  )
}

export default App
