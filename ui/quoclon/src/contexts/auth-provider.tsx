import { createContext, useContext, useEffect, useState } from "react";
import { type Player } from "../game/player";
import { generatePPID } from "../app";

interface AuthPlayerContext {
    user: Player | null;
    setUser: (user: Player) => void;
}

export const AuthPlayerContext = createContext<AuthPlayerContext | undefined>(undefined);

interface AuthProviderProps {
    children: React.ReactNode;
    user?: Player;
}

export const AuthProvider = ({ children, user: initialUser }: AuthProviderProps) => {
    const [user, setUser] = useState<Player | null>(initialUser ?? null);

    useEffect(() => {
        async function fetchUser() {
            await new Promise(resolve => setTimeout(resolve, 1000));
            const ppid = generatePPID();
            setUser({name: ppid, id: -1, ppid: ppid});
        }

        if (!initialUser) {
            fetchUser();
        }

    }, [initialUser]);

    return (
        <AuthPlayerContext value={{ user: user, setUser: setUser }}>
            {children}
        </AuthPlayerContext>
    )
}

export const useAuth = () => {
    const context = useContext(AuthPlayerContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within a AuthProvider');
    }
    return context;
}