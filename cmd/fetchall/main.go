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

	// --- Read Base Data CSV files ---
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

	// --- Read Asset Data CSV files (handle missing files gracefully) ---
	assetSummaryRecords, errAssetSummary := readCsvData("asset_summary.csv")
	if errAssetSummary != nil && !os.IsNotExist(errAssetSummary) {
		return nil, fmt.Errorf("failed to read asset_summary.csv: %w", errAssetSummary)
	} else if os.IsNotExist(errAssetSummary) {
		log.Println("Warning: asset_summary.csv not found in cache. Asset details will be missing.")
		assetSummaryRecords = [][]string{} // Empty slice
	}

	hardwareRecords, errHardware := readCsvData("hardware_assets.csv")
	if errHardware != nil && !os.IsNotExist(errHardware) {
		return nil, fmt.Errorf("failed to read hardware_assets.csv: %w", errHardware)
	} else if os.IsNotExist(errHardware) {
		log.Println("Warning: hardware_assets.csv not found in cache. Hardware details will be missing.")
		hardwareRecords = [][]string{} // Empty slice
	}

	softwareRecords, errSoftware := readCsvData("software_assets.csv")
	if errSoftware != nil && !os.IsNotExist(errSoftware) {
		return nil, fmt.Errorf("failed to read software_assets.csv: %w", errSoftware)
	} else if os.IsNotExist(errSoftware) {
		log.Println("Warning: software_assets.csv not found in cache. Software details will be missing.")
		softwareRecords = [][]string{} // Empty slice
	}

	// --- Process records into usable maps for efficient lookup ---

	// Map clientID -> Client Name
	clientMap := make(map[int]string)
	for _, rec := range clientRecords {
		if len(rec) < 2 {
			log.Printf("Warning: Skipping malformed client record in cache: %v", rec)
			continue
		}
		clientID, err := strconv.Atoi(rec[0])
		if err != nil {
			log.Printf("Warning: Skipping client record with invalid ID: %v", rec)
			continue
		}
		clientMap[clientID] = rec[1]
	}

	// Map clientID -> list of Sites (ID, Name)
	sitesByClient := make(map[int][]SiteDetail)
	for _, rec := range siteRecords {
		if len(rec) < 3 {
			log.Printf("Warning: Skipping malformed site record in cache: %v", rec)
			continue
		}
		siteID, errS := strconv.Atoi(rec[0])
		clientID, errC := strconv.Atoi(rec[2])
		if errS != nil || errC != nil {
			log.Printf("Warning: Skipping site record with invalid numeric data: %v", rec)
			continue
		}
		sitesByClient[clientID] = append(sitesByClient[clientID], SiteDetail{ID: siteID, Name: rec[1]})
	}

	// --- Process Asset Data into Maps ---

	// Map deviceID -> Asset Summary Data (using AssetDetails struct for convenience)
	assetSummaryMap := make(map[int]nsight.AssetDetails)
	for _, rec := range assetSummaryRecords {
		if len(rec) < 37 { // Expected number of columns in asset_summary.csv
			log.Printf("Warning: Skipping malformed asset summary record in cache: %v", rec)
			continue
		}
		deviceID, err := strconv.Atoi(rec[0])
		if err != nil {
			log.Printf("Warning: Skipping asset summary record with invalid device ID: %v", rec)
			continue
		}
		ram, _ := strconv.ParseInt(rec[15], 10, 64) // Ignore error, default to 0
		assetSummaryMap[deviceID] = nsight.AssetDetails{
			Client:       rec[1],
			ChassisType:  rec[2],
			IP:           rec[3],
			MAC1:         rec[4],
			MAC2:         rec[5],
			MAC3:         rec[6],
			User:         rec[7],
			Manufacturer: rec[8],
			Model:        rec[9],
			OS:           rec[10],
			SerialNumber: rec[11],
			ProductKey:   rec[12],
			Role:         rec[13],
			ServicePack:  rec[14],
			RAM:          ram,
			ScanTime:     rec[16],
			Custom1:      nsight.CustomField{Name: rec[17], Value: rec[18]},
			Custom2:      nsight.CustomField{Name: rec[19], Value: rec[20]},
			Custom3:      nsight.CustomField{Name: rec[21], Value: rec[22]},
			Custom4:      nsight.CustomField{Name: rec[23], Value: rec[24]},
			Custom5:      nsight.CustomField{Name: rec[25], Value: rec[26]},
			Custom6:      nsight.CustomField{Name: rec[27], Value: rec[28]},
			Custom7:      nsight.CustomField{Name: rec[29], Value: rec[30]},
			Custom8:      nsight.CustomField{Name: rec[31], Value: rec[32]},
			Custom9:      nsight.CustomField{Name: rec[33], Value: rec[34]},
			Custom10:     nsight.CustomField{Name: rec[35], Value: rec[36]},
			// Hardware and Software lists will be populated later
		}
	}

	// Map deviceID -> []HardwareItem
	hardwareMap := make(map[int][]nsight.HardwareItem)
	for _, rec := range hardwareRecords {
		if len(rec) < 9 { // Expected columns in hardware_assets.csv
			log.Printf("Warning: Skipping malformed hardware asset record in cache: %v", rec)
			continue
		}
		deviceID, err := strconv.Atoi(rec[0])
		if err != nil {
			log.Printf("Warning: Skipping hardware asset record with invalid device ID: %v", rec)
			continue
		}
		hwID, _ := strconv.Atoi(rec[1])
		hwType, _ := strconv.Atoi(rec[3])
		hwDeleted, _ := strconv.Atoi(rec[7])
		hwModified, _ := strconv.Atoi(rec[8])
		hardwareMap[deviceID] = append(hardwareMap[deviceID], nsight.HardwareItem{
			HardwareID:   hwID,
			Name:         rec[2],
			Type:         hwType,
			Manufacturer: rec[4],
			Details:      rec[5],
			Status:       rec[6],
			Deleted:      hwDeleted,
			Modified:     hwModified,
		})
	}

	// Map deviceID -> []SoftwareItem
	softwareMap := make(map[int][]nsight.SoftwareItem)
	for _, rec := range softwareRecords {
		if len(rec) < 8 { // Expected columns in software_assets.csv
			log.Printf("Warning: Skipping malformed software asset record in cache: %v", rec)
			continue
		}
		deviceID, err := strconv.Atoi(rec[0])
		if err != nil {
			log.Printf("Warning: Skipping software asset record with invalid device ID: %v", rec)
			continue
		}
		swID, _ := strconv.Atoi(rec[1])
		swDeleted, _ := strconv.Atoi(rec[6])
		swModified, _ := strconv.Atoi(rec[7])
		softwareMap[deviceID] = append(softwareMap[deviceID], nsight.SoftwareItem{
			SoftwareID:  swID,
			Name:        rec[2],
			Version:     rec[3],
			InstallDate: rec[4],
			Type:        rec[5],
			Deleted:     swDeleted,
			Modified:    swModified,
		})
	}

	// --- Process Base Data (Servers and Workstations) and combine with Asset Data ---

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

		// Look up and assign asset details from maps
		var assetInfoPtr *nsight.AssetDetails
		if summary, ok := assetSummaryMap[serverID]; ok {
			summary.Hardware = hardwareMap[serverID] // Assign hardware list
			summary.Software = softwareMap[serverID] // Assign software list
			assetInfoPtr = &summary                  // Assign pointer to the combined struct
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
			AssetInfo:    assetInfoPtr,                // Assign asset details from cache
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

		// Look up and assign asset details from maps
		var assetInfoPtr *nsight.AssetDetails
		if summary, ok := assetSummaryMap[wsID]; ok {
			summary.Hardware = hardwareMap[wsID] // Assign hardware list
			summary.Software = softwareMap[wsID] // Assign software list
			assetInfoPtr = &summary              // Assign pointer to the combined struct
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
			AssetInfo:    assetInfoPtr,                // Assign asset details from cache
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
				// Ensure Servers and Workstations are initialized to empty slices if nil
				servers := serversBySite[site.ID]
				if servers == nil {
					servers = []ServerDetail{}
				}
				workstations := workstationsBySite[site.ID]
				if workstations == nil {
					workstations = []WorkstationDetail{}
				}
				siteDetail := SiteDetail{
					ID:           site.ID,
					Name:         site.Name,
					Servers:      servers,
					Workstations: workstations,
				}
				clientDetail.Sites = append(clientDetail.Sites, siteDetail)
			}
		}
		// Ensure Sites slice is not nil if it remained empty
		if clientDetail.Sites == nil {
			clientDetail.Sites = []SiteDetail{}
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
			"clients":         {"client_id", "name"},
			"sites":           {"site_id", "name", "client_id"},
			"servers":         {"server_id", "name", "os", "ip", "online", "user", "manufacturer", "model", "serial_number", "last_boot_time", "site_id", "client_id"},
			"workstations":    {"workstation_id", "name", "os", "ip", "online", "user", "manufacturer", "model", "serial_number", "last_boot_time", "site_id", "client_id"},
			"asset_summary":   {"device_id", "client_name", "chassistype", "ip_asset", "mac1", "mac2", "mac3", "user_asset", "manufacturer_asset", "model_asset", "os_asset", "serialnumber_asset", "productkey", "role", "servicepack", "ram", "scantime", "custom1_name", "custom1_value", "custom2_name", "custom2_value", "custom3_name", "custom3_value", "custom4_name", "custom4_value", "custom5_name", "custom5_value", "custom6_name", "custom6_value", "custom7_name", "custom7_value", "custom8_name", "custom8_value", "custom9_name", "custom9_value", "custom10_name", "custom10_value"},
			"hardware_assets": {"device_id", "hardware_id", "name", "type", "manufacturer", "details", "status", "deleted", "modified"},
			"software_assets": {"device_id", "software_id", "name", "version", "install_date", "type", "deleted", "modified"},
		}
		writers := make(map[string]*csv.Writer)
		files := make(map[string]*os.File)
		defer func() {
			log.Println("Flushing and closing CSV files...")
			for name, writer := range writers {
				writer.Flush()
				if ferr := writer.Error(); ferr != nil {
					log.Printf("Error flushing CSV writer for %s: %v", name, ferr)
				}
			}
			for name, file := range files {
				if ferr := file.Close(); ferr != nil {
					log.Printf("Error closing file %s.csv: %v", name, ferr)
				}
			}
			log.Println("Finished flushing and closing CSV files.")
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

					// Write server to CSV (primary data)
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

					// Write asset details to separate CSVs if fetched successfully
					if assetDetails != nil {
						deviceIDStr := strconv.Itoa(server.ServerID)
						// Write asset summary
						if err := writers["asset_summary"].Write([]string{
							deviceIDStr,
							assetDetails.Client, assetDetails.ChassisType, assetDetails.IP, assetDetails.MAC1, assetDetails.MAC2, assetDetails.MAC3,
							assetDetails.User, assetDetails.Manufacturer, assetDetails.Model, assetDetails.OS, assetDetails.SerialNumber,
							assetDetails.ProductKey, assetDetails.Role, assetDetails.ServicePack, strconv.FormatInt(assetDetails.RAM, 10), assetDetails.ScanTime,
							assetDetails.Custom1.Name, assetDetails.Custom1.Value, assetDetails.Custom2.Name, assetDetails.Custom2.Value,
							assetDetails.Custom3.Name, assetDetails.Custom3.Value, assetDetails.Custom4.Name, assetDetails.Custom4.Value,
							assetDetails.Custom5.Name, assetDetails.Custom5.Value, assetDetails.Custom6.Name, assetDetails.Custom6.Value,
							assetDetails.Custom7.Name, assetDetails.Custom7.Value, assetDetails.Custom8.Name, assetDetails.Custom8.Value,
							assetDetails.Custom9.Name, assetDetails.Custom9.Value, assetDetails.Custom10.Name, assetDetails.Custom10.Value,
						}); err != nil {
							log.Printf("Warning: Failed to write asset summary for device %s to CSV: %v", deviceIDStr, err)
						}
						// Write hardware items
						for _, item := range assetDetails.Hardware {
							if err := writers["hardware_assets"].Write([]string{
								deviceIDStr, strconv.Itoa(item.HardwareID), item.Name, strconv.Itoa(item.Type), item.Manufacturer, item.Details, item.Status, strconv.Itoa(item.Deleted), strconv.Itoa(item.Modified),
							}); err != nil {
								log.Printf("Warning: Failed to write hardware asset %d for device %s to CSV: %v", item.HardwareID, deviceIDStr, err)
							}
						}
						// Write software items
						for _, item := range assetDetails.Software {
							if err := writers["software_assets"].Write([]string{
								deviceIDStr, strconv.Itoa(item.SoftwareID), item.Name, item.Version, item.InstallDate, item.Type, strconv.Itoa(item.Deleted), strconv.Itoa(item.Modified),
							}); err != nil {
								log.Printf("Warning: Failed to write software asset %d for device %s to CSV: %v", item.SoftwareID, deviceIDStr, err)
							}
						}
					}

					// Append server detail to site (including asset info pointer)
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

					// Write workstation to CSV (primary data)
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

					// Write asset details to separate CSVs if fetched successfully
					if assetDetails != nil {
						deviceIDStr := strconv.Itoa(ws.WorkstationID)
						// Write asset summary
						if err := writers["asset_summary"].Write([]string{
							deviceIDStr,
							assetDetails.Client, assetDetails.ChassisType, assetDetails.IP, assetDetails.MAC1, assetDetails.MAC2, assetDetails.MAC3,
							assetDetails.User, assetDetails.Manufacturer, assetDetails.Model, assetDetails.OS, assetDetails.SerialNumber,
							assetDetails.ProductKey, assetDetails.Role, assetDetails.ServicePack, strconv.FormatInt(assetDetails.RAM, 10), assetDetails.ScanTime,
							assetDetails.Custom1.Name, assetDetails.Custom1.Value, assetDetails.Custom2.Name, assetDetails.Custom2.Value,
							assetDetails.Custom3.Name, assetDetails.Custom3.Value, assetDetails.Custom4.Name, assetDetails.Custom4.Value,
							assetDetails.Custom5.Name, assetDetails.Custom5.Value, assetDetails.Custom6.Name, assetDetails.Custom6.Value,
							assetDetails.Custom7.Name, assetDetails.Custom7.Value, assetDetails.Custom8.Name, assetDetails.Custom8.Value,
							assetDetails.Custom9.Name, assetDetails.Custom9.Value, assetDetails.Custom10.Name, assetDetails.Custom10.Value,
						}); err != nil {
							log.Printf("Warning: Failed to write asset summary for device %s to CSV: %v", deviceIDStr, err)
						}
						// Write hardware items
						for _, item := range assetDetails.Hardware {
							if err := writers["hardware_assets"].Write([]string{
								deviceIDStr, strconv.Itoa(item.HardwareID), item.Name, strconv.Itoa(item.Type), item.Manufacturer, item.Details, item.Status, strconv.Itoa(item.Deleted), strconv.Itoa(item.Modified),
							}); err != nil {
								log.Printf("Warning: Failed to write hardware asset %d for device %s to CSV: %v", item.HardwareID, deviceIDStr, err)
							}
						}
						// Write software items
						for _, item := range assetDetails.Software {
							if err := writers["software_assets"].Write([]string{
								deviceIDStr, strconv.Itoa(item.SoftwareID), item.Name, item.Version, item.InstallDate, item.Type, strconv.Itoa(item.Deleted), strconv.Itoa(item.Modified),
							}); err != nil {
								log.Printf("Warning: Failed to write software asset %d for device %s to CSV: %v", item.SoftwareID, deviceIDStr, err)
							}
						}
					}

					// Append workstation detail to site (including asset info pointer)
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
