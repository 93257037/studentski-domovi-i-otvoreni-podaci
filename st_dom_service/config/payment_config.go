package config

import (
	"os"
	"strconv"
)

// PaymentConfig holds payment-related configuration
type PaymentConfig struct {
	DefaultAmount        float64
	DefaultDueDay        int  // Day of month for payment due date (e.g., 15)
	AutoCreateOnApproval bool // Whether to auto-create payment when approving
}

// GetPaymentConfig returns the payment configuration
// Values can be overridden via environment variables
func GetPaymentConfig() PaymentConfig {
	config := PaymentConfig{
		DefaultAmount:        100.0, // Default €100
		DefaultDueDay:        15,    // Default: 15th of the month
		AutoCreateOnApproval: true,  // Default: auto-create enabled
	}

	// Override from environment if set
	if amountStr := os.Getenv("PAYMENT_DEFAULT_AMOUNT"); amountStr != "" {
		if amount, err := strconv.ParseFloat(amountStr, 64); err == nil && amount > 0 {
			config.DefaultAmount = amount
		}
	}

	if dueDayStr := os.Getenv("PAYMENT_DEFAULT_DUE_DAY"); dueDayStr != "" {
		if dueDay, err := strconv.Atoi(dueDayStr); err == nil && dueDay >= 1 && dueDay <= 31 {
			config.DefaultDueDay = dueDay
		}
	}

	if autoCreateStr := os.Getenv("PAYMENT_AUTO_CREATE"); autoCreateStr != "" {
		if autoCreate, err := strconv.ParseBool(autoCreateStr); err == nil {
			config.AutoCreateOnApproval = autoCreate
		}
	}

	return config
}

