package services

import (
	"encoding/json"
	"fmt"
	"io"
	"kingdup/db"
	"net/http"
	"os"
)

type PrintfulListItem struct {
	ID int64 `json:"id"`
}

type PrintfulListResponse struct {
	Result []PrintfulListItem `json:"result"`
}

type SyncProduct struct {
	ID           int64  `json:"id"`
	ExternalID   string `json:"external_id"`
	Name         string `json:"name"`
	Synced       int    `json:"synced"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type PrintfulVariant struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	RetailPrice string `json:"retail_price"`
	Files       []struct {
		Type       string `json:"type"`
		PreviewURL string `json:"preview_url"`
	} `json:"files"`
}

type PrintfulProduct struct {
	ID           int64             `json:"id"`
	ExternalID   string            `json:"external_id"`
	Name         string            `json:"name"`
	Synced       int               `json:"synced"`
	ThumbnailURL string            `json:"thumbnail_url"`
	Variants     []PrintfulVariant `json:"variants"`
}

type PrintfulProductResponse struct {
	Result struct {
		SyncProduct  SyncProduct       `json:"sync_product"`
		SyncVariants []PrintfulVariant `json:"sync_variants"`
	} `json:"result"`
}

func SyncProductsFromPrintful() error {
	apiKey := os.Getenv("PRINTFUL_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("missing PRINTFUL_API_KEY")
	}

	baseURL := "https://api.printful.com"

	// Step 1: Get list of product IDs
	listURL := baseURL + "/store/products"

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch product list, status code: %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)

	var listResp PrintfulListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return err
	}

	// Step 2: Fetch each product's full info
	for _, item := range listResp.Result {
		detailURL := fmt.Sprintf("%s/store/products/%d", baseURL, item.ID)

		req, _ := http.NewRequest("GET", detailURL, nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Error fetching product %d: %v\n", item.ID, err)
			continue
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			fmt.Printf("❌ Failed to fetch product %d details\n", item.ID)
			continue
		}

		body, _ := io.ReadAll(res.Body)

		var productResp PrintfulProductResponse
		if err := json.Unmarshal(body, &productResp); err != nil {
			fmt.Printf("❌ Failed to parse product %d JSON: %v\n", item.ID, err)
			continue
		}

		p := productResp.Result.SyncProduct

		product := db.Product{
			PrintfulID:   p.ID,
			ExternalID:   p.ExternalID,
			Name:         p.Name,
			ThumbnailURL: p.ThumbnailURL,
			Synced:       p.Synced == 1,
		}

		// Upsert product
		err = db.DB.
			Where("printful_id = ?", p.ID).
			Assign(product).
			FirstOrCreate(&product).Error

		if err != nil {
			fmt.Printf("❌ Failed to save product %d: %v\n", p.ID, err)
			continue
		}

		// Upsert variants
		for _, v := range productResp.Result.SyncVariants {
			thumb := ""
			for _, file := range v.Files {
				if file.Type == "preview" && file.PreviewURL != "" {
					thumb = file.PreviewURL
					break
				}
			}

			variant := db.Variant{
				PrintfulID:   v.ID,
				ProductID:    product.ID,
				Name:         v.Name,
				SKU:          v.SKU,
				RetailPrice:  v.RetailPrice,
				ThumbnailURL: thumb,
			}

			err := db.DB.
				Where("printful_id = ?", v.ID).
				Assign(variant).
				FirstOrCreate(&variant).Error

			if err != nil {
				fmt.Printf("❌ Failed to save variant %d: %v\n", v.ID, err)
			}
		}

	}

	fmt.Println("✅ Synced all Printful products and variants")
	return nil
}
