import { board } from "../../board.type"
import "./board2.css"

const classes = {
    Blank: "blank",
    Square: "square",
    Wall: "wall"
}

export const Board2 = () => {
    return (
        <div className="board">
            {board.map((row, indexRow) => 
                <div key={`row-${indexRow}`} className={row.type === "Full" ? "row-full-container" : "row-wall-container"}>
                    {row.cells.map((cell, indexCol) => 
                        <div id={cell.id} key={cell.id} className={classes[cell.type]}>
                            {cell.type}
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}