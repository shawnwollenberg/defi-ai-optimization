import { useState, useEffect } from 'react'
import { ethers } from 'ethers'
import './WalletConnect.css'

interface WalletState {
  address: string | null
  connected: boolean
}

export default function WalletConnect() {
  const [wallet, setWallet] = useState<WalletState>({
    address: null,
    connected: false,
  })
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    // Check if wallet is already connected
    const savedAddress = localStorage.getItem('wallet_address')
    if (savedAddress) {
      setWallet({ address: savedAddress, connected: true })
    }
  }, [])

  const connectWallet = async () => {
    if (typeof window.ethereum === 'undefined') {
      alert('Please install MetaMask or another Web3 wallet')
      return
    }

    setLoading(true)
    try {
      const provider = new ethers.BrowserProvider(window.ethereum)
      const accounts = await provider.send('eth_requestAccounts', [])
      const address = accounts[0]

      setWallet({ address, connected: true })
      localStorage.setItem('wallet_address', address)

      // Authenticate with backend
      await authenticateWithBackend(address)
    } catch (error) {
      console.error('Error connecting wallet:', error)
      alert('Failed to connect wallet')
    } finally {
      setLoading(false)
    }
  }

  const disconnectWallet = () => {
    setWallet({ address: null, connected: false })
    localStorage.removeItem('wallet_address')
    localStorage.removeItem('auth_token')
  }

  const authenticateWithBackend = async (address: string) => {
    try {
      const message = `Sign in to DeFi Optimizer\n\nAddress: ${address}\nTimestamp: ${Date.now()}`
      
      const provider = new ethers.BrowserProvider(window.ethereum)
      const signer = await provider.getSigner()
      const signature = await signer.signMessage(message)

      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/auth/wallet`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          wallet_address: address,
          signature,
          message,
        }),
      })

      if (response.ok) {
        const data = await response.json()
        localStorage.setItem('auth_token', data.token)
      }
    } catch (error) {
      console.error('Error authenticating:', error)
    }
  }

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`
  }

  if (wallet.connected && wallet.address) {
    return (
      <div className="wallet-connected">
        <span className="wallet-address">{formatAddress(wallet.address)}</span>
        <button onClick={disconnectWallet} className="disconnect-btn">
          Disconnect
        </button>
      </div>
    )
  }

  return (
    <button
      onClick={connectWallet}
      disabled={loading}
      className="connect-btn"
    >
      {loading ? 'Connecting...' : 'Connect Wallet'}
    </button>
  )
}

