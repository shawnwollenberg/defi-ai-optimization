import { useEffect, useRef, useState } from 'react'

interface WebSocketMessage {
  type: string
  payload: any
}

export function useWebSocket(url: string) {
  const [messages, setMessages] = useState<WebSocketMessage[]>([])
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    if (!token) return

    const wsUrl = url.replace('http', 'ws')
    const ws = new WebSocket(`${wsUrl}?token=${token}`)

    ws.onopen = () => {
      setConnected(true)
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        setMessages((prev) => [...prev, message])
      } catch (error) {
        console.error('Error parsing WebSocket message:', error)
      }
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.onclose = () => {
      setConnected(false)
      console.log('WebSocket disconnected')
    }

    wsRef.current = ws

    return () => {
      ws.close()
    }
  }, [url])

  const sendMessage = (message: WebSocketMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    }
  }

  return { messages, connected, sendMessage }
}

