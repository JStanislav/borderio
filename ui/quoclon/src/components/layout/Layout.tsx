import { Outlet } from "react-router"
import "./layout.css"
import banner from "../../assets/header-banner.png"
import { Footer } from "./Footer"
import { AuthProvider } from "../../contexts/auth-provider"

export const Layout = () => {
    return <AuthProvider>
        <div className="layout-container">
            <header className="header">
                <a href="/"><img src={banner} alt="Banner" /></a>
            </header>
            <main className="page-content">
                <Outlet />        
            </main>  
            <Footer />
        </div>
    </AuthProvider>
}