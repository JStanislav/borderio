import { useState } from 'react';
import './App.css'
import { GameFrame } from './Gameframe';
import { connect } from './server-conn';

export interface GameState {
  type: string,
  playerOne: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
  },
  playerTwo: {
    id: number,
    name: string,
    position: {
      row: number,
      col: number
    }
  }
  walls: Array<{
    cellA: {
      row: number,
      column: number
    },
    cellB: {
      row: number,
      column: number
    },
  }>
}

function App() {
  const [gameState, setGameState] = useState<GameState>({
  type: "gameState",
  playerOne: {
    id: 1,
    name: "P1",
    position: {
      row: 0,
      col: 0
    }
  },
  playerTwo: {
    id: 2,
    name: "P2",
    position: {
      row: 8,
      col: 8
    }
  },
  walls: []
  })


  const onMessage = (ev: MessageEvent) => {
    console.log("message arrived");
    console.log("ev", ev.data)
    setGameState(JSON.parse(ev.data));
  }


  const startConnection = () => {
    // starts socket connection
    connect(onMessage);
  }

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <GameFrame gameState={gameState}/>
      <button onClick={startConnection}>Start connection</button>
    </div>
  )
}

export default App
