package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"nsight-proxy/internal/nsight"
)

// formatUnixTimestamp converts a Unix timestamp string to "DD.MM.YYYY HH:MM:SS" format.
// Returns an empty string if the input is invalid or conversion fails.
func formatUnixTimestamp(timestampStr string) string {
	if timestampStr == "" {
		return ""
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Printf("Warning: Failed to parse timestamp '%s': %v", timestampStr, err)
		return "" // Or return timestampStr ? Or a placeholder like "Invalid Date"?
	}
	if timestamp == 0 { // Often used as a nil time value
		return ""
	}
	t := time.Unix(timestamp, 0)
	return t.Format("02.01.2006 15:04:05")
}

// --- Output Structures for Nested JSON ---

type ServerDetail struct {
	ID           int                  `json:"server_id"`
	Name         string               `json:"server_name"`
	Online       bool                 `json:"online"`
	OS           string               `json:"os,omitempty"`
	IP           string               `json:"ip,omitempty"`
	User         string               `json:"user,omitempty"`
	Manufacturer string               `json:"manufacturer,omitempty"`
	Model        string               `json:"model,omitempty"`
	DeviceSerial string               `json:"serial_number,omitempty"`
	LastBootTime string               `json:"last_boot_time,omitempty"`
	AssetInfo    *nsight.AssetDetails `json:"asset_details,omitempty"`
}

type WorkstationDetail struct {
	ID           int                  `json:"workstation_id"`
	Name         string               `json:"workstation_name"`
	Online       bool                 `json:"online"`
	OS           string               `json:"os,omitempty"`
	IP           string               `json:"ip,omitempty"`
	User         string               `json:"user,omitempty"`
	Manufacturer string               `json:"manufacturer,omitempty"`
	Model        string               `json:"model,omitempty"`
	DeviceSerial string               `json:"serial_number,omitempty"`
	LastBootTime string               `json:"last_boot_time,omitempty"`
	AssetInfo    *nsight.AssetDetails `json:"asset_details,omitempty"`
}

type SiteDetail struct {
	ID           int                 `json:"site_id"`
	Name         string              `json:"site_name"`
	Servers      []ServerDetail      `json:"servers,omitempty"`
	Workstations []WorkstationDetail `json:"workstations,omitempty"`
}

type ClientDetail struct {
	ID    int          `json:"client_id"`
	Name  string       `json:"client_name"`
	Sites []SiteDetail `json:"sites,omitempty"`
}

// --- CSV Reading Functions ---

