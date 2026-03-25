# Borderio

Borderio is a web-based implementation of a two-player abstract strategy board game inspired by [Quoridor](https://en.wikipedia.org/wiki/Quoridor). It is an online platform with automatic matchmaking, real-time communication via WebSockets, and a game engine written in Go.

> This project is not affiliated with or endorsed by Gigamic. Quoridor is a registered trademark of Gigamic.

---

## Overview

Players can search for a match from the browser. When two players are queued, the server pairs them and a game session begins. All game state is managed server-side and communicated to each client in real time. There is no persistent storage — sessions live entirely in memory.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Game engine & server | Go ≥ 1.24 |
| Frontend | React + TypeScript (Vite) |
| Styling | Plain CSS |
| Communication | WebSockets |

---

## Project Structure

```
borderio/
├── server/
│   ├── main.go
│   ├── game/         # Game rules, state definitions, and logic
│   ├── graph/        # Board representation using graph data structures
│   ├── player/       # Player definitions and utilities
│   ├── utils/        # Position helpers and shared utilities
│   └── websocket/
│       └── messages/ # DTOs and message structure definitions
└── client/
    ├── src/
    │   ├── components/
    │   ├── game/     # Game definitions and client-side utilities
    │   ├── server/   # Server communication and WebSocket logic
    │   └── App.tsx
    └── index.html
```

---

## Requirements

- [Go](https://go.dev/dl/) ≥ 1.24
- [Node.js](https://nodejs.org/) ≥ 24.2 (recommended)
- npm or any compatible package manager

---

## Installation

Clone the repository:

```bash
git clone https://github.com/JStanislav/borderio.git
cd borderio
```

Install frontend dependencies:

```bash
cd ui/quoclon
npm install
```

---

## Running the Project

### Backend

From the repository root:

```bash
cd server
go run main.go
```

The server starts on `http://localhost:8080` by default.

### Frontend

From the `ui/quoclon` directory:

```bash
npm run dev
```

The app will be available at `http://localhost:5173`.

---

## Goals

- Provide a fully playable online Quoridor experience in the browser.
- Implement a correct and efficient game engine, including wall placement validation using graph-based pathfinding.
- Support automatic matchmaking so players can find opponents without manual coordination.
- Keep the architecture simple and stateless (no database), using in-memory session management.
- Serve as a foundation for future features such as game history, rankings, or AI opponents.

---

## License

MIT