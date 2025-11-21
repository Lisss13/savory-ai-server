package service

import (
	"fmt"
	"strings"
)

func generateBotResponse(userMessage string) (string, error) {

	// Simple responses based on keywords
	switch {
	case containsAny(userMessage, []string{"hello", "hi", "hey", "greetings"}):
		return fmt.Sprintf("Hello! Welcome to table %s. How can I help you today?", "Arseniy"), nil
	case containsAny(userMessage, []string{"menu", "food", "eat", "dish", "meal"}):
		return "I can help you with our menu. Would you like to see our specials or the full menu?", nil
	case containsAny(userMessage, []string{"order", "want", "get", "bring"}):
		return "I'll be happy to take your order. What would you like to have?", nil
	case containsAny(userMessage, []string{"bill", "check", "pay", "payment"}):
		return "I'll arrange for your bill right away. You can pay at the counter or through our app.", nil
	case containsAny(userMessage, []string{"thank", "thanks", "appreciate"}):
		return "You're welcome! Is there anything else I can help you with?", nil
	case containsAny(userMessage, []string{"bye", "goodbye", "see you", "leaving"}):
		return "Thank you for visiting us! We hope to see you again soon.", nil
	default:
		return "I'm here to assist you with your dining experience. Can you please provide more details about what you need?", nil
	}
}

// Helper function to check if a string contains any of the given keywords
func containsAny(s string, keywords []string) bool {
	for _, keyword := range keywords {
		if contains(s, keyword) {
			return true
		}
	}
	return false
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
