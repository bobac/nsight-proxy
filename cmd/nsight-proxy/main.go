package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"nsight-proxy/internal/nsight"
)

// ProxyServer handles API requests
type ProxyServer struct {
	server string
}

// NewProxyServer creates a new proxy server instance
func NewProxyServer() (*ProxyServer, error) {
	// Only require server configuration, API key will come from requests
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	server := os.Getenv("NSIGHT_SERVER")
	if server == "" {
		return nil, fmt.Errorf("NSIGHT_SERVER must be set in .env file or environment variables")
	}
	
	return &ProxyServer{server: server}, nil
}

// handleAPI routes API requests based on service parameter
func (ps *ProxyServer) handleAPI(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		return
	}

	// Only accept GET requests
	if r.Method != "GET" {
		http.Error(w, `{"error": "Only GET method is supported"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract required parameters
	service := r.URL.Query().Get("service")
	if service == "" {
		http.Error(w, `{"error": "Missing service parameter"}`, http.StatusBadRequest)
		return
	}

	apiKey := r.URL.Query().Get("apikey")
	if apiKey == "" {
		http.Error(w, `{"error": "Missing apikey parameter"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Handling request for service: %s", service)

	// Create API client with provided credentials
	client, err := nsight.NewApiClientWithCredentials(apiKey, ps.server)
	if err != nil {
		log.Printf("Error creating API client: %v", err)
		http.Error(w, `{"error": "Failed to create API client"}`, http.StatusInternalServerError)
		return
	}

	// Route to appropriate handler based on service
	var result interface{}

	switch service {
	case "list_clients":
		result, err = client.FetchClients()
	
	case "list_sites":
		clientID, parseErr := strconv.Atoi(r.URL.Query().Get("clientid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid clientid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchSites(clientID)
	
	case "list_servers":
		siteID, parseErr := strconv.Atoi(r.URL.Query().Get("siteid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid siteid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchServers(siteID)
	
	case "list_workstations":
		siteID, parseErr := strconv.Atoi(r.URL.Query().Get("siteid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid siteid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchWorkstations(siteID)
	
	case "list_devices":
		siteID, parseErr := strconv.Atoi(r.URL.Query().Get("siteid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid siteid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchDevicesBySite(siteID)
	
	case "list_devices_at_client":
		clientID, parseErr := strconv.Atoi(r.URL.Query().Get("clientid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid clientid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchDevices(clientID)
	
	case "list_device_asset_details":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchDeviceAssetDetails(deviceID)
	
	case "list_failing_checks":
		result, err = client.FetchFailingChecks()
	
	case "list_checks":
		deviceIDStr := r.URL.Query().Get("deviceid")
		siteIDStr := r.URL.Query().Get("siteid")
		
		if deviceIDStr != "" {
			deviceID, parseErr := strconv.Atoi(deviceIDStr)
			if parseErr != nil {
				http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
				return
			}
			result, err = client.FetchChecks(deviceID)
		} else if siteIDStr != "" {
			siteID, parseErr := strconv.Atoi(siteIDStr)
			if parseErr != nil {
				http.Error(w, `{"error": "Invalid siteid parameter"}`, http.StatusBadRequest)
				return
			}
			result, err = client.FetchChecksBySite(siteID)
		} else {
			http.Error(w, `{"error": "Missing deviceid or siteid parameter"}`, http.StatusBadRequest)
			return
		}
	
	case "list_device_monitoring_details":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchDeviceMonitoringDetails(deviceID)
	
	case "list_agentless_assets":
		siteID, parseErr := strconv.Atoi(r.URL.Query().Get("siteid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid siteid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchAgentlessAssets(siteID)
	
	case "list_hardware":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchHardware(deviceID)
	
	case "list_software":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchSoftware(deviceID)
	
	case "list_license_groups":
		result, err = client.FetchLicenseGroups()
	
	case "list_patches":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchPatches(deviceID)
	
	case "list_antivirus_products":
		result, err = client.FetchAntivirusProducts()
	
	case "list_antivirus_definitions":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchAntivirusDefinitions(deviceID)
	
	case "list_quarantine":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		result, err = client.FetchQuarantineList(deviceID)
	
	case "list_performance_history":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		checkID, parseErr := strconv.Atoi(r.URL.Query().Get("checkid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid checkid parameter"}`, http.StatusBadRequest)
			return
		}
		startDate := r.URL.Query().Get("startdate")
		endDate := r.URL.Query().Get("enddate")
		result, err = client.FetchPerformanceHistory(deviceID, checkID, startDate, endDate)
	
	case "list_drive_space_history":
		deviceID, parseErr := strconv.Atoi(r.URL.Query().Get("deviceid"))
		if parseErr != nil {
			http.Error(w, `{"error": "Invalid deviceid parameter"}`, http.StatusBadRequest)
			return
		}
		startDate := r.URL.Query().Get("startdate")
		endDate := r.URL.Query().Get("enddate")
		result, err = client.FetchDriveSpaceHistory(deviceID, startDate, endDate)
	
	case "list_templates":
		result, err = client.FetchTemplates()
	
	default:
		http.Error(w, fmt.Sprintf(`{"error": "Unsupported service: %s"}`, service), http.StatusBadRequest)
		return
	}

	// Handle any error from the API call
	if err != nil {
		log.Printf("Error calling API service %s: %v", service, err)
		http.Error(w, fmt.Sprintf(`{"error": "API call failed: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Convert result to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshaling JSON for service %s: %v", service, err)
		http.Error(w, `{"error": "Failed to convert response to JSON"}`, http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Write(jsonData)
}

// healthCheck provides a simple health check endpoint
func (ps *ProxyServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "nsight-proxy"}`))
}

func main() {
	log.Println("Starting N-Sight JSON Proxy Server...")

	// Create proxy server instance
	proxy, err := NewProxyServer()
	if err != nil {
		log.Fatalf("Failed to initialize proxy server: %v", err)
	}

	// Set up routes
	http.HandleFunc("/api/", proxy.handleAPI)
	http.HandleFunc("/health", proxy.healthCheck)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"service": "N-Sight JSON Proxy", "version": "1.0", "endpoints": ["/api/", "/health"]}`))
	})

	// Start server on port 80
	log.Println("Server starting on port 80...")
	log.Println("API endpoint: http://localhost/api/?service=<service_name>&<parameters>")
	log.Println("Health check: http://localhost/health")
	
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}