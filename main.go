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
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GitHubWebhookPayload struct {
	Action string `json:"action"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		Name string `json:"name"`
	} `json:"repository"`
}

func handleGitHubWebhook(c echo.Context) error {
	// Read raw request body
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
	}

	// Print the raw body (useful for debugging)
	fmt.Println("Raw GitHub Webhook Body:", string(body))

	// Parse payload
	payload := new(GitHubWebhookPayload)
	if err := c.Bind(payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid webhook payload"})
	}

	// Log event details
	fmt.Printf("📢 Received GitHub Event: Repository=%s, Pusher=%s, Action=%s\n",
		payload.Repository.Name, payload.Pusher.Name, payload.Action)

	return c.JSON(http.StatusOK, map[string]string{"message": "GitHub Webhook received successfully"})
}

func main() {
	e := echo.New()

	// Define the GitHub webhook endpoint
	e.POST("/github-webhook", handleGitHubWebhook)

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
