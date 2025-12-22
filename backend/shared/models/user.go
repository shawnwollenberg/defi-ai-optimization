package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a platform user
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	WalletAddress string `gorm:"uniqueIndex;not null" json:"wallet_address"`
	Email         string `gorm:"index" json:"email,omitempty"`
	
	// Subscription
	SubscriptionTier string    `gorm:"default:free" json:"subscription_tier"` // free, basic, premium
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at,omitempty"`
	
	// Preferences
	Preferences map[string]interface{} `gorm:"type:jsonb" json:"preferences,omitempty"`
	
	// Relationships
	Portfolios      []Portfolio      `gorm:"foreignKey:UserID" json:"portfolios,omitempty"`
	AutomationRules []AutomationRule `gorm:"foreignKey:UserID" json:"automation_rules,omitempty"`
	Transactions    []Transaction    `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}

// Portfolio represents a user's DeFi portfolio
type Portfolio struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"index;not null" json:"user_id"`

	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description,omitempty"`
	
	// Position data
	TotalValueUSD    float64 `gorm:"default:0" json:"total_value_usd"`
	TotalCollateral  float64 `gorm:"default:0" json:"total_collateral"`
	TotalDebt        float64 `gorm:"default:0" json:"total_debt"`
	HealthFactor     float64 `gorm:"default:0" json:"health_factor"`
	
	// Relationships
	Positions []Position `gorm:"foreignKey:PortfolioID" json:"positions,omitempty"`
	Snapshots []PortfolioSnapshot `gorm:"foreignKey:PortfolioID" json:"snapshots,omitempty"`
}

// Position represents a DeFi position (lending, borrowing, staking)
type Position struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PortfolioID uint   `gorm:"index;not null" json:"portfolio_id"`
	
	Protocol    string `gorm:"not null" json:"protocol"` // aave, compound, eigenlayer
	Chain       string `gorm:"not null" json:"chain"`    // ethereum, base
	Asset       string `gorm:"not null" json:"asset"`    // USDC, ETH, etc.
	PositionType string `gorm:"not null" json:"position_type"` // lending, borrowing, staking
	
	// Position details
	Amount      float64 `gorm:"default:0" json:"amount"`
	APY         float64 `gorm:"default:0" json:"apy"`
	Address     string  `gorm:"not null" json:"address"` // contract address
	
	// Risk metrics
	LiquidationRisk float64 `gorm:"default:0" json:"liquidation_risk"`
	LastRiskCheck   *time.Time `json:"last_risk_check,omitempty"`
}

// PortfolioSnapshot represents a historical snapshot of a portfolio
type PortfolioSnapshot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	PortfolioID uint `gorm:"index;not null" json:"portfolio_id"`
	
	TotalValueUSD   float64 `json:"total_value_usd"`
	TotalCollateral float64 `json:"total_collateral"`
	TotalDebt       float64 `json:"total_debt"`
	HealthFactor    float64 `json:"health_factor"`
	
	// Snapshot data as JSON
	SnapshotData map[string]interface{} `gorm:"type:jsonb" json:"snapshot_data"`
}

// AutomationRule represents an automation rule for portfolio management
type AutomationRule struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"index;not null" json:"user_id"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	
	// Trigger configuration
	TriggerType string                 `gorm:"not null" json:"trigger_type"` // apy_drop, health_factor, risk_threshold
	TriggerConfig map[string]interface{} `gorm:"type:jsonb" json:"trigger_config"`
	
	// Action configuration
	ActionType string                 `gorm:"not null" json:"action_type"` // rebalance, withdraw, deposit
	ActionConfig map[string]interface{} `gorm:"type:jsonb" json:"action_config"`
	
	// Execution tracking
	LastExecutedAt *time.Time `json:"last_executed_at,omitempty"`
	ExecutionCount int        `gorm:"default:0" json:"execution_count"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"index;not null" json:"user_id"`

	// Transaction details
	TxHash      string `gorm:"uniqueIndex;not null" json:"tx_hash"`
	Chain       string `gorm:"not null" json:"chain"`
	FromAddress string `gorm:"not null" json:"from_address"`
	ToAddress   string `gorm:"not null" json:"to_address"`
	
	// Transaction metadata
	Type        string  `gorm:"not null" json:"type"` // rebalance, deposit, withdraw
	Status      string  `gorm:"default:pending" json:"status"` // pending, confirmed, failed
	Value       float64 `gorm:"default:0" json:"value"`
	GasUsed     uint64  `json:"gas_used,omitempty"`
	GasPrice    string  `json:"gas_price,omitempty"`
	
	// Related automation rule
	AutomationRuleID *uint `gorm:"index" json:"automation_rule_id,omitempty"`
	
	// Transaction data
	TxData map[string]interface{} `gorm:"type:jsonb" json:"tx_data,omitempty"`
}

// Subscription represents subscription and payment tracking
type Subscription struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	Tier           string     `gorm:"not null" json:"tier"` // free, basic, premium
	StripeCustomerID string   `gorm:"index" json:"stripe_customer_id,omitempty"`
	StripeSubscriptionID string `gorm:"index" json:"stripe_subscription_id,omitempty"`
	
	Status         string     `gorm:"default:active" json:"status"` // active, cancelled, expired
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
	
	// Performance tracking
	TotalSavedLosses float64 `gorm:"default:0" json:"total_saved_losses"`
	PerformanceFee   float64 `gorm:"default:0" json:"performance_fee"`
}

