import { Board } from "./components/Board/NewBoard"
import { WallPicker } from "./components/Board/WallPicker"
import "./gameframe.css"

export const GameFrame = () => {
    return (
        <div className="game-frame">
            <WallPicker walls={9} position="top"/>
            <Board />
            <WallPicker walls={9} position="bottom"/>
        </div>
    )
}