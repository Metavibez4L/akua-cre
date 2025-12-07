package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// TestWebhookIntegration tests the full HTTP webhook delivery flow
func TestWebhookIntegration(t *testing.T) {
	receivedPayloads := []EventPayload{}
	var mu sync.Mutex

	// Create mock webhook server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify webhook secret header
		secret := r.Header.Get("X-CRE-Webhook-Secret")
		if secret != "test-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid secret"})
			return
		}

		// Verify Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got: %s", contentType)
		}

		// Parse payload
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var payload EventPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Errorf("failed to unmarshal webhook payload: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Store received payload
		mu.Lock()
		receivedPayloads = append(receivedPayloads, payload)
		mu.Unlock()

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      true,
			"eventId": "test-event-id",
		})
	}))
	defer server.Close()

	// Test sending various event types
	tests := []struct {
		name    string
		payload EventPayload
	}{
		{
			name: "EscrowCreated",
			payload: EventPayload{
				EventName:      "EscrowCreated",
				Network:        "arbitrumSepolia",
				ChainName:      "ethereum-testnet-sepolia-arbitrum-1",
				EscrowID:       "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
				Seller:         "0x1234567890123456789012345678901234567890",
				Buyer:          "0x0987654321098765432109876543210987654321",
				CommodityToken: "0xabcdef1234567890abcdef1234567890abcdef12",
				TokenAmount:    "1000",
				PricePerToken:  "50",
				TxHash:         "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				BlockNumber:    "12345678",
				Timestamp:      1234567890,
			},
		},
		{
			name: "EscrowFunded",
			payload: EventPayload{
				EventName:   "EscrowFunded",
				Network:     "arbitrumSepolia",
				ChainName:   "ethereum-testnet-sepolia-arbitrum-1",
				EscrowID:    "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
				Funder:      "0x1234567890123456789012345678901234567890",
				TokenAmount: "50000",
				TxHash:      "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				BlockNumber: "12345679",
				Timestamp:   1234567900,
			},
		},
		{
			name: "ReceiptIssued",
			payload: EventPayload{
				EventName:     "ReceiptIssued",
				Network:       "arbitrumSepolia",
				ChainName:     "ethereum-testnet-sepolia-arbitrum-1",
				ReceiptID:     "0x201f1e1d1c1b1a191817161514131211100f0e0d0c0b0a090807060504030201",
				Depositor:     "0x1234567890123456789012345678901234567890",
				Warehouse:     "0xabcdef1234567890abcdef1234567890abcdef12",
				CommodityType: "Cocoa Beans",
				Quantity:      "10000",
				DepositDate:   "1640000000",
				ExpiryDate:    "1672000000",
				TxHash:        "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
				BlockNumber:   "12345680",
				Timestamp:     1234567910,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send webhook request
			payload, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(payload))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-CRE-Webhook-Secret", "test-secret")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status 200, got: %d", resp.StatusCode)
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if ok, exists := result["ok"].(bool); !exists || !ok {
				t.Error("expected 'ok' field to be true")
			}
		})
	}

	// Verify all payloads were received
	mu.Lock()
	defer mu.Unlock()

	if len(receivedPayloads) != len(tests) {
		t.Errorf("expected %d webhook calls, got %d", len(tests), len(receivedPayloads))
	}

	// Verify each payload has required fields
	for i, payload := range receivedPayloads {
		if payload.EventName == "" {
			t.Errorf("payload %d: eventName should not be empty", i)
		}
		if payload.Network == "" {
			t.Errorf("payload %d: network should not be empty", i)
		}
		if payload.ChainName == "" {
			t.Errorf("payload %d: chainName should not be empty", i)
		}
		if payload.TxHash == "" {
			t.Errorf("payload %d: txHash should not be empty", i)
		}
	}
}

// TestWebhookAuthentication tests webhook authentication scenarios
func TestWebhookAuthentication(t *testing.T) {
	tests := []struct {
		name           string
		secret         string
		expectedStatus int
	}{
		{
			name:           "Valid Secret",
			secret:         "correct-secret",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Secret",
			secret:         "wrong-secret",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing Secret",
			secret:         "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("X-CRE-Webhook-Secret")
		if secret != "correct-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := EventPayload{
				EventName: "TestEvent",
				Network:   "arbitrumSepolia",
				ChainName: "ethereum-testnet-sepolia-arbitrum-1",
			}

			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.secret != "" {
				req.Header.Set("X-CRE-Webhook-Secret", tt.secret)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got: %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

// TestWebhookErrorHandling tests webhook error scenarios
func TestWebhookErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		serverStatus   int
		serverResponse string
		expectError    bool
	}{
		{
			name:           "Success Response",
			serverStatus:   http.StatusOK,
			serverResponse: `{"ok": true}`,
			expectError:    false,
		},
		{
			name:           "Bad Request",
			serverStatus:   http.StatusBadRequest,
			serverResponse: `{"error": "Invalid payload"}`,
			expectError:    true,
		},
		{
			name:           "Internal Server Error",
			serverStatus:   http.StatusInternalServerError,
			serverResponse: `{"error": "Server error"}`,
			expectError:    true,
		},
		{
			name:           "Service Unavailable",
			serverStatus:   http.StatusServiceUnavailable,
			serverResponse: `{"error": "Service unavailable"}`,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			payload := EventPayload{
				EventName: "TestEvent",
				Network:   "arbitrumSepolia",
				ChainName: "ethereum-testnet-sepolia-arbitrum-1",
			}

			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if tt.expectError {
				if resp.StatusCode < 400 {
					t.Error("expected error status code (>= 400)")
				}
			} else {
				if resp.StatusCode >= 400 {
					t.Errorf("expected success status code (< 400), got: %d", resp.StatusCode)
				}
			}
		})
	}
}

// TestWebhookPayloadValidation tests payload validation
func TestWebhookPayloadValidation(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		expectSuccess bool
	}{
		{
			name: "Valid Payload",
			payload: `{
				"eventName": "EscrowCreated",
				"network": "arbitrumSepolia",
				"chainName": "ethereum-testnet-sepolia-arbitrum-1",
				"escrowId": "0x01020304..."
			}`,
			expectSuccess: true,
		},
		{
			name: "Missing EventName",
			payload: `{
				"network": "arbitrumSepolia",
				"chainName": "ethereum-testnet-sepolia-arbitrum-1"
			}`,
			expectSuccess: false,
		},
		{
			name: "Missing Network",
			payload: `{
				"eventName": "EscrowCreated",
				"chainName": "ethereum-testnet-sepolia-arbitrum-1"
			}`,
			expectSuccess: false,
		},
		{
			name:          "Invalid JSON",
			payload:       `{invalid json`,
			expectSuccess: false,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		
		var payload EventPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}

		if payload.EventName == "" || payload.Network == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "eventName and network are required",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if tt.expectSuccess {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("expected status 200, got: %d", resp.StatusCode)
				}
			} else {
				if resp.StatusCode == http.StatusOK {
					t.Error("expected error status code, got 200")
				}
			}
		})
	}
}

