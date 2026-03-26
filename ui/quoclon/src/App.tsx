import { useState } from 'react';
import { GameFrame } from './components/game/Gameframe';
import { getDefaultGameState, type GameState } from './game/GameState';
import { startConnection  } from './server/server';
import { Toaster } from 'react-hot-toast';

function App() {
  const [gameState, setGameState] = useState<GameState>(getDefaultGameState())

const connectToServer = () => {
  startConnection(setGameState);
}

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <GameFrame gameState={gameState}/>
      <button onClick={connectToServer}>Start connection</button>
      <Toaster />
    </div>
  )
}

export default App
