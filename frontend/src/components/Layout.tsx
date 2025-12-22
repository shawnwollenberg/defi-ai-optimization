import { Link, useLocation } from 'react-router-dom'
import WalletConnect from './WalletConnect'
import './Layout.css'

interface LayoutProps {
  children: React.ReactNode
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation()

  return (
    <div className="layout">
      <header className="header">
        <div className="header-content">
          <h1 className="logo">DeFi Optimizer</h1>
          <nav className="nav">
            <Link
              to="/"
              className={location.pathname === '/' ? 'active' : ''}
            >
              Dashboard
            </Link>
            <Link
              to="/portfolios"
              className={location.pathname === '/portfolios' ? 'active' : ''}
            >
              Portfolios
            </Link>
            <Link
              to="/automation"
              className={location.pathname === '/automation' ? 'active' : ''}
            >
              Automation
            </Link>
            <Link
              to="/settings"
              className={location.pathname === '/settings' ? 'active' : ''}
            >
              Settings
            </Link>
          </nav>
          <WalletConnect />
        </div>
      </header>
      <main className="main-content">{children}</main>
    </div>
  )
}

