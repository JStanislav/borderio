import { createContext, useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { allPlayersReady, getDefaultGameState, type GameState } from './game/GameState';
import { startConnection  } from './server/server';
import { Toaster } from 'react-hot-toast';
import { useNavigate, useParams, useSearchParams } from 'react-router';
import { canDisplayStartButton, generatePPID } from './app';
import { send, closeConn } from './server/server-conn';
import { DefaultPlayer, type Player } from './game/player';
import { DefaultLobby, type Lobby } from './game/lobby/lobby';
import type { MatchConfiguration } from './game/MatchConfiguration';

const getReadyText = (ready: boolean): string => {
  return ready ? "Ready" : "Not Ready";
}

export const LobbyContext = createContext<Lobby>(DefaultLobby);
export const PlayerContext = createContext<Player>(DefaultPlayer);

function App() {
    const [gameState, setGameState] = useState<GameState>(getDefaultGameState())
    const [player, setPlayer] = useState<Player>(DefaultPlayer)
    const [lobby, setLobby] = useState<Lobby>(DefaultLobby)
    const [matchConfiguration, setMatchConfiguration] = useState<MatchConfiguration>({ playerAmount: 2 })
    const navigate = useNavigate();

    const { id } = useParams()
    const [searchParams] = useSearchParams()

  useEffect(() => {
    if (id !== undefined && searchParams.get("action") !== null) {
      const ppid = generatePPID();
      setPlayer({...player, ppid: ppid});
      
      const action = searchParams.get("action") as "create" | "join"

      startConnection(id, action, ppid, setGameState, setPlayer, setLobby, setMatchConfiguration, redirectToHome);
    }

    return () => {
      closeConn()
    }
  }, [])

  const redirectToHome = () => navigate("/");

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
    <LobbyContext value={lobby}>
      <PlayerContext value={player}>
      {
        allPlayersReady(gameState) ?
        <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>
          <GameFrame gameState={gameState}/>
        </div>
        :
        <div>
          Waiting for others players to be ready...
          {lobby.players.map((lobbyPlayer, index) => <div key={index}>{`${lobbyPlayer.name}${lobbyPlayer.host ? "[H]" : ""}  ${getReadyText(lobbyPlayer.ready)}`}</div>)}
          <button onClick={toggleReady}>{player.ready ? "Unready" : "Ready"}</button>
          {canDisplayStartButton(lobby, matchConfiguration, player) && <button onClick={onClickStartGame}>Start</button>}
        </div>
      }
      <button onClick={redirectToHome}>Leave Game</button>
        <Toaster />
      </PlayerContext>
    </LobbyContext>
  )
}

export default App
