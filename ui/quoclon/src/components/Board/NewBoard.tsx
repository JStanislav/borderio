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

export const Board = ({players}: {players: Players}) => {
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
                        >
                            {players.map(player => 
                                (player.position.row === indexRow && player.position.col === colIdx) ? (
                                    <img src={player.id === 1 ? playerOne : playerTwo} alt={player.name} />
                                ) : null
                            )}
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}