interface Config {
    httpProtocol: string;
    wsProtocol: string; 

    serverUrl: string;
    serverWSUrl: string;

    port: number;
}

const LoadOrDefaultValue = (key: string, defaultValue: string): string => {
    const envProcess = import.meta.env

    const value = (envProcess && envProcess[key] !== undefined)
        ? envProcess[key]
        : defaultValue;
    
        console.log(`Config: ${key} = ${value}`);
        return value 
};

const env = LoadOrDefaultValue("VITE_NODE_ENV", "development");
    
export const config: Config = {    
    httpProtocol: env === "production" ? "https://" : "http://",
    wsProtocol: env === "production" ? "wss://" : "ws://",

    serverUrl: LoadOrDefaultValue("VITE_SERVER_URL", "localhost:8080"),
    serverWSUrl: LoadOrDefaultValue("VITE_SERVER_WS_URL", "localhost:8080"),

    port: parseInt(LoadOrDefaultValue("VITE_APP_PORT", "5173"))
}