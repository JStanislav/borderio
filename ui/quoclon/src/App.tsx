import './App.css'
import { Board } from './components/Board'
import { Board2 } from './components/Board2/NewBoard'

function App() {

  return (
    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>      
      <h1>Quoridor</h1>
      <Board2 />
    </div>
  )
}

export default App
