package twilio

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func SendMessage(ctx context.Context, to, body string) error {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	whatsappNumber := os.Getenv("TWILIO_WHATSAPP_NUMBER")

	if accountSID == "" || authToken == "" || whatsappNumber == "" {
		return fmt.Errorf("twilio credentials missing: ensure environment variables are loaded")
	}

	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	formData := url.Values{}
	formData.Set("To", "whatsapp:"+to)
	formData.Set("From", whatsappNumber)
	formData.Set("Body", body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("request timed out or cancel: %w", err)
		}
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Message sent successfully: %s\n", string(respBody))
		return nil
	}
	return fmt.Errorf("failed to send message, status: %d, response: %s", resp.StatusCode, string(respBody))
}
