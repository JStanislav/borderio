import './App.css'
import { Board } from './components/Board/NewBoard'
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
      <Board />
      <button onClick={startConnection}>Start connection</button>
      <button onClick={sendData}>Send data</button>
    </div>
  )
}

export default App
