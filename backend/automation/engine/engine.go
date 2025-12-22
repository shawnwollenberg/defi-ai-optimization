package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
)

// Engine manages automation rules and executes actions
type Engine struct {
	defiServiceURL   string
	walletServiceURL string
	mlServiceURL     string
	httpClient       *http.Client
}

// NewEngine creates a new automation engine
func NewEngine(defiServiceURL, walletServiceURL, mlServiceURL string) *Engine {
	return &Engine{
		defiServiceURL:   defiServiceURL,
		walletServiceURL: walletServiceURL,
		mlServiceURL:     mlServiceURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Start begins monitoring and executing automation rules
func (e *Engine) Start(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Println("Automation engine started")

	// Initial check
	e.processRules(ctx)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			e.processRules(ctx)
		}
	}
}

// processRules processes all enabled automation rules
func (e *Engine) processRules(ctx context.Context) {
	var rules []models.AutomationRule
	if err := database.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		log.Printf("Error fetching automation rules: %v", err)
		return
	}

	for _, rule := range rules {
		if err := e.evaluateRule(ctx, rule); err != nil {
			log.Printf("Error evaluating rule %d: %v", rule.ID, err)
			continue
		}
	}
}

// evaluateRule evaluates a single automation rule
func (e *Engine) evaluateRule(ctx context.Context, rule models.AutomationRule) error {
	// Check if trigger conditions are met
	triggered, err := e.checkTrigger(ctx, rule)
	if err != nil {
		return fmt.Errorf("error checking trigger: %w", err)
	}

	if !triggered {
		return nil
	}

	// Execute action
	if err := e.executeAction(ctx, rule); err != nil {
		return fmt.Errorf("error executing action: %w", err)
	}

	// Update rule execution tracking
	now := time.Now()
	rule.LastExecutedAt = &now
	rule.ExecutionCount++
	if err := database.DB.Save(&rule).Error; err != nil {
		log.Printf("Error updating rule execution: %v", err)
	}

	return nil
}

// checkTrigger checks if a rule's trigger conditions are met
func (e *Engine) checkTrigger(ctx context.Context, rule models.AutomationRule) (bool, error) {
	switch rule.TriggerType {
	case "apy_drop":
		return e.checkAPYDrop(ctx, rule)
	case "health_factor":
		return e.checkHealthFactor(ctx, rule)
	case "risk_threshold":
		return e.checkRiskThreshold(ctx, rule)
	default:
		return false, fmt.Errorf("unknown trigger type: %s", rule.TriggerType)
	}
}

// checkAPYDrop checks if APY has dropped below threshold
func (e *Engine) checkAPYDrop(ctx context.Context, rule models.AutomationRule) (bool, error) {
	config := rule.TriggerConfig
	if config == nil {
		return false, fmt.Errorf("trigger config is nil")
	}

	protocol, ok := config["protocol"].(string)
	if !ok {
		return false, fmt.Errorf("protocol not specified in trigger config")
	}

	asset, ok := config["asset"].(string)
	if !ok {
		return false, fmt.Errorf("asset not specified in trigger config")
	}

	threshold, ok := config["threshold"].(float64)
	if !ok {
		return false, fmt.Errorf("threshold not specified in trigger config")
	}

	chain := "ethereum"
	if c, ok := config["chain"].(string); ok {
		chain = c
	}

	// Fetch current APY from DeFi service
	url := fmt.Sprintf("%s/api/v1/protocols/%s/apy?asset=%s&chain=%s", e.defiServiceURL, protocol, asset, chain)
	resp, err := e.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch APY: status %d", resp.StatusCode)
	}

	var result struct {
		APY float64 `json:"apy"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Check if APY is below threshold
	return result.APY < threshold, nil
}

// checkHealthFactor checks if health factor is below threshold
func (e *Engine) checkHealthFactor(ctx context.Context, rule models.AutomationRule) (bool, error) {
	config := rule.TriggerConfig
	if config == nil {
		return false, fmt.Errorf("trigger config is nil")
	}

	threshold, ok := config["threshold"].(float64)
	if !ok {
		return false, fmt.Errorf("threshold not specified in trigger config")
	}

	protocol, ok := config["protocol"].(string)
	if !ok {
		protocol = "aave" // Default
	}

	// Get user from rule
	var user models.User
	if err := database.DB.First(&user, rule.UserID).Error; err != nil {
		return false, err
	}

	chain := "ethereum"
	if c, ok := config["chain"].(string); ok {
		chain = c
	}

	// Fetch health factor from DeFi service
	url := fmt.Sprintf("%s/api/v1/protocols/%s/health-factor?user_address=%s&chain=%s",
		e.defiServiceURL, protocol, user.WalletAddress, chain)
	resp, err := e.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch health factor: status %d", resp.StatusCode)
	}

	var result struct {
		HealthFactor float64 `json:"health_factor"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Check if health factor is below threshold
	return result.HealthFactor < threshold, nil
}

