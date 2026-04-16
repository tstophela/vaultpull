package vault

import (
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// Config holds the parameters needed to create a Vault client.
type Config struct {
	Address string
	Token   string
	Timeout time.Duration
}

// NewClient creates and authenticates a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	vCfg := vaultapi.DefaultConfig()
	vCfg.Address = cfg.Address
	vCfg.HttpClient = &http.Client{Timeout: cfg.Timeout}

	client, err := vaultapi.NewClient(vCfg)
	if err != nil {
		return nil, fmt.Errorf("vault: failed to create client: %w", err)
	}

	client.SetToken(cfg.Token)

	return &Client{logical: client.Logical()}, nil
}

// ReadSecrets reads key/value secrets from the given Vault path.
// It supports both KV v1 and KV v2 (data/ prefix handling is caller's responsibility).
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault: no secret found at %q", path)
	}

	data, ok := secret.Data["data"]
	if !ok {
		// KV v1 — data is at the top level
		data = secret.Data
	}

	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: unexpected data format at %q", path)
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
