//package main
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//
//	"github.com/labstack/echo/v4"
//)
//
//// WebhookPayload struct to parse incoming webhook data
//type WebhookPayload struct {
//	Event   string `json:"event"`
//	OrderID string `json:"order_id"`
//	Status  string `json:"status"`
//}
//
//func handleWebhook(c echo.Context) error {
//	payload := new(WebhookPayload)
//	if err := c.Bind(payload); err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
//	}
//
//	// Process the webhook event
//	fmt.Printf("Received Webhook: Event=%s, OrderID=%s, Status=%s\n", payload.Event, payload.OrderID, payload.Status)
//
//	// Simulate updating order status in the database
//	if payload.Event == "payment_success" {
//		fmt.Printf("✅ Order %s marked as PAID\n", payload.OrderID)
//	}
//
//	return c.JSON(http.StatusOK, map[string]string{"message": "Webhook received"})
//}
//
//func main() {
//	e := echo.New()
//	e.POST("/webhook", handleWebhook)
//
//	fmt.Println("🚀 Webhook server running on port 8080")
//	log.Fatal(e.Start(":8080"))
//}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// GitHubWebhookPayload defines the expected structure of the webhook payload
type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	After      string `json:"after"`
	Before     string `json:"before"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Print the raw payload
	fmt.Println("Raw GitHub Webhook Body:", string(body))

	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Process the webhook data
	fmt.Printf("Repository: %s\n", payload.Repository.FullName)
	fmt.Printf("Pusher: %s (%s)\n", payload.Pusher.Name, payload.Pusher.Email)
	fmt.Printf("Branch: %s\n", payload.Ref)
	fmt.Println("Commits:")
	for _, commit := range payload.Commits {
		fmt.Printf("- %s: %s (%s)\n", commit.ID[:7], commit.Message, commit.Author.Name)
		fmt.Printf("  Commit URL: %s\n", commit.URL)
	}

	// Respond to GitHub
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	port := 8080
	fmt.Printf("⇨ HTTP server started on :%d\n", port)
	// Vishal
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
