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

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result      float64 `json:"result"`
	Success     bool    `json:"success"`
	Description string  `json:"description"`
}

// enableCORS allows the browser to talk to the server
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, desc, err := evaluateExpression(req.Expression)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CalculationResponse{
		Result:      result,
		Success:     err == nil,
		Description: desc,
	})
}

// evaluateExpression logic from your original code
func evaluateExpression(expr string) (float64, string, error) {
	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")

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

	num, err := strconv.ParseFloat(expr, 64)
	if err == nil {
		return num, "Value parsed", nil
	}
	return 0, "", fmt.Errorf("invalid format")
}

// performOperation logic from your original code
func performOperation(num1, num2 float64, op string) (float64, string, error) {
	switch op {
	case "+":
		return num1 + num2, "Addition completed", nil
	case "-":
		return num1 - num2, "Subtraction completed", nil
	case "*":
		return num1 * num2, "Multiplication completed", nil
	case "/":
		if num2 == 0 {
			return 0, "", fmt.Errorf("cannot divide by zero")
		}
		return num1 / num2, "Division completed", nil
	case "%":
		return math.Mod(num1, num2), "Modulo completed", nil
	case "^":
		return math.Pow(num1, num2), "Power completed", nil
	default:
		return 0, "", fmt.Errorf("unsupported op")
	}
}

func main() {
	http.HandleFunc("/calculate", CalculateHandler)
	fmt.Println("🚀 Apple-Style Calc Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
