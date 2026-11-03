import "./gametitle.css"

interface GameTitleProps {
    players: string[],
    currentUser: string,
}

export const GameTitle = ({ players, currentUser } : GameTitleProps) => {

    const span = <span className="current-user-name-title">{currentUser}</span>
    const title = players.map((p, i) => {
        const text = p === currentUser ? span : p
        if (i === players.length - 1) {
            return text
        }
        return <span>{text} vs </span>
    })

    return <h1>
        {title}
    </h1>
}