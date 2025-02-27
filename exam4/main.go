package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type DateRequest struct {
	Input string `json:"input"`
}

type DateResponse struct {
	Year  string `json:"year"`
	Month string `json:"month"`
	Day   string `json:"day"`
}

func parseDate(input string) (DateResponse, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "คุณเป็น AI ที่ช่วยแปลงข้อมูลวันเดือนปีให้อยู่ในรูปแบบที่ถูกต้อง (YYYY-MM-DD)",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("แปลงข้อมูลวันที่: %s", input),
				},
			},
		},
	)
	if err != nil {
		return DateResponse{}, err
	}

	output := strings.TrimSpace(resp.Choices[0].Message.Content)
	dateParts := strings.Split(output, "-")

	if len(dateParts) != 3 {
		return DateResponse{"-", "-", "-"}, fmt.Errorf("invalid date format")
	}

	return DateResponse{
		Year:  dateParts[0],
		Month: dateParts[1],
		Day:   dateParts[2],
	}, nil
}

func dateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dateResp, err := parseDate(req.Input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dateResp)
}

func main() {
	http.HandleFunc("/parse-date", dateHandler)
	port := ":8080"
	fmt.Println("Server is running on port" + port)
	http.ListenAndServe(port, nil)
	
}
