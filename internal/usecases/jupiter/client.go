package jupiter

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/humangrass/price-keeper/config"

	"github.com/mailru/easyjson"
)

type Client struct {
	cfg     config.Jupiter
	httpCli *http.Client
}

func NewClient(cfg config.Jupiter) *Client {
	return &Client{
		cfg: cfg,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetPrices() error {
	url := c.buildURL()

	resp, err := c.doRequest(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := c.validateResponse(resp); err != nil {
		return err
	}

	priceResp, err := c.decodeResponse(resp)
	if err != nil {
		return err
	}

	return c.processResponse(priceResp)
}

func (c *Client) buildURL() string {
	url := fmt.Sprintf("https://lite-api.jup.ag/price/v2?ids=%s", c.cfg.TokenID)
	if !c.cfg.ExtraInfo && len(c.cfg.VSTokenID) > 0 {
		url += fmt.Sprintf("&vsToken=%s", c.cfg.VSTokenID)
	}
	return url
}

func (c *Client) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	return c.httpCli.Do(req)
}

func (c *Client) validateResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) decodeResponse(resp *http.Response) (*PriceResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var priceResp PriceResponse
	if err := easyjson.Unmarshal(body, &priceResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &priceResp, nil
}

func (c *Client) processResponse(resp *PriceResponse) error {
	for token, data := range resp.Data {
		log.Printf("%s - %+v\n", token, data)
	}
	log.Printf("Request took %.4f seconds", resp.TimeTaken)
	return nil
}
