import { useState } from 'react';
import './App.css'
import { GameFrame } from './Gameframe';
import { connect, send } from './server-conn';

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
  }
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

  const sendData = () => {
    // sends data to server
    console.log("sending data")
    send("hello", "world");
  }

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <GameFrame gameState={gameState}/>
      <button onClick={startConnection}>Start connection</button>
      <button onClick={sendData}>Send data</button>
    </div>
  )
}

export default App
