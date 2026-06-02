package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/arkag/cope-and-hope-weather/models"
)

func TestLiveAPI(t *testing.T) {
	apiEndpoint := os.Getenv("API_ENDPOINT")
	if apiEndpoint == "" {
		t.Skip("Skipping E2E test: API_ENDPOINT not set")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	t.Run("Cope Mode Live Test", func(t *testing.T) {
		req, err := http.NewRequest("GET", apiEndpoint+"/weather?city=Seattle&mode=cope", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var apiResp models.APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			t.Fatalf("failed to decode json: %v", err)
		}

		if apiResp.Mode != "cope" {
			t.Errorf("expected mode cope, got %s", apiResp.Mode)
		}
		if apiResp.Requested.City != "Seattle" {
			t.Errorf("expected city Seattle, got %s", apiResp.Requested.City)
		}
	})
}
