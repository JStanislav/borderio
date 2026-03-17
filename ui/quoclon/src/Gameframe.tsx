import { Board } from "./components/Board/NewBoard"
import { WallPicker } from "./components/Board/WallPicker"
import "./gameframe.css"

export const GameFrame = () => {

    const players = [
        {
            id: 1,
            name: "P1",
            position: {
                row: 3,
                col: 3
            }
        },
        {
            id: 2,
            name: "P2",
            position: {
                row: 5,
                col: 5
            }
        }
    ]

    return (
        <div className="game-frame">
            <WallPicker walls={9} position="top"/>
            <Board players={players}/>
            <WallPicker walls={9} position="bottom"/>
        </div>
    )
}