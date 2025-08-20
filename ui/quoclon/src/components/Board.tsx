import { getBoard } from "../types"
import "./board.css"



export const Board = () => {
    return (
        <div className="board-grid">
            {getBoard().map((row, indexRow) => {
                return (
                    <div id={`row-${indexRow}`} className="grid-row">
                        {row.map((cell, indexCol) => {
                            return (
                                <div id={`cell-${indexRow}-${indexCol}`} className="cell">
                                    <div className={`grid-box ${cell.bottomWalled ? "box-bottom-walled" : ""} ${cell.rightWalled ? "right-walled" : ""}`} key={`${indexRow}-${indexCol}`}>
                                        {cell.type}
                                    </div>
                                    <div className={`left-border ${cell.leftWalled ? "visible" : ""}`}></div>
                                    <div className={`top-border ${cell.topWalled ? "visible" : ""}`}></div>
                                    <div className={`right-border ${cell.rightWalled ? "visible" : ""}`}></div>
                                    <div className={`bottom-border ${cell.bottomWalled ? "visible" : ""}`}></div>
                               </div>
                                
                            )
                        })}
                    </div>
                )
            })}
        </div>
    )
}