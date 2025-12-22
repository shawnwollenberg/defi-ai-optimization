import { useQuery } from '@tanstack/react-query'
import './Portfolios.css'

export default function Portfolios() {
  const { data: portfolios, isLoading } = useQuery({
    queryKey: ['portfolios'],
    queryFn: async () => {
      const token = localStorage.getItem('auth_token')
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/portfolios`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!response.ok) throw new Error('Failed to fetch portfolios')
      return response.json()
    },
  })

  if (isLoading) {
    return <div className="portfolios">Loading...</div>
  }

  return (
    <div className="portfolios">
      <div className="portfolios-header">
        <h1>Portfolios</h1>
        <button className="btn-primary">Create Portfolio</button>
      </div>

      <div className="portfolios-grid">
        {portfolios?.map((portfolio: any) => (
          <div key={portfolio.id} className="portfolio-card">
            <h2>{portfolio.name}</h2>
            <p className="portfolio-description">{portfolio.description}</p>
            <div className="portfolio-stats">
              <div className="stat">
                <span className="stat-label">Total Value</span>
                <span className="stat-value">${portfolio.total_value_usd?.toLocaleString() || '0.00'}</span>
              </div>
              <div className="stat">
                <span className="stat-label">Health Factor</span>
                <span className="stat-value">{portfolio.health_factor?.toFixed(2) || 'N/A'}</span>
              </div>
            </div>
            <div className="portfolio-actions">
              <button className="btn-secondary">View Details</button>
              <button className="btn-secondary">Manage</button>
            </div>
          </div>
        )) || (
          <div className="empty-state">
            <p>No portfolios yet. Create your first portfolio to get started!</p>
          </div>
        )}
      </div>
    </div>
  )
}

