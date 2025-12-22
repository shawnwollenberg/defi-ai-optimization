import { useQuery } from '@tanstack/react-query'
import './Settings.css'

export default function Settings() {
  const { data: subscription } = useQuery({
    queryKey: ['subscription'],
    queryFn: async () => {
      const token = localStorage.getItem('auth_token')
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/subscription`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!response.ok) return { tier: 'free', status: 'active' }
      return response.json()
    },
  })

  return (
    <div className="settings">
      <h1>Settings</h1>

      <div className="settings-section">
        <h2>Subscription</h2>
        <div className="subscription-card">
          <div className="subscription-info">
            <h3>Current Plan: {subscription?.tier || 'Free'}</h3>
            <p>Status: {subscription?.status || 'active'}</p>
          </div>
          <button className="btn-primary">Upgrade Plan</button>
        </div>

        <div className="plans-grid">
          <div className="plan-card">
            <h3>Basic</h3>
            <p className="plan-price">$10/month</p>
            <ul className="plan-features">
              <li>Basic automation</li>
              <li>Risk monitoring</li>
              <li>Email alerts</li>
            </ul>
            <button className="btn-secondary">Select</button>
          </div>
          <div className="plan-card featured">
            <h3>Premium</h3>
            <p className="plan-price">$50/month</p>
            <ul className="plan-features">
              <li>Advanced automation</li>
              <li>AI risk forecasting</li>
              <li>Priority support</li>
              <li>Performance analytics</li>
            </ul>
            <button className="btn-primary">Select</button>
          </div>
        </div>
      </div>

      <div className="settings-section">
        <h2>Preferences</h2>
        <div className="preferences-form">
          <label>
            <span>Email Notifications</span>
            <input type="checkbox" defaultChecked />
          </label>
          <label>
            <span>Risk Alert Threshold</span>
            <input type="number" defaultValue={0.5} step={0.1} min={0} max={1} />
          </label>
        </div>
      </div>
    </div>
  )
}

