package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

// InitStripe initializes Stripe with API key
func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// CreateCheckoutSession creates a Stripe checkout session
func CreateCheckoutSession(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Tier string `json:"tier" binding:"required"` // basic, premium
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Determine price based on tier
	var priceID string
	switch req.Tier {
	case "basic":
		priceID = os.Getenv("STRIPE_BASIC_PRICE_ID")
	case "premium":
		priceID = os.Getenv("STRIPE_PREMIUM_PRICE_ID")
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tier"})
		return
	}

	// Create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(os.Getenv("FRONTEND_URL") + "/settings?success=true"),
		CancelURL:  stripe.String(os.Getenv("FRONTEND_URL") + "/settings?canceled=true"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: stripe.String(user.Email),
		Metadata: map[string]string{
			"user_id": string(rune(user.ID)),
			"tier":    req.Tier,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sess.ID,
		"url":        sess.URL,
	})
}

// HandleStripeWebhook handles Stripe webhook events
func HandleStripeWebhook(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEvent(body, sigHeader, webhookSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		// Handle successful checkout
		// Update user subscription in database
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
			// Update subscription based on session
			// This is a simplified version
		}

	case "customer.subscription.updated":
		// Handle subscription update
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err == nil {
			// Update subscription in database
		}

	case "customer.subscription.deleted":
		// Handle subscription cancellation
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err == nil {
			// Cancel subscription in database
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

