import { board } from "../../board.type"
import "./board.css"
import playerTwo from  "../../assets/players/player_one.png"
import playerOne from "../../assets/players/player_two.png"

const classes = {
    Blank: "blank",
    Square: "square",
    Wall: "wall"
}

type Players = Player[]

type Player = {
    id: number,
    name: string,
    position: {
        row: number,
        col: number
    }
}


const onDragOver = (ev: React.DragEvent<HTMLDivElement>) => {
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "move";
}

const onDragStart = (ev: React.DragEvent<HTMLImageElement>, playerId: number, row: number, col: number) => {
    console.log("drag start")
    ev.dataTransfer.setData("row", row.toString())
    ev.dataTransfer.setData("col", col.toString());
    ev.dataTransfer.setData("playerId", playerId.toString());
}

const isDraggable = (row: number, col: number) => {
    return (row % 2 !== 0) && (col % 2 !== 0);
}

export const Board = ({players, requestPlayerMove}: {players: Players, requestPlayerMove: (playerId: number, row: number, col: number) => void}) => {

    const onDrop = (ev: React.DragEvent<HTMLDivElement>, rowEnd: number, colEnd: number) => {
        ev.preventDefault();
        const row = ev.dataTransfer.getData("row");
        const col = ev.dataTransfer.getData("col");
        const playerId = ev.dataTransfer.getData("playerId");
        console.log("data start", { row: parseInt(row), col: parseInt(col) });
        console.log("data end", { rowEnd, colEnd });
        console.log("dropped");

        requestPlayerMove(parseInt(playerId), rowEnd, colEnd);
    }
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
                                ${cell.fillType ? "filled " + cell.fillType : ""}
                            `}
                            
                            onDrop={e => onDrop(e, indexRow, colIdx)}
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
                                        onDragStart={e => onDragStart(e, player.id, indexRow, colIdx)}
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