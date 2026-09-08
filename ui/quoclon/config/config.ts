interface Config {
    serverUrl: string;
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

export const config: Config = {
    serverUrl: LoadOrDefaultValue("VITE_SERVER_URL", "localhost:8080"),
    port: parseInt(LoadOrDefaultValue("VITE_APP_PORT", "5173"))
}