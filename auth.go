package gobriva

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"
)

// authenticate performs OAuth2 authentication to get access token
func (c *Client) authenticate(ctx context.Context) error {
	// Create signature for token request
	timestamp := c.generateTimestamp()
	payload := c.clientID + "|" + timestamp

	// Parse private key
	block, _ := pem.Decode([]byte(c.privateKey))
	if block == nil {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	var privateKey *rsa.PrivateKey
	var err error
	if parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		privateKey = parsedKey
	} else if parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
			privateKey = rsaKey
		} else {
			return fmt.Errorf("private key is not RSA")
		}
	} else {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create signature
	hashed := sha256.Sum256([]byte(payload))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return fmt.Errorf("failed to sign payload: %w", err)
	}

	// Encode signature to base64
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	// Create token request
	tokenReq := TokenRequest{
		GrantType: "client_credentials",
	}

	reqBody, err := json.Marshal(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to marshal token request: %w", err)
	}

	// Create HTTP request
	fullURL := c.baseURL + "/snap/v1.0/access-token/b2b"
	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SIGNATURE", signatureB64)
	req.Header.Set("X-CLIENT-KEY", c.clientID)
	req.Header.Set("X-TIMESTAMP", timestamp)

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	// Parse response
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		json.Unmarshal(respBody, &errorResp)
		return &APIError{
			ResponseCode:    errorResp.ResponseCode,
			ResponseMessage: errorResp.ResponseMessage,
		}
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return fmt.Errorf("failed to unmarshal token response: %w", err)
	}

	// Store token
	c.accessToken = authResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)

	return nil
}
