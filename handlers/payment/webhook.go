package payment

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

func StripeWebhookHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Read error"})
			return
		}

		sigHeader := c.GetHeader("Stripe-Signature")
		endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

		event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Signature verification failed"})
			return
		}

		if event.Type == "checkout.session.completed" {
			var session stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook JSON parse error"})
				return
			}

			orderID := session.ClientReferenceID
			if orderID != "" {
				_, err := db.Exec(`UPDATE orders SET status = 'paid' WHERE id = $1`, orderID)
				if err != nil {
					fmt.Println("❌ Failed to update order status:", err)
				} else {
					fmt.Println("✅ Order", orderID, "marked as paid")
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"received": true})
	}
}
