import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import './Automation.css'

export default function Automation() {
  const [showCreateForm, setShowCreateForm] = useState(false)

  const { data: rules } = useQuery({
    queryKey: ['automation-rules'],
    queryFn: async () => {
      const token = localStorage.getItem('auth_token')
      const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/automation/rules`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!response.ok) throw new Error('Failed to fetch automation rules')
      return response.json()
    },
  })

  return (
    <div className="automation">
      <div className="automation-header">
        <h1>Automation Rules</h1>
        <button className="btn-primary" onClick={() => setShowCreateForm(true)}>
          Create Rule
        </button>
      </div>

      <div className="rules-list">
        {rules?.map((rule: any) => (
          <div key={rule.id} className="rule-card">
            <div className="rule-header">
              <h3>{rule.name}</h3>
              <label className="toggle">
                <input type="checkbox" checked={rule.enabled} readOnly />
                <span className="slider"></span>
              </label>
            </div>
            <p className="rule-description">{rule.description}</p>
            <div className="rule-details">
              <div className="detail">
                <span className="detail-label">Trigger:</span>
                <span className="detail-value">{rule.trigger_type}</span>
              </div>
              <div className="detail">
                <span className="detail-label">Action:</span>
                <span className="detail-value">{rule.action_type}</span>
              </div>
              <div className="detail">
                <span className="detail-label">Executions:</span>
                <span className="detail-value">{rule.execution_count}</span>
              </div>
            </div>
            <div className="rule-actions">
              <button className="btn-secondary">Edit</button>
              <button className="btn-secondary">Delete</button>
            </div>
          </div>
        )) || (
          <div className="empty-state">
            <p>No automation rules yet. Create your first rule to automate your DeFi strategy!</p>
          </div>
        )}
      </div>
    </div>
  )
}

