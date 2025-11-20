import './App.css'
import { Board } from './components/Board/NewBoard'

function App() {

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <Board />
    </div>
  )
}

export default App
