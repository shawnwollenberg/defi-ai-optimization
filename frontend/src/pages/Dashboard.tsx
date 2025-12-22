import { useQuery } from '@tanstack/react-query'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import './Dashboard.css'

export default function Dashboard() {
  const { data: portfolio } = useQuery({
    queryKey: ['portfolio'],
    queryFn: async () => {
      const token = localStorage.getItem('auth_token')
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/portfolios`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!response.ok) throw new Error('Failed to fetch portfolio')
      const data = await response.json()
      return data[0] || null
    },
  })

  const { data: risk } = useQuery({
    queryKey: ['risk'],
    queryFn: async () => {
      const token = localStorage.getItem('auth_token')
      const walletAddress = localStorage.getItem('wallet_address')
      if (!walletAddress) return null

      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/risk/forecast`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          user_address: walletAddress,
          positions: [],
        }),
      })
      if (!response.ok) return null
      return response.json()
    },
    enabled: !!portfolio,
  })

  // Mock data for chart
  const chartData = [
    { date: '2024-01-01', value: 10000 },
    { date: '2024-01-02', value: 10200 },
    { date: '2024-01-03', value: 10100 },
    { date: '2024-01-04', value: 10300 },
    { date: '2024-01-05', value: 10500 },
  ]

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>
      
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Value</h3>
          <p className="stat-value">${portfolio?.total_value_usd?.toLocaleString() || '0.00'}</p>
        </div>
        <div className="stat-card">
          <h3>Health Factor</h3>
          <p className="stat-value">{portfolio?.health_factor?.toFixed(2) || 'N/A'}</p>
        </div>
        <div className="stat-card">
          <h3>Liquidation Risk</h3>
          <p className={`stat-value ${risk?.liquidation_risk > 0.5 ? 'risk-high' : 'risk-low'}`}>
            {(risk?.liquidation_risk * 100).toFixed(1) || '0.0'}%
          </p>
        </div>
        <div className="stat-card">
          <h3>Active Rules</h3>
          <p className="stat-value">0</p>
        </div>
      </div>

      <div className="chart-section">
        <h2>Portfolio Value Trend</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#3a4166" />
            <XAxis dataKey="date" stroke="#a0aec0" />
            <YAxis stroke="#a0aec0" />
            <Tooltip contentStyle={{ backgroundColor: '#1a1f3a', border: '1px solid #3a4166' }} />
            <Line type="monotone" dataKey="value" stroke="#667eea" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {risk && (
        <div className="risk-section">
          <h2>Risk Assessment</h2>
          <div className="risk-card">
            <p className="risk-level">Risk Level: <span className={`risk-${risk.risk_level}`}>{risk.risk_level}</span></p>
            <p className="risk-confidence">Confidence: {(risk.confidence * 100).toFixed(0)}%</p>
            <ul className="recommendations">
              {risk.recommendations?.map((rec: string, i: number) => (
                <li key={i}>{rec}</li>
              ))}
            </ul>
          </div>
        </div>
      )}
    </div>
  )
}