// checkRiskThreshold checks if risk exceeds threshold
func (e *Engine) checkRiskThreshold(ctx context.Context, rule models.AutomationRule) (bool, error) {
	config := rule.TriggerConfig
	if config == nil {
		return false, fmt.Errorf("trigger config is nil")
	}

	threshold, ok := config["threshold"].(float64)
	if !ok {
		return false, fmt.Errorf("threshold not specified in trigger config")
	}

	// Get user from rule
	var user models.User
	if err := database.DB.First(&user, rule.UserID).Error; err != nil {
		return false, err
	}

	// Fetch risk forecast from ML service
	// This is a simplified version - in production, you'd fetch actual positions
	url := fmt.Sprintf("%s/api/v1/risk/forecast", e.mlServiceURL)
	reqBody := map[string]interface{}{
		"user_address": user.WalletAddress,
		"positions":    []interface{}{},
	}

	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/json")
	// Note: In production, you'd properly set the request body

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch risk forecast: status %d", resp.StatusCode)
	}

	var result struct {
		LiquidationRisk float64 `json:"liquidation_risk"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Check if risk exceeds threshold
	return result.LiquidationRisk > threshold, nil
}

// executeAction executes the action specified in the rule
func (e *Engine) executeAction(ctx context.Context, rule models.AutomationRule) error {
	switch rule.ActionType {
	case "rebalance":
		return e.executeRebalance(ctx, rule)
	case "withdraw":
		return e.executeWithdraw(ctx, rule)
	case "deposit":
		return e.executeDeposit(ctx, rule)
	default:
		return fmt.Errorf("unknown action type: %s", rule.ActionType)
	}
}

// executeRebalance executes a rebalancing action
func (e *Engine) executeRebalance(ctx context.Context, rule models.AutomationRule) error {
	config := rule.ActionConfig
	if config == nil {
		return fmt.Errorf("action config is nil")
	}

	fromProtocol, _ := config["from_protocol"].(string)
	toProtocol, _ := config["to_protocol"].(string)
	asset, _ := config["asset"].(string)
	amount, _ := config["amount"].(float64)

	log.Printf("Executing rebalance: %s from %s to %s, amount: %f %s",
		asset, fromProtocol, toProtocol, amount, asset)

	// Get user
	var user models.User
	if err := database.DB.First(&user, rule.UserID).Error; err != nil {
		return err
	}

	// TODO: Implement actual rebalancing logic
	// 1. Withdraw from source protocol
	// 2. Deposit to target protocol
	// 3. Record transaction

	// For now, just log the action
	log.Printf("Rebalancing executed for user %s", user.WalletAddress)

	return nil
}

// executeWithdraw executes a withdrawal action
func (e *Engine) executeWithdraw(ctx context.Context, rule models.AutomationRule) error {
	config := rule.ActionConfig
	if config == nil {
		return fmt.Errorf("action config is nil")
	}

	log.Printf("Executing withdraw action for rule %d", rule.ID)
	// TODO: Implement withdrawal logic
	return nil
}

// executeDeposit executes a deposit action
func (e *Engine) executeDeposit(ctx context.Context, rule models.AutomationRule) error {
	config := rule.ActionConfig
	if config == nil {
		return fmt.Errorf("action config is nil")
	}

	log.Printf("Executing deposit action for rule %d", rule.ID)
	// TODO: Implement deposit logic
	return nil
}

