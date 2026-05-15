import { useState } from "react"
import { CodeBox } from "./codebox"
import "./home.css"
import { generateGameCode } from "./home";
import { Link, useNavigate } from "react-router";
import { GameExist } from "../../server/server";
import toast, { Toaster } from "react-hot-toast";


export const Home = () => {
    const [gameCode, setGameCode] = useState<string | null>(null);
    const [joinCode, setJoinCode] = useState<string>("");

    const navigate = useNavigate();
    
    async function PingGame() {
        const exist = await GameExist(joinCode);
        if (exist) {
            navigate(`/game/${joinCode}?action=join`);
        } else {
            setJoinCode("");
            toast.error("Game does not exist");
        }
    }

    return <div className="home-container">
        <div className="codes-container">
            <CodeBox title="Create game">
                <div>Generate a code and invite your friend to play!</div>
                {!gameCode ? 
                    <button onClick={() => setGameCode(generateGameCode())}>Generate code</button>
                :
                    <Link to={`/game/${gameCode}?action=create`}>
                        <button disabled={!gameCode}>Join</button>
                    </Link>
                }

                {gameCode && <span className="code">{gameCode}</span>}
            </CodeBox>

            <CodeBox title="Join game">
                <div>Enter a code to join a game</div>
                <input placeholder="Enter code here" value={joinCode} onChange={(e) => setJoinCode(e.target.value)} />
                <button onClick={PingGame} disabled={!joinCode}>Join game</button>
            </CodeBox>
        </div>
        <Toaster />

    </div>
}