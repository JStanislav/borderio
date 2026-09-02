import { Outlet } from "react-router"
import "./layout.css"
import banner from "../../assets/header-banner.png"
import { Footer } from "./Footer"
import { AuthProvider } from "../../contexts/auth-provider"
import { Navbar } from "./Navbar"
import { GameAwareLink } from "../GameAwareLink"

export const Layout = () => {
    return <AuthProvider>
        <div className="layout-container">
            <header className="header">
                <GameAwareLink to="/"><img src={banner} alt="Banner" /></GameAwareLink>
                <Navbar />
            </header>
            <main className="page-content">
                <Outlet />        
            </main>  
            <Footer />
        </div>
    </AuthProvider>
}