import { setWall, TOTAL_BOARD_DIM, type Row } from "./board.type";

export const onDragOver = (ev: React.DragEvent<HTMLDivElement>) => {
    ev.preventDefault();
    ev.dataTransfer.dropEffect = "move";
}

export const onDragStart = (ev: React.DragEvent<HTMLImageElement>, playerId: number) => {
    ev.dataTransfer.setData("playerId", playerId.toString());
}

export const isDraggable = (row: number, col: number) => {
    return (row % 2 !== 0) && (col % 2 !== 0);
}

const isVerticalCell = (cellId: number, totalBoardDim: number) => {
    return (cellId % (2*totalBoardDim) ) >= totalBoardDim 
}


export const onClick = (ev: React.MouseEvent<HTMLDivElement>, requestWallPlacement: (playerId: number, row: number, col: number, orientation: "horizontal" | "vertical") => void) => {
    const cellId = Number(ev.currentTarget.id.split("-")[1]);
    if (!cellId) return;
    if (cellId % 2 === 0) return; // only allow clicking on wall placeable cells
    const isVertical = isVerticalCell(cellId, TOTAL_BOARD_DIM);
    
    requestWallPlacement(1, 
        Math.floor(cellId / TOTAL_BOARD_DIM),
        (cellId % TOTAL_BOARD_DIM), 
        isVertical ? "vertical" : "horizontal");
}

export const onDrop = (ev: React.DragEvent<HTMLDivElement>, rowEnd: number, colEnd: number, requestPlayerMove: (playerId: number, row: number, col: number) => void) => {
    ev.preventDefault();
    const playerId = ev.dataTransfer.getData("playerId");

    requestPlayerMove(parseInt(playerId), rowEnd, colEnd);
}

export const setActiveWalls = (board: Array<Row>, activeWalls: { row: number, col: number, orientation: "horizontal" | "vertical" }[]) => {
    activeWalls.forEach(wall => {
        setWall(board, wall.row, wall.col, 2);})
}

export const isFinishLine = (row: number): boolean=> {
    return row === 1 || row === TOTAL_BOARD_DIM + 2
}