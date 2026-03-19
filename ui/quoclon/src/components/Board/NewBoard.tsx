import { board, setWall, TOTAL_BOARD_DIM } from "../../board.type"
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

const isVerticalCell = (cellId: number, totalBoardDim: number) => {
    return (cellId % (2*totalBoardDim) ) >= totalBoardDim 
}

interface Props {
    players: Players,
    requestPlayerMove: (playerId: number, row: number, col: number) => void,
    requestWallPlacement: (playerId: number, row: number, col: number, orientation: "horizontal" | "vertical") => void,
    activeWalls: { row: number, col: number, orientation: "horizontal" | "vertical" }[],
}

export const Board = ({players, requestPlayerMove, requestWallPlacement, activeWalls}: Props) => {

        activeWalls?.forEach(wall => {
            setWall(board, wall.row, wall.col, 2);})

    const onClick = (ev: React.MouseEvent<HTMLDivElement>) => {
        const cellId = Number(ev.currentTarget.id.split("-")[1]);

        if (!cellId) return;
        if (cellId % 2 === 0) return; // only allow clicking on wall placeable cells

        const isVertical = isVerticalCell(cellId, TOTAL_BOARD_DIM);

        if(isVertical) {
            console.log("clicked vertical cell", cellId);
        } else {
            console.log("clicked horizontal cell", cellId);
        }
        
        requestWallPlacement(1, 
            Math.floor(cellId / TOTAL_BOARD_DIM),
            (cellId % TOTAL_BOARD_DIM), 
            isVertical ? "vertical" : "horizontal");
    }

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

                            onClick={onClick}
                            
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