// readCsvData reads a specified CSV file, skipping the header
func readCsvData(filename string) ([][]string, error) {
	path := filepath.Join("data", filename)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("cache file %s not found. Run fetchall without -cache first", path)
		}
		return nil, fmt.Errorf("failed to open cache file %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header row
	if err != nil {
		if err == io.EOF {
			return [][]string{}, nil // Empty file is valid
		}
		return nil, fmt.Errorf("failed to read header from %s: %w", path, err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from %s: %w", path, err)
	}
	return records, nil
}

// buildResultFromCache reconstructs the nested structure from CSV data
func buildResultFromCache() ([]ClientDetail, error) {
	log.Println("Building result from CSV cache...")

	// Read all CSV files
	clientRecords, err := readCsvData("clients.csv")
	if err != nil {
		return nil, err
	}
	siteRecords, err := readCsvData("sites.csv")
	if err != nil {
		return nil, err
	}
	serverRecords, err := readCsvData("servers.csv")
	if err != nil {
		return nil, err
	}
	workstationRecords, err := readCsvData("workstations.csv")
	if err != nil {
		return nil, err
	}

	// --- Process records into usable maps for efficient lookup ---

	// Map clientID -> Client Name
	clientMap := make(map[int]string)
	for _, rec := range clientRecords {
		if len(rec) < 2 {
			continue
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			continue
		}
		clientMap[id] = rec[1]
	}

	// Map clientID -> list of Sites (ID, Name)
	sitesByClient := make(map[int][]SiteDetail)
	for _, rec := range siteRecords {
		if len(rec) < 3 {
			continue
		}
		siteID, errS := strconv.Atoi(rec[0])
		clientID, errC := strconv.Atoi(rec[2])
		if errS != nil || errC != nil {
			continue
		}
		sitesByClient[clientID] = append(sitesByClient[clientID], SiteDetail{ID: siteID, Name: rec[1]})
	}

	// Map siteID -> list of Servers
	serversBySite := make(map[int][]ServerDetail)
	for _, rec := range serverRecords {
		// Check for the extended number of columns (now 12: ID, Name, OS, IP, Online, User, Manufacturer, Model, Serial, LastBootTime, SiteID, ClientID)
		if len(rec) < 12 { // Updated count
			log.Printf("Warning: Skipping malformed server record in cache: %v", rec)
			continue
		}
		serverID, errSv := strconv.Atoi(rec[0])
		onlineInt, errO := strconv.Atoi(rec[4])
		siteID, errSi := strconv.Atoi(rec[10]) // Site ID is now at index 10
		if errSv != nil || errO != nil || errSi != nil {
			log.Printf("Warning: Skipping server record with invalid numeric data: %v", rec)
			continue
		}
		serversBySite[siteID] = append(serversBySite[siteID], ServerDetail{
			ID:           serverID,
			Name:         rec[1],
			OS:           rec[2],
			IP:           rec[3],
			Online:       onlineInt == 1,
			User:         rec[5],
			Manufacturer: rec[6],
			Model:        rec[7],
			DeviceSerial: rec[8],
			LastBootTime: formatUnixTimestamp(rec[9]), // Format timestamp
		})
	}

	// Map siteID -> list of Workstations
	workstationsBySite := make(map[int][]WorkstationDetail)
	for _, rec := range workstationRecords {
		// Check for the extended number of columns (now 12)
		if len(rec) < 12 { // Updated count
			log.Printf("Warning: Skipping malformed workstation record in cache: %v", rec)
			continue
		}
		wsID, errW := strconv.Atoi(rec[0])
		onlineInt, errO := strconv.Atoi(rec[4])
		siteID, errS := strconv.Atoi(rec[10]) // Site ID is now at index 10
		if errW != nil || errO != nil || errS != nil {
			log.Printf("Warning: Skipping workstation record with invalid numeric data: %v", rec)
			continue
		}
		workstationsBySite[siteID] = append(workstationsBySite[siteID], WorkstationDetail{
			ID:           wsID,
			Name:         rec[1],
			OS:           rec[2],
			IP:           rec[3],
			Online:       onlineInt == 1,
			User:         rec[5],
			Manufacturer: rec[6],
			Model:        rec[7],
			DeviceSerial: rec[8],
			LastBootTime: formatUnixTimestamp(rec[9]), // Format timestamp
		})
	}

	// --- Build the final nested structure ---

	var finalResult []ClientDetail
	for clientID, clientName := range clientMap {
		clientDetail := ClientDetail{
			ID:    clientID,
			Name:  clientName,
			Sites: []SiteDetail{},
		}

		sites, ok := sitesByClient[clientID]
		if ok {
			for _, site := range sites {
				siteDetail := SiteDetail{
					ID:           site.ID,
					Name:         site.Name,
					Servers:      serversBySite[site.ID],      // Will be nil if no servers, which is fine for JSON omitempty
					Workstations: workstationsBySite[site.ID], // Will be nil if no workstations
				}
				clientDetail.Sites = append(clientDetail.Sites, siteDetail)
			}
		}

		finalResult = append(finalResult, clientDetail)
	}

	log.Println("Successfully built result from CSV cache.")
	return finalResult, nil
}

// --- CSV Writing Functions ---

// Helper to open or create a CSV file and return the writer
func openCsvWriter(path string, header []string) (*csv.Writer, *os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create directory %s: %w", filepath.Dir(path), err)
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file %s: %w", path, err)
	}
	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("failed to write header to %s: %w", path, err)
	}
	writer.Flush() // Ensure header is written immediately
	return writer, file, nil
}

