import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import { Home } from './components/home/home.tsx'
import './style/color.css'
import './style/index.css'
import { Layout } from './components/layout/Layout.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route index path="/" element={<Home />} />
            <Route path="game/:id" element={<App />} />
            <Route path="about" element={<div>About!!!</div>} />    
          </Route>
        </Routes>
      </BrowserRouter>
  </StrictMode>,
)
