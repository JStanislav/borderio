import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { allPlayersReady, getDefaultGameState, type GameState } from './game/GameState';
import { startConnection  } from './server/server';
import { Toaster } from 'react-hot-toast';
import { useParams, useSearchParams } from 'react-router';
import { canDisplayStartButton, generatePPID } from './app';
import { send, closeConn } from './server/server-conn';
import { DefaultPlayer, type Player } from './game/player';
import type { LobbyPlayer } from './server/messages';

const getReadyText = (ready: boolean): string => {
  return ready ? "Ready" : "Not Ready";
}

const GameFrameContext = createContext<GameState>(getDefaultGameState());
export const PlayerContext = createContext<Player>(DefaultPlayer);

function App() {
  const [gameState, setGameState] = useState<GameState>(getDefaultGameState())
    const [player, setPlayer] = useState<Player>(DefaultPlayer)
    const [lobbyPlayers, setLobbyPlayers] = useState<LobbyPlayer[]>([])
  
    const { id } = useParams()
    const [searchParams] = useSearchParams()

  useEffect(() => {
    if (id !== undefined && searchParams.get("action") !== null) {
      const ppid = generatePPID();
      setPlayer({...player, ppid: ppid});
      
      const action = searchParams.get("action") as "create" | "join"

      startConnection(id, action, ppid, setGameState, setPlayer, setLobbyPlayers);
    }

    return () => {
      closeConn()
    }
  }, [])

  const toggleReady = () => {
      setPlayer({...player, ready: !player.ready});
      const type = "playerReady";
      const data = {playerId: player.id, ppid: player.ppid, ready: !player.ready };
      send(type, data);
  }
  const onClickStartGame = () => {
    const type = "startGame";
    const data = {ppid: player.ppid};
    send(type, data);
  }


  return (
    <GameFrameContext value={gameState}>
      <PlayerContext value={player}>
      {
        allPlayersReady(gameState) ?
        <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
          <GameFrame gameState={gameState}/>
        </div>
        :
        <div>
          Waiting for others players to be ready...
          {lobbyPlayers.map((lobbyPlayer, index) => <div key={index}>{lobbyPlayer.name + getReadyText(lobbyPlayer.ready)}</div>)}
          <button onClick={toggleReady}>{player.ready ? "Unready" : "Ready"}</button>
          {canDisplayStartButton(lobbyPlayers, player) && <button onClick={onClickStartGame}>Start</button>}
        </div>
      }
        <Toaster />
      </PlayerContext>
    </GameFrameContext>
  )
}

export default App
