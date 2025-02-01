// без ошибок: curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d "{\"expression\": \"2+2*2\"}" -i
// ошибка 422: curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d "{\"expression\": \"\"}" -i
// ошибка 500: curl -X POST "http://localhost:8080/api/v1/calculate" -H "Content-Type: application/json" -d \"\" -i

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"calculator/pkg/calculator"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculateHandler(writer http.ResponseWriter, request *http.Request) {
	var userInput CalculationRequest

	var body = json.NewDecoder(request.Body)
	var err = body.Decode(&userInput)

	if err != nil {
		http.Error(writer, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	var expression = userInput.Expression

	if expression == "" {
		http.Error(writer, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculator.Calc(expression)

	resultString := strconv.FormatFloat(result, 'f', -1, 64)

	if err != nil {
		http.Error(writer, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(CalculationResponse{Result: resultString})
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
