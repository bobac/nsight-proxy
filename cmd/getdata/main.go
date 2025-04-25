package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"nsight-proxy/internal/nsight"
)

// -- Output Structures (Specific to getdata tool) --

// ClientOutput defines the structure for JSON output for clients
type ClientOutput struct {
	ClientID int    `json:"client_id"`
	Name     string `json:"client_name"`
}

// SiteOutput defines the structure for JSON output for sites
type SiteOutput struct {
	SiteID int    `json:"site_id"`
	Name   string `json:"site_name"`
}

// ServerOutput defines the structure for JSON output for servers
type ServerOutput struct {
	ServerID     int    `json:"server_id"`
	Name         string `json:"server_name"`
	Online       bool   `json:"online"`
	OS           string `json:"os,omitempty"`
	IP           string `json:"ip,omitempty"`
	User         string `json:"user,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	DeviceSerial string `json:"serial_number,omitempty"`
}

// WorkstationOutput defines the structure for JSON output for workstations
type WorkstationOutput struct {
	WorkstationID int    `json:"workstation_id"`
	Name          string `json:"workstation_name"`
	Online        bool   `json:"online"`
	OS            string `json:"os,omitempty"`
	IP            string `json:"ip,omitempty"`
	User          string `json:"user,omitempty"`
	Manufacturer  string `json:"manufacturer,omitempty"`
	Model         string `json:"model,omitempty"`
	DeviceSerial  string `json:"serial_number,omitempty"`
}

// -- Helper Functions (Now using nsight.ApiClient) --

// getClientIDByName fetches clients and finds the ID for a given client name
func getClientIDByName(apiClient *nsight.ApiClient, clientName string) (int, error) {
	fmt.Println("Fetching client list to find ID for name:", clientName)
	clients, err := apiClient.FetchClients()
	if err != nil {
		return 0, fmt.Errorf("could not fetch client list to find ID: %w", err)
	}

	for _, client := range clients {
		if client.Name == clientName {
			fmt.Printf("Found Client ID %d for name '%s'\n", client.ClientID, clientName)
			return client.ClientID, nil
		}
	}

	return 0, fmt.Errorf("client with name '%s' not found", clientName)
}

// getSiteIDByName fetches clients and sites to find the ID for a given site name
func getSiteIDByName(apiClient *nsight.ApiClient, siteName string) (int, error) {
	fmt.Println("Fetching clients to search for site name:", siteName)
	clients, err := apiClient.FetchClients()
	if err != nil {
		return 0, fmt.Errorf("could not fetch client list to find site ID: %w", err)
	}

	fmt.Println("Iterating through clients and their sites...")
	for _, client := range clients {
		fmt.Printf("Checking sites for client %d (%s)...\n", client.ClientID, client.Name)
		sites, err := apiClient.FetchSites(client.ClientID)
		if err != nil {
			log.Printf("Warning: Failed to fetch sites for client %d (%s): %v", client.ClientID, client.Name, err)
			continue
		}
		for _, site := range sites {
			if site.Name == siteName {
				fmt.Printf("Found Site ID %d for name '%s' under client %d (%s)\n", site.SiteID, siteName, client.ClientID, client.Name)
				return site.SiteID, nil
			}
		}
	}

	return 0, fmt.Errorf("site with name '%s' not found across any client", siteName)
}

// printClientsAsJSON fetches clients and prints them in the desired JSON format
func printClientsAsJSON(apiClient *nsight.ApiClient) {
	fmt.Println("Running getdata tool for list_clients...")
	clients, err := apiClient.FetchClients()
	if err != nil {
		log.Fatalf("Error fetching clients: %v", err)
	}

	var output []ClientOutput
	for _, client := range clients {
		output = append(output, ClientOutput{
			ClientID: client.ClientID,
			Name:     client.Name,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling client data to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

// printSitesAsJSON fetches sites for a given client ID and prints them as JSON
func printSitesAsJSON(apiClient *nsight.ApiClient, clientID int) {
	fmt.Printf("Running getdata tool for list_sites (Client ID: %d)...\n", clientID)
	sites, err := apiClient.FetchSites(clientID)
	if err != nil {
		log.Fatalf("Error fetching sites: %v", err)
	}

	var output []SiteOutput
	for _, site := range sites {
		output = append(output, SiteOutput{
			SiteID: site.SiteID,
			Name:   site.Name,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling site data to JSON: %v", err)
	}
	if len(sites) == 0 {
		fmt.Println("[]")
	} else {
		fmt.Println(string(jsonData))
	}
}

// printServersAsJSON fetches servers for a given site ID and prints them as JSON
func printServersAsJSON(apiClient *nsight.ApiClient, siteID int) {
	fmt.Printf("Running getdata tool for list_servers (Site ID: %d)...\n", siteID)
	servers, err := apiClient.FetchServers(siteID)
	if err != nil {
		log.Fatalf("Error fetching servers: %v", err)
	}

	var output []ServerOutput
	for _, serverItem := range servers {
		output = append(output, ServerOutput{
			ServerID:     serverItem.ServerID,
			Name:         serverItem.Name,
			Online:       serverItem.Online == 1,
			OS:           serverItem.OS,
			IP:           serverItem.IP,
			User:         serverItem.User,
			Manufacturer: serverItem.Manufacturer,
			Model:        serverItem.Model,
			DeviceSerial: serverItem.DeviceSerial,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling server data to JSON: %v", err)
	}
	if len(servers) == 0 {
		fmt.Println("[]")
	} else {
		fmt.Println(string(jsonData))
	}
}

// printWorkstationsAsJSON fetches workstations for a given site ID and prints them as JSON
func printWorkstationsAsJSON(apiClient *nsight.ApiClient, siteID int) {
	fmt.Printf("Running getdata tool for list_workstations (Site ID: %d)...\n", siteID)
	workstations, err := apiClient.FetchWorkstations(siteID)
	if err != nil {
		log.Fatalf("Error fetching workstations: %v", err)
	}

	var output []WorkstationOutput
	for _, ws := range workstations {
		output = append(output, WorkstationOutput{
			WorkstationID: ws.WorkstationID,
			Name:          ws.Name,
			Online:        ws.Online == 1,
			OS:            ws.OS,
			IP:            ws.IP,
			User:          ws.User,
			Manufacturer:  ws.Manufacturer,
			Model:         ws.Model,
			DeviceSerial:  ws.DeviceSerial,
		})
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling workstation data to JSON: %v", err)
	}
	if len(workstations) == 0 {
		fmt.Println("[]")
	} else {
		fmt.Println(string(jsonData))
	}
}

func main() {
	apiClient, err := nsight.NewApiClient()
	if err != nil {
		log.Fatalf("Failed to initialize API client: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/getdata/main.go <service_name> [parameters...]")
		fmt.Println("Available services:")
		fmt.Println("  list_clients")
		fmt.Println("  list_sites <client_id | \"client_name\">")
		fmt.Println("  list_servers <site_id | \"site_name\">")
		fmt.Println("  list_workstations <site_id | \"site_name\">")
		os.Exit(1)
	}

	serviceName := os.Args[1]

	switch serviceName {
	case "list_clients":
		if len(os.Args) != 2 {
			log.Fatalf("Usage: go run cmd/getdata/main.go list_clients")
		}
		printClientsAsJSON(apiClient)
	case "list_sites":
		if len(os.Args) != 3 {
			log.Fatalf("Usage: go run cmd/getdata/main.go list_sites <client_id | \"client_name\">")
		}
		clientIdentifier := os.Args[2]
		var clientID int

		clientID, err = strconv.Atoi(clientIdentifier)
		if err != nil {
			fmt.Printf("'%s' is not a valid integer ID, attempting to find client by name.\n", clientIdentifier)
			clientID, err = getClientIDByName(apiClient, clientIdentifier)
			if err != nil {
				log.Fatalf("Error finding client ID: %v", err)
			}
		}
		printSitesAsJSON(apiClient, clientID)
	case "list_servers":
		if len(os.Args) != 3 {
			log.Fatalf("Usage: go run cmd/getdata/main.go list_servers <site_id | \"site_name\">")
		}
		siteIdentifier := os.Args[2]
		var siteID int

		siteID, err = strconv.Atoi(siteIdentifier)
		if err != nil {
			fmt.Printf("'%s' is not a valid integer ID, attempting to find site by name across all clients...\n", siteIdentifier)
			siteID, err = getSiteIDByName(apiClient, siteIdentifier)
			if err != nil {
				log.Fatalf("Error finding site ID: %v", err)
			}
		}
		printServersAsJSON(apiClient, siteID)
	case "list_workstations":
		if len(os.Args) != 3 {
			log.Fatalf("Usage: go run cmd/getdata/main.go list_workstations <site_id | \"site_name\">")
		}
		siteIdentifier := os.Args[2]
		var siteID int

		siteID, err = strconv.Atoi(siteIdentifier)
		if err != nil {
			fmt.Printf("'%s' is not a valid integer ID, attempting to find site by name across all clients...\n", siteIdentifier)
			siteID, err = getSiteIDByName(apiClient, siteIdentifier)
			if err != nil {
				log.Fatalf("Error finding site ID: %v", err)
			}
		}
		printWorkstationsAsJSON(apiClient, siteID)
	default:
		log.Fatalf("Error: Unknown service '%s'", serviceName)
	}
}
