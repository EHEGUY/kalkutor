package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// CalculationRequest represents the incoming calculation request
type CalculationRequest struct {
	Expression string  `json:"expression,omitempty"`
	Num1       float64 `json:"num1,omitempty"`
	Num2       float64 `json:"num2,omitempty"`
	Operation  string  `json:"operation,omitempty"`
}

// CalculationResponse represents the calculation result
type CalculationResponse struct {
	Result      float64 `json:"result"`
	Expression  string  `json:"expression"`
	Description string  `json:"description"`
	Success     bool    `json:"success"`
}

// ErrorResponse represents an error message
type ErrorResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

// enableCORS adds CORS headers to responses
func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// HomeHandler serves the welcome message
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Advanced Calculator API - Ready to compute!",
		"version": "2.0",
		"endpoints": "/calculate (POST), /health (GET)",
	})
}

// HealthHandler checks server health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "calculator-api",
	})
}

// CalculateHandler performs calculations
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Only POST method is allowed",
			Success: false,
		})
		return
	}
	
	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid request format",
			Success: false,
		})
		return
	}
	
	// Handle expression-based calculation
	if req.Expression != "" {
		result, desc, err := evaluateExpression(req.Expression)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   err.Error(),
				Success: false,
			})
			return
		}
		
		json.NewEncoder(w).Encode(CalculationResponse{
			Result:      result,
			Expression:  req.Expression,
			Description: desc,
			Success:     true,
		})
		return
	}
	
	// Handle traditional operation-based calculation
	result, desc, err := performOperation(req.Num1, req.Num2, req.Operation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   err.Error(),
			Success: false,
		})
		return
	}
	
	expression := fmt.Sprintf("%.2f %s %.2f", req.Num1, req.Operation, req.Num2)
	json.NewEncoder(w).Encode(CalculationResponse{
		Result:      result,
		Expression:  expression,
		Description: desc,
		Success:     true,
	})
}

// performOperation executes the mathematical operation
func performOperation(num1, num2 float64, op string) (float64, string, error) {
	var result float64
	var description string
	
	switch op {
	case "+", "add":
		result = num1 + num2
		description = "Addition completed"
	case "-", "subtract":
		result = num1 - num2
		description = "Subtraction completed"
	case "*", "multiply":
		result = num1 * num2
		description = "Multiplication completed"
	case "/", "divide":
		if num2 == 0 {
			return 0, "", fmt.Errorf("cannot divide by zero")
		}
		result = num1 / num2
		description = "Division completed"
	case "%", "mod":
		if num2 == 0 {
			return 0, "", fmt.Errorf("cannot perform modulo with zero")
		}
		result = math.Mod(num1, num2)
		description = "Modulo completed"
	case "^", "power":
		result = math.Pow(num1, num2)
		description = "Power calculation completed"
	default:
		return 0, "", fmt.Errorf("unsupported operation: %s", op)
	}
	
	return result, description, nil
}

// evaluateExpression parses and evaluates a mathematical expression
func evaluateExpression(expr string) (float64, string, error) {
	expr = strings.TrimSpace(expr)
	
	// Simple expression parser (supports basic operations)
	// For production, consider using a proper parser library
	
	// Find the operation
	ops := []string{"+", "-", "*", "/", "%", "^"}
	for _, op := range ops {
		idx := strings.LastIndex(expr, op)
		if idx > 0 && idx < len(expr)-1 {
			leftStr := strings.TrimSpace(expr[:idx])
			rightStr := strings.TrimSpace(expr[idx+1:])
			
			left, err1 := strconv.ParseFloat(leftStr, 64)
			right, err2 := strconv.ParseFloat(rightStr, 64)
			
			if err1 == nil && err2 == nil {
				result, desc, err := performOperation(left, right, op)
				return result, desc, err
			}
		}
	}
	
	// Try to parse as single number
	num, err := strconv.ParseFloat(expr, 64)
	if err == nil {
		return num, "Value parsed", nil
	}
	
	return 0, "", fmt.Errorf("invalid expression format")
}

func main() {
	mux := http.NewServeMux()
	
	// Register handlers
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/calculate", CalculateHandler)
	
	// Server configuration
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
	log.Println("🚀 Calculator API server starting on port 8080")
	log.Println("📍 Endpoints:")
	log.Println("   GET  /        - API information")
	log.Println("   GET  /health  - Health check")
	log.Println("   POST /calculate - Perform calculations")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
