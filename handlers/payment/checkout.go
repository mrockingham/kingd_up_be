package payment

import (
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
	Items     []CheckoutItem `json:"items"`
	Email     string         `json:"email,omitempty"`
	IsGuest   bool           `json:"is_guest"`
	ReturnURL string         `json:"return_url"`
	CancelURL string         `json:"cancel_url"`
}

func CreateCheckoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

		var lineItems []*stripe.CheckoutSessionLineItemParams
		for _, item := range req.Items {
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(item.Name),
					},
					UnitAmount: stripe.Int64(int64(item.Price * 100)), // cents
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
