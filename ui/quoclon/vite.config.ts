import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export const config = {
  serverUrl: process.env.VITE_SERVER_URL || "localhost:8080",
  port: parseInt(process.env.VITE_PORT || "5173")
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: config.port
  }
})
