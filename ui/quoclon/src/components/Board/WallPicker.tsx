import "./wallpicker.css"

interface WallPickerProps {
    walls: number;
    position: "top" | "bottom";
}

export const WallPicker = ({walls, position}: WallPickerProps) => {
    
    return (
        <div className={`picker-container ${position}-position`}>
            <div className="wall-picker">
                {
                    walls > 0 &&
                        <div key={`wall-0-${position}`} className="wall-to-be-picked-top"/>
                }
                {
                    walls > 1 &&
                        Array.from({length: walls}, (_, index) => 
                        <div key={`wall-${index}`} className="wall-to-be-picked"/>)
                }
            </div>
        </div>
    )
}