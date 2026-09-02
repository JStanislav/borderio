import { Link, useNavigate } from "react-router";
import { useAuth } from "../contexts/auth-provider";
import { useState } from "react";
import { ExitGameDialog } from "./ExitGameDialog";

interface GameAwareLinkProps {
  to: string;
  children: React.ReactNode;
}

/* 
    GameAwareLink is a component that behaves like a Link. But if a user is in a game,
    it will render a dialog to confirm the action and warn the user that it's leaving the game.

    If the user is not in a game, it will render the Link as usual.

    Note that if redirection is triggered with navigate(), this component will not catch the redirection.
*/
export const GameAwareLink = ({to, children, ...props}: GameAwareLinkProps) => {
    const { user } = useAuth();
    const [showDialog, setShowDialog] = useState(false)
    const navigate = useNavigate();


    if (!(user?.inGame)) {
        return <Link to={to} {...props}>{children}</Link>
    }

    const onClick = (e: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
        e.preventDefault();
        setShowDialog(true);
    }

    return (<>
            <Link to={to} onClick={onClick} {...props}>
                {children}
            </Link>

            <ExitGameDialog open={showDialog} onClose={() => setShowDialog(false)} onConfirm={() => {
                setShowDialog(false);
                navigate(to);
            }} />
        </>
    )

}