func main() {
	// Define and parse flags
	cacheMode := flag.Bool("cache", false, "Read data from CSV cache instead of fetching from API")
	flag.Parse()

	// Determine output filename (non-flag argument)
	outputFilename := ""
	if flag.NArg() > 0 {
		outputFilename = flag.Arg(0)
	}

	var finalResult []ClientDetail
	var err error

	if *cacheMode {
		// --- Cache Mode ---
		finalResult, err = buildResultFromCache()
		if err != nil {
			log.Fatalf("Error building result from cache: %v", err)
		}
	} else {
		// --- API Fetch Mode ---
		log.Println("Starting fetchall process from API...")
		apiClient, err := nsight.NewApiClient()
		if err != nil {
			log.Fatalf("Failed to initialize API client: %v", err)
		}

		// Prepare CSV Writers (only needed in non-cache mode)
		csvHeaders := map[string][]string{
			"clients":      {"client_id", "name"},
			"sites":        {"site_id", "name", "client_id"},
			"servers":      {"server_id", "name", "os", "ip", "online", "user", "manufacturer", "model", "serial_number", "last_boot_time", "site_id", "client_id"},
			"workstations": {"workstation_id", "name", "os", "ip", "online", "user", "manufacturer", "model", "serial_number", "last_boot_time", "site_id", "client_id"},
		}
		writers := make(map[string]*csv.Writer)
		files := make(map[string]*os.File)
		defer func() {
			for name, writer := range writers {
				writer.Flush()
				if ferr := writer.Error(); ferr != nil {
					log.Printf("Error flushing CSV writer for %s: %v", name, ferr)
				}
			}
			for name, file := range files {
				if ferr := file.Close(); ferr != nil {
					log.Printf("Error closing file %s: %v", name, ferr)
				}
			}
		}()

		for name, header := range csvHeaders {
			path := filepath.Join("data", name+".csv")
			writer, file, ferr := openCsvWriter(path, header)
			if ferr != nil {
				log.Fatalf("Failed to setup CSV writer for %s: %v", name, ferr)
			}
			writers[name] = writer
			files[name] = file
			log.Printf("Opened %s for writing.", path)
		}

		// Fetch and Process Data from API
		log.Println("Fetching clients from API...")
		clients, err := apiClient.FetchClients()
		if err != nil {
			log.Fatalf("Failed to fetch clients: %v", err)
		}
		log.Printf("Fetched %d clients.", len(clients))

		// finalResult is built during the fetch loop
		finalResult = []ClientDetail{}

		for _, client := range clients {
			// Write client to CSV
			if err := writers["clients"].Write([]string{strconv.Itoa(client.ClientID), client.Name}); err != nil {
				log.Printf("Warning: Failed to write client %d to CSV: %v", client.ClientID, err)
			}

			log.Printf("Fetching sites for client %d (%s)...", client.ClientID, client.Name)
			sites, err := apiClient.FetchSites(client.ClientID)
			if err != nil {
				log.Printf("Warning: Failed to fetch sites for client %d: %v. Skipping client.", client.ClientID, err)
				continue // Skip this client if sites can't be fetched
			}
			log.Printf("Fetched %d sites for client %d.", len(sites), client.ClientID)

			clientDetail := ClientDetail{
				ID:    client.ClientID,
				Name:  client.Name,
				Sites: []SiteDetail{},
			}

			for _, site := range sites {
				// Write site to CSV
				if err := writers["sites"].Write([]string{
					strconv.Itoa(site.SiteID),
					site.Name,
					strconv.Itoa(client.ClientID),
				}); err != nil {
					log.Printf("Warning: Failed to write site %d to CSV: %v", site.SiteID, err)
				}

				log.Printf("Fetching servers for site %d (%s)...", site.SiteID, site.Name)
				servers, err := apiClient.FetchServers(site.SiteID)
				if err != nil {
					log.Printf("Warning: Failed to fetch servers for site %d: %v", site.SiteID, err)
				}
				log.Printf("Fetched %d servers for site %d.", len(servers), site.SiteID)

				log.Printf("Fetching workstations for site %d (%s)...", site.SiteID, site.Name)
				workstations, err := apiClient.FetchWorkstations(site.SiteID)
				if err != nil {
					log.Printf("Warning: Failed to fetch workstations for site %d: %v", site.SiteID, err)
				}
				log.Printf("Fetched %d workstations for site %d.", len(workstations), site.SiteID)

				siteDetail := SiteDetail{
					ID:           site.SiteID,
					Name:         site.Name,
					Servers:      []ServerDetail{},
					Workstations: []WorkstationDetail{},
				}

				// Process and write servers
				for _, server := range servers {
					// Format the timestamp
					formattedBootTime := formatUnixTimestamp(server.LastBootTime)

					// Fetch asset details
					assetDetails, err := apiClient.FetchDeviceAssetDetails(server.ServerID)
					if err != nil {
						log.Printf("Warning: Failed to fetch asset details for server %d: %v", server.ServerID, err)
						// assetDetails will be nil, so AssetInfo will be omitted in JSON
					}

					// Write server to CSV (extended data) - CSV does not include asset details
					if err := writers["servers"].Write([]string{
						strconv.Itoa(server.ServerID),
						server.Name,
						server.OS,
						server.IP,
						strconv.Itoa(server.Online),
						server.User,
						server.Manufacturer,
						server.Model,
						server.DeviceSerial,
						formattedBootTime, // Use formatted time
						strconv.Itoa(site.SiteID),
						strconv.Itoa(client.ClientID),
					}); err != nil {
						log.Printf("Warning: Failed to write server %d to CSV: %v", server.ServerID, err)
					}
					// Append server detail to site (extended data)
					siteDetail.Servers = append(siteDetail.Servers, ServerDetail{
						ID:           server.ServerID,
						Name:         server.Name,
						Online:       server.Online == 1,
						OS:           server.OS,
						IP:           server.IP,
						User:         server.User,
						Manufacturer: server.Manufacturer,
						Model:        server.Model,
						DeviceSerial: server.DeviceSerial,
						LastBootTime: formattedBootTime, // Use formatted time
						AssetInfo:    assetDetails,      // Assign fetched asset details
					})
				}

				// Process and write workstations
				for _, ws := range workstations {
					// Format the timestamp
					formattedBootTime := formatUnixTimestamp(ws.LastBootTime)

					// Fetch asset details
					assetDetails, err := apiClient.FetchDeviceAssetDetails(ws.WorkstationID)
					if err != nil {
						log.Printf("Warning: Failed to fetch asset details for workstation %d: %v", ws.WorkstationID, err)
						// assetDetails will be nil, so AssetInfo will be omitted in JSON
					}

					// Write workstation to CSV (extended data) - CSV does not include asset details
					if err := writers["workstations"].Write([]string{
						strconv.Itoa(ws.WorkstationID),
						ws.Name,
						ws.OS,
						ws.IP,
						strconv.Itoa(ws.Online),
						ws.User,
						ws.Manufacturer,
						ws.Model,
						ws.DeviceSerial,
						formattedBootTime, // Use formatted time
						strconv.Itoa(site.SiteID),
						strconv.Itoa(client.ClientID),
					}); err != nil {
						log.Printf("Warning: Failed to write workstation %d to CSV: %v", ws.WorkstationID, err)
					}
					// Append workstation detail to site (extended data)
					siteDetail.Workstations = append(siteDetail.Workstations, WorkstationDetail{
						ID:           ws.WorkstationID,
						Name:         ws.Name,
						Online:       ws.Online == 1,
						OS:           ws.OS,
						IP:           ws.IP,
						User:         ws.User,
						Manufacturer: ws.Manufacturer,
						Model:        ws.Model,
						DeviceSerial: ws.DeviceSerial,
						LastBootTime: formattedBootTime, // Use formatted time
						AssetInfo:    assetDetails,      // Assign fetched asset details
					})
				}

				clientDetail.Sites = append(clientDetail.Sites, siteDetail)
			} // end site loop

			finalResult = append(finalResult, clientDetail)
			writers["clients"].Flush()      // Flush periodically for clients
			writers["sites"].Flush()        // Flush periodically for sites
			writers["servers"].Flush()      // Flush periodically for servers
			writers["workstations"].Flush() // Flush periodically for workstations
		} // end client loop

		log.Println("Finished fetching data from API.")
	} // End of else block (API fetch mode)

	// --- Output Final JSON ---
	log.Println("Marshalling final JSON...")
	finalJsonData, err := json.MarshalIndent(finalResult, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal final result to JSON: %v", err)
	}

	// Output to file or stdout based on outputFilename
	if outputFilename != "" {
		log.Printf("Writing JSON output to file: %s", outputFilename)
		if err := os.WriteFile(outputFilename, finalJsonData, 0644); err != nil {
			log.Fatalf("Failed to write JSON to file %s: %v", outputFilename, err)
		}
	} else {
		fmt.Println(string(finalJsonData))
	}

	log.Println("Fetchall process completed successfully.")
}
