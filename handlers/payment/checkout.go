package payment

import (
	"database/sql"
	"fmt"
	"kingdup/db"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

type CheckoutItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type CheckoutRequest struct {
	Items           []CheckoutItem `json:"items"`
	Email           string         `json:"email,omitempty"`
	IsGuest         bool           `json:"is_guest"`
	UserID          *int           `json:"user_id,omitempty"` // 👈 new
	ShippingAddress string         `json:"shipping_address"`
	ReturnURL       string         `json:"return_url"`
	CancelURL       string         `json:"cancel_url"`
}

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func CreateCheckoutHandler(sqlDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

		// Map items to OrderItemInput
		var orderItems []db.OrderItemInput
		for _, item := range req.Items {
			orderItems = append(orderItems, db.OrderItemInput{
				ProductName: item.Name,
				Quantity:    item.Quantity,
				PriceEach:   item.Price,
			})
		}

		ctx := c.Request.Context()

		// Save the order before redirecting to Stripe
		var orderID int64
		var err error

		orderID, err = db.CreateOrderWithItems(ctx, sqlDB, db.OrderInput{
			UserID:          toInt64Ptr(req.UserID),
			GuestEmail:      toPtr(req.Email),
			GuestName:       nil,
			Status:          "pending",
			TotalAmount:     calculateTotal(orderItems),
			ShippingAddress: req.ShippingAddress,
			Items:           orderItems,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order", "details": err.Error()})
			return
		}

		// Prepare Stripe Line Items
		var lineItems []*stripe.CheckoutSessionLineItemParams
		for _, item := range req.Items {
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(item.Name),
					},
					UnitAmount: stripe.Int64(int64(item.Price * 100)),
				},
				Quantity: stripe.Int64(int64(item.Quantity)),
			})
		}

		params := &stripe.CheckoutSessionParams{
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
			LineItems:          lineItems,
			Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL:         stripe.String(req.ReturnURL),
			CancelURL:          stripe.String(req.CancelURL),
			ClientReferenceID:  stripe.String(fmt.Sprintf("%d", orderID)), // optional but useful
		}

		if req.Email != "" {
			params.CustomerEmail = stripe.String(req.Email)
		}

		s, err := session.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": s.URL})
	}
}

func toInt64Ptr(i *int) *int64 {
	if i == nil {
		return nil
	}
	converted := int64(*i)
	return &converted
}

func calculateTotal(items []db.OrderItemInput) float64 {
	var total float64
	for _, item := range items {
		total += item.PriceEach * float64(item.Quantity)
	}
	return total
}
