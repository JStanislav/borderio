import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import { Home } from './components/home/home.tsx'
import './style/color.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <div className="body-container">
      
      <BrowserRouter>
        <Routes>
            <Route index path="/" element={<Home />} />
            <Route path="game/:id" element={<App />} />
            <Route path="about" element={<div>About!!!</div>} />
        </Routes>
      </BrowserRouter>

    </div>
  </StrictMode>,
)
