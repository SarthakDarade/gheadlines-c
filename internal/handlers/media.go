package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gheadlines/config"
	"io"
	"mime/multipart"
	"net/http"
)

// CloudflareResponse represents the response from Cloudflare Images API
type CloudflareResponse struct {
	Success bool `json:"success"`
	Result  struct {
		ID       string   `json:"id"`
		Filename string   `json:"filename"`
		Variants []string `json:"variants"`
	} `json:"result"`
	Errors []interface{} `json:"errors"`
}

// UploadImageHandler handles image uploads to Cloudflare
func UploadImageHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload to Cloudflare
		imageURL, err := uploadToCloudflare(file, header, cfg)
		if err != nil {
			http.Error(w, "Error uploading to Cloudflare: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"location": imageURL, // TinyMCE expects "location"
			"url":      imageURL,
		})
	}
}

func uploadToCloudflare(file multipart.File, header *multipart.FileHeader, cfg *config.Config) (string, error) {
	// If no Cloudflare config, return a dummy URL for testing
	if cfg.CloudflareAccountID == "" || cfg.CloudflareImagesToken == "" {
		return "https://images.unsplash.com/photo-1504711434969-e33886168f5c?w=800&q=80", nil
	}

	// Prepare multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	// Create request
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/images/v1", cfg.CloudflareAccountID)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.CloudflareImagesToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read body for error details
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("cloudflare API error: %s - %s", resp.Status, string(respBody))
	}

	// Parse response
	var cfResp CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return "", err
	}

	if !cfResp.Success {
		return "", fmt.Errorf("cloudflare upload failed: %v", cfResp.Errors)
	}

	// Return the first variant (usually 'public' or similar)
	if len(cfResp.Result.Variants) > 0 {
		return cfResp.Result.Variants[0], nil
	}

	return "", fmt.Errorf("no image variants returned")
}
