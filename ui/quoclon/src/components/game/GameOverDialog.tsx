import { useContext, useEffect, useRef } from "react";
import { getGameOverText } from "./player.type";
import { LobbyContext, PlayerContext } from "../../App.tsx";
import "./gameoverdialog.css"
import { useNavigate } from "react-router";

interface Props {
    winnerPlayerName: string,
}

export const GameOverDialog = ({winnerPlayerName}: Props) => {
    const dialogRef = useRef<HTMLDialogElement>(null);
    const lobbyContext = useContext(LobbyContext);
    const playerContext = useContext(PlayerContext);

    const navigate = useNavigate();

    const leaveGame = () => {
        
        navigate("/");
    }
    
    useEffect(() => {
        if (lobbyContext.winnerPlayerId !== undefined) {
            onOpenModal();
        }
    }, [lobbyContext.winnerPlayerId])

    const onOpenModal = () => {
        dialogRef.current?.showModal();
    }

    const onCloseModal = () => {
        dialogRef.current?.close();
    }

    
    return (
        <div className="dialog-container">
            <dialog ref={dialogRef} className="game-over-dialog" closedby="any">
                <p>{getGameOverText(winnerPlayerName, lobbyContext.winnerPlayerId === playerContext.id)}</p>
                <div className="game-over-dialog-actions">
                    <button onClick={leaveGame}>Go back to home</button>
                    <button onClick={onCloseModal}>Close</button>
                </div>
            </dialog>        
        </div>
    )
}