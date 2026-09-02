import { useEffect, useRef } from "react";
import "./game/gameoverdialog.css"

interface ExitGameDialogProps {
    onConfirm: () => void;
    onClose: () => void;
    open: boolean;
}

export const ExitGameDialog = ({ onConfirm, open, onClose }: ExitGameDialogProps) => {
    const dialogRef = useRef<HTMLDialogElement>(null);
    
    useEffect(() => {
        if (open) {
            openModal();
        } else {
            closeModal();
        }
    }, [open])

    // Catch the close event, even with the backdrop click, and call prop close function
    useEffect(() => {
        const dialog = dialogRef.current;
        if (!dialog) return;
        dialog.addEventListener("close", onClose)
    }, [])

    
    const openModal = () => {
        dialogRef.current?.showModal();
    }

    const closeModal = () => {
        dialogRef.current?.close();
    }

    
    return (
        <div className="dialog-container">
            <dialog ref={dialogRef} className="game-over-dialog" closedby="any">
                <p>You are about to leave the game if you continue. Are you sure?</p>
                <div className="game-over-dialog-actions">
                    <button onClick={onClose}>Stay</button>
                    <button onClick={onConfirm}>Leave game</button>
                </div>
            </dialog>        
        </div>
    )
}