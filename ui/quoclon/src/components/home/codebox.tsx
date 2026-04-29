import "./codebox.css"

interface CodeBoxProps {
    title?: string
    children?: React.ReactNode
}

export const CodeBox = (props: CodeBoxProps) => {
    return <div className="codebox-container">
        <h4 className="codebox-title">{props.title}</h4>
        <div className="codebox-box">
            {props.children}
        </div>
    </div>
}