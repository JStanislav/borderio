import './App.css'
import { GameFrame } from './Gameframe';
import { connect, send } from './server-conn';

function App() {

  const startConnection = () => {
    // starts socket connection
    connect();
  }

  const sendData = () => {
    // sends data to server
    console.log("sending data")
    send("hello", "world");
  }

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <GameFrame />
      <button onClick={startConnection}>Start connection</button>
      <button onClick={sendData}>Send data</button>
    </div>
  )
}

export default App
