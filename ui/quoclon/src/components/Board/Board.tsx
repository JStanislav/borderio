import { board, classes } from "./board.type"
import "./board.css"
import playerTwo from  "../../assets/players/player_one.png"
import playerOne from "../../assets/players/player_two.png"
import type { Players } from "../game/player.type"
import { isDraggable, isFinishLine, onClick, onDragOver, onDragStart, onDrop, setActiveWalls } from "./board.utils"

interface Props {
    players: Players,
    requestPlayerMove: (playerId: number, row: number, col: number) => void,
    requestWallPlacement: (playerId: number, row: number, col: number, orientation: "horizontal" | "vertical") => void,
    activeWalls: { row: number, col: number, orientation: "horizontal" | "vertical" }[],
    currentTurnPlayerId: number
}

export const Board = ({players, requestPlayerMove, requestWallPlacement, activeWalls, currentTurnPlayerId}: Props) => {

    setActiveWalls(board, activeWalls);

    return (
        <div className="board">
            {board.map((row, indexRow) => 
                <div key={`row-${indexRow}`} className={row.type === "Full" ? "row-full-container" : "row-wall-container"}>
                    {row.cells.map((cell, colIdx) => 
                        <div 
                            id={cell.id} key={cell.id}
                            className={`
                                ${classes[cell.type]} \
                                ${(colIdx % 2) === 0 ? "narrow-col" : "wide-col"} \
                                ${cell.fillType ? "filled " + cell.fillType : ""} \ 
                                ${isFinishLine(indexRow) ? "finish-line" : ""}
                            `}

                            onClick={e => onClick(e, requestWallPlacement, currentTurnPlayerId)}
                            
                            onDrop={e => onDrop(e, indexRow, colIdx, requestPlayerMove)}
                            onDragOver={isDraggable(indexRow, colIdx) ? onDragOver : undefined}
                        >
                            {players.map(player => 
                                (player.position.row === indexRow && player.position.col === colIdx) ? (
                                    <img
                                        key={`player-${player.id}`}
                                        className="img-player"
                                        src={player.id === 1 ? playerOne : playerTwo} 
                                        alt={player.name}
                                        draggable
                                        onDragStart={e => onDragStart(e, player.id)}
                                    />
                                ) : null
                            )}
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}