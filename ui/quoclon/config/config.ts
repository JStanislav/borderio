interface Config {
    serverUrl: string;
    port: number;
}

const LoadOrDefaultValue = (key: string, defaultValue: string): string => {
    const envProcess = (globalThis as any).process
    const value = (envProcess && envProcess.env && envProcess.env[key] !== undefined)
        ? envProcess.env[key]
        : defaultValue;
    
        console.log(`Config: ${key} = ${value}`);
        return value 
};

export const config: Config = {
    serverUrl: LoadOrDefaultValue("VITE_SERVER_URL", "localhost:8080"),
    port: parseInt(LoadOrDefaultValue("VITE_APP_PORT", "5173"))
}