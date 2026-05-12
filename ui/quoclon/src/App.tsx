import { useEffect, useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { getDefaultGameState, type GameState } from './game/GameState';
import { startConnection  } from './server/server';
import { Toaster } from 'react-hot-toast';
import { useParams, useSearchParams } from 'react-router';
import { generatePPID } from './app';

function App() {
  const [gameState, setGameState] = useState<GameState>(getDefaultGameState())
  
    let { id } = useParams()
    let [searchParams] = useSearchParams()

  useEffect(() => {
    console.log("id: ", id, "action: ", searchParams.get("action"))

    if (id !== undefined && searchParams.get("action") !== null) {
      startConnection(id, searchParams.get("action") as "create" | "join", generatePPID(), setGameState);
    }
  }, [])

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <GameFrame gameState={gameState}/>
      <Toaster />
    </div>
  )
}

export default App
