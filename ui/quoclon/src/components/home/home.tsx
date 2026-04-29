import { useState } from "react"
import { CodeBox } from "./codebox"
import "./home.css"
import { generateGameCode } from "./home";
import { Link } from "react-router";


export const Home = () => {
    const [gameCode, setGameCode] = useState<string | null>(null);

    return <div className="home-container">
        <div className="codes-container">
            <CodeBox title="Create game">
                <div>Generate a code and invite your friend to play!</div>
                {!gameCode ? 
                    <button onClick={() => setGameCode(generateGameCode())}>Generate code</button>
                :
                    <Link to={`/game/${gameCode}`}>
                        <button>Join</button>
                    </Link>
                }

                {gameCode && <span className="code">{gameCode}</span>}
            </CodeBox>

            <CodeBox title="Join game">
                <div>Enter a code to join a game</div>
                <input placeholder="Enter code here" />
                <button>Join game</button>
            </CodeBox>
        </div>
    </div>
}