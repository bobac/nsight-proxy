package nsight

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/net/html/charset"
)

// ApiClient holds the configuration and provides methods for API calls
type ApiClient struct {
	apiKey string
	server string
}

// NewApiClient creates a new ApiClient, loading configuration from .env
func NewApiClient() (*ApiClient, error) {
	err := godotenv.Load()
	if err != nil {
		// It's often okay if .env is missing, rely on environment variables
		log.Println("Warning: Could not load .env file:", err)
	}

	apiKey := os.Getenv("NSIGHT_API_KEY")
	server := os.Getenv("NSIGHT_SERVER")

	if apiKey == "" || server == "" {
		return nil, errors.New("NSIGHT_API_KEY and NSIGHT_SERVER must be set in .env file or environment variables")
	}

	return &ApiClient{apiKey: apiKey, server: server}, nil
}

// callAPI performs the HTTP GET request and returns the response body bytes
func (c *ApiClient) callAPI(service string, params map[string]string) ([]byte, error) {
	base, err := url.Parse(fmt.Sprintf("https://%s/api/", c.server))
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	q := base.Query()
	q.Set("apikey", c.apiKey)
	q.Set("service", service)
	for key, value := range params {
		q.Set(key, value)
	}
	base.RawQuery = q.Encode()

	apiUrl := base.String()
	fmt.Println("Requesting URL:", apiUrl) // Print URL for debugging

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from API (%s): %w", service, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body (%s): %w", service, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API (%s) returned non-OK status: %s\nResponse: %s", service, resp.Status, string(bodyBytes))
	}
	return bodyBytes, nil
}

// decodeXML parses the XML body using the correct charset reader
func decodeXML(bodyBytes []byte, target interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(target)
	if err != nil && err != io.EOF { // Ignore EOF if the structure allows empty results
		return fmt.Errorf("error parsing XML response: %w\nRaw Response:\n%s", err, string(bodyBytes))
	}
	return nil
}

// FetchClients fetches the list of all clients
func (c *ApiClient) FetchClients() ([]Client, error) {
	bodyBytes, err := c.callAPI("list_clients", nil)
	if err != nil {
		return nil, err
	}
	var result ClientResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchSites fetches the list of sites for a given client ID
func (c *ApiClient) FetchSites(clientID int) ([]Site, error) {
	params := map[string]string{"clientid": fmt.Sprintf("%d", clientID)}
	bodyBytes, err := c.callAPI("list_sites", params)
	if err != nil {
		// Consider if a specific error means 'no sites' vs. a real problem
		// For now, assume any error is problematic for fetching
		return nil, fmt.Errorf("fetching sites for client %d: %w", clientID, err)
	}
	var result SiteResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		// Check if the error suggests an empty list vs. malformed XML
		// This simple check might need refinement based on API behavior for empty lists
		if len(result.Items) == 0 && err.Error() == "EOF" { // Common case for empty list with Decode
			log.Printf("No sites found for client %d.", clientID)
			return []Site{}, nil // Return empty slice if no sites
		}
		return nil, fmt.Errorf("decoding sites for client %d: %w", clientID, err)
	}
	return result.Items, nil
}

// FetchServers fetches the list of servers for a given site ID
func (c *ApiClient) FetchServers(siteID int) ([]Server, error) {
	params := map[string]string{"siteid": fmt.Sprintf("%d", siteID)}
	bodyBytes, err := c.callAPI("list_servers", params)
	if err != nil {
		return nil, fmt.Errorf("fetching servers for site %d: %w", siteID, err)
	}
	var result ServerResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		if len(result.Items) == 0 && err.Error() == "EOF" {
			log.Printf("No servers found for site %d.", siteID)
			return []Server{}, nil
		}
		return nil, fmt.Errorf("decoding servers for site %d: %w", siteID, err)
	}
	return result.Items, nil
}

// FetchWorkstations fetches the list of workstations for a given site ID
func (c *ApiClient) FetchWorkstations(siteID int) ([]Workstation, error) {
	params := map[string]string{"siteid": fmt.Sprintf("%d", siteID)}
	bodyBytes, err := c.callAPI("list_workstations", params)
	if err != nil {
		return nil, fmt.Errorf("fetching workstations for site %d: %w", siteID, err)
	}
	var result WorkstationResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		if len(result.Items) == 0 && err.Error() == "EOF" {
			log.Printf("No workstations found for site %d.", siteID)
			return []Workstation{}, nil
		}
		return nil, fmt.Errorf("decoding workstations for site %d: %w", siteID, err)
	}
	return result.Items, nil
}
