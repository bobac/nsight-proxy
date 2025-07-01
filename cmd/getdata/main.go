package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"nsight-proxy/internal/nsight"
)



// -- Main service function --
func main() {
	apiClient, err := nsight.NewApiClient()
	if err != nil {
		log.Fatalf("Failed to initialize API client: %v", err)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	serviceName := os.Args[1]
	args := os.Args[2:]

	switch serviceName {
	// -- Basic Entity Listing --
	case "list_clients":
		handleListClients(apiClient, args)
	case "list_sites":
		handleListSites(apiClient, args)
	case "list_servers":
		handleListServers(apiClient, args)
	case "list_workstations":
		handleListWorkstations(apiClient, args)
	case "list_devices":
		handleListDevices(apiClient, args)
	case "list_devices_at_client":
		handleListDevicesAtClient(apiClient, args)
	case "list_agentless_assets":
		handleListAgentlessAssets(apiClient, args)

	// -- Check and Monitoring --
	case "list_failing_checks":
		handleListFailingChecks(apiClient, args)
	case "list_checks":
		handleListChecks(apiClient, args)
	case "list_device_monitoring_details":
		handleListDeviceMonitoringDetails(apiClient, args)
	case "list_check_configuration", "list_check_configuration_windows", "list_check_configuration_mac", "list_check_configuration_linux":
		handleListCheckConfiguration(apiClient, serviceName, args)
	case "list_outages":
		handleListOutages(apiClient, args)
	case "clear_check":
		handleClearCheck(apiClient, args)
	case "add_check_note":
		handleAddCheckNote(apiClient, args)

	// -- Asset Tracking --
	case "list_hardware":
		handleListHardware(apiClient, args)
	case "list_software":
		handleListSoftware(apiClient, args)
	case "list_device_asset_details":
		handleListDeviceAssetDetails(apiClient, args)
	case "list_license_groups":
		handleListLicenseGroups(apiClient, args)

	// -- Patch Management --
	case "list_patches":
		handleListPatches(apiClient, args)
	case "approve_patch":
		handleApprovePatches(apiClient, args)
	case "ignore_patch":
		handleIgnorePatches(apiClient, args)

	// -- Antivirus --
	case "list_antivirus_products":
		handleListAntivirusProducts(apiClient, args)
	case "list_antivirus_definitions":
		handleListAntivirusDefinitions(apiClient, args)
	case "list_quarantine":
		handleListQuarantine(apiClient, args)
	case "start_scan":
		handleStartAntivirusScan(apiClient, args)

	// -- Performance and History --
	case "list_performance_history":
		handleListPerformanceHistory(apiClient, args)
	case "list_drive_space_history":
		handleListDriveSpaceHistory(apiClient, args)

	// -- Templates --
	case "list_templates":
		handleListTemplates(apiClient, args)

	// -- Backup & Recovery --
	case "list_backup_sessions":
		handleListBackupSessions(apiClient, args)

	// -- Settings --
	case "list_wall_chart_settings":
		handleListWallChartSettings(apiClient, args)
	case "list_general_settings":
		handleListGeneralSettings(apiClient, args)

	// -- Tasks and Users --
	case "list_active_directory_users":
		handleListActiveDirectoryUsers(apiClient, args)
	case "run_task_now":
		handleRunTaskNow(apiClient, args)

	// -- Site Management --
	case "add_client":
		handleAddClient(apiClient, args)
	case "add_site":
		handleAddSite(apiClient, args)
	case "get_site_installation_package":
		handleGetSiteInstallationPackage(apiClient, args)

	default:
		log.Fatalf("Error: Unknown service '%s'", serviceName)
	}
}

// -- Utility Functions --

func printUsage() {
	fmt.Println("Usage: go run cmd/getdata/main.go <service_name> [parameters...]")
	fmt.Println()
	fmt.Println("Basic Entity Listing:")
	fmt.Println("  list_clients")
	fmt.Println("  list_sites <client_id | \"client_name\">")
	fmt.Println("  list_servers <site_id | \"site_name\">")
	fmt.Println("  list_workstations <site_id | \"site_name\">")
	fmt.Println("  list_devices <site_id>")
	fmt.Println("  list_devices_at_client <client_id>")
	fmt.Println("  list_agentless_assets <site_id>")
	fmt.Println()
	fmt.Println("Check and Monitoring:")
	fmt.Println("  list_failing_checks")
	fmt.Println("  list_checks <device_id | site_id>")
	fmt.Println("  list_device_monitoring_details <device_id>")
	fmt.Println("  list_check_configuration <device_id> [os]")
	fmt.Println("  list_check_configuration_windows <device_id>")
	fmt.Println("  list_check_configuration_mac <device_id>")
	fmt.Println("  list_check_configuration_linux <device_id>")
	fmt.Println("  list_outages <site_id> <start_date> <end_date>")
	fmt.Println("  clear_check <check_id>")
	fmt.Println("  add_check_note <check_id> \"<note>\"")
	fmt.Println()
	fmt.Println("Asset Tracking:")
	fmt.Println("  list_hardware <device_id>")
	fmt.Println("  list_software <device_id>")
	fmt.Println("  list_device_asset_details <device_id>")
	fmt.Println("  list_license_groups")
	fmt.Println()
	fmt.Println("Patch Management:")
	fmt.Println("  list_patches <device_id>")
	fmt.Println("  approve_patch <device_id> <patch_id1,patch_id2,...>")
	fmt.Println("  ignore_patch <device_id> <patch_id1,patch_id2,...>")
	fmt.Println()
	fmt.Println("Antivirus:")
	fmt.Println("  list_antivirus_products")
	fmt.Println("  list_antivirus_definitions <device_id>")
	fmt.Println("  list_quarantine <device_id>")
	fmt.Println("  start_scan <device_id> <scan_type>")
	fmt.Println()
	fmt.Println("Performance and History:")
	fmt.Println("  list_performance_history <device_id> <check_id> <start_date> <end_date>")
	fmt.Println("  list_drive_space_history <device_id> <start_date> <end_date>")
	fmt.Println()
	fmt.Println("Templates:")
	fmt.Println("  list_templates")
	fmt.Println()
	fmt.Println("Backup & Recovery:")
	fmt.Println("  list_backup_sessions <device_id>")
	fmt.Println()
	fmt.Println("Settings:")
	fmt.Println("  list_wall_chart_settings")
	fmt.Println("  list_general_settings")
	fmt.Println()
	fmt.Println("Tasks and Users:")
	fmt.Println("  list_active_directory_users <device_id>")
	fmt.Println("  run_task_now <task_id>")
	fmt.Println()
	fmt.Println("Site Management:")
	fmt.Println("  add_client \"<name>\" \"<contact_name>\" \"<contact_email>\"")
	fmt.Println("  add_site <client_id> \"<name>\" \"<contact_name>\" \"<contact_email>\"")
	fmt.Println("  get_site_installation_package <site_id> <package_type>")
}

func outputJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling data to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func getClientIDByName(apiClient *nsight.ApiClient, clientName string) (int, error) {
	clients, err := apiClient.FetchClients()
	if err != nil {
		return 0, fmt.Errorf("could not fetch client list to find ID: %w", err)
	}

	for _, client := range clients {
		if client.Name == clientName {
			return client.ClientID, nil
		}
	}

	return 0, fmt.Errorf("client with name '%s' not found", clientName)
}

func getSiteIDByName(apiClient *nsight.ApiClient, siteName string) (int, error) {
	clients, err := apiClient.FetchClients()
	if err != nil {
		return 0, fmt.Errorf("could not fetch client list to find site ID: %w", err)
	}

	for _, client := range clients {
		sites, err := apiClient.FetchSites(client.ClientID)
		if err != nil {
			continue
		}
		for _, site := range sites {
			if site.Name == siteName {
				return site.SiteID, nil
			}
		}
	}

	return 0, fmt.Errorf("site with name '%s' not found across any client", siteName)
}

func resolveClientID(apiClient *nsight.ApiClient, identifier string) (int, error) {
	if clientID, err := strconv.Atoi(identifier); err == nil {
		return clientID, nil
	}
	return getClientIDByName(apiClient, identifier)
}

func resolveSiteID(apiClient *nsight.ApiClient, identifier string) (int, error) {
	if siteID, err := strconv.Atoi(identifier); err == nil {
		return siteID, nil
	}
	return getSiteIDByName(apiClient, identifier)
}

func resolveDeviceID(identifier string) (int, error) {
	return strconv.Atoi(identifier)
}

// -- Handler Functions --

func handleListClients(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_clients")
	}
	clients, err := apiClient.FetchClients()
	if err != nil {
		log.Fatalf("Error fetching clients: %v", err)
	}
	outputJSON(clients)
}

func handleListSites(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_sites <client_id | \"client_name\">")
	}
	clientID, err := resolveClientID(apiClient, args[0])
	if err != nil {
		log.Fatalf("Error resolving client: %v", err)
	}
	sites, err := apiClient.FetchSites(clientID)
	if err != nil {
		log.Fatalf("Error fetching sites: %v", err)
	}
	outputJSON(sites)
}

func handleListServers(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_servers <site_id | \"site_name\">")
	}
	siteID, err := resolveSiteID(apiClient, args[0])
	if err != nil {
		log.Fatalf("Error resolving site: %v", err)
	}
	servers, err := apiClient.FetchServers(siteID)
	if err != nil {
		log.Fatalf("Error fetching servers: %v", err)
	}
	outputJSON(servers)
}

func handleListWorkstations(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_workstations <site_id | \"site_name\">")
	}
	siteID, err := resolveSiteID(apiClient, args[0])
	if err != nil {
		log.Fatalf("Error resolving site: %v", err)
	}
	workstations, err := apiClient.FetchWorkstations(siteID)
	if err != nil {
		log.Fatalf("Error fetching workstations: %v", err)
	}
	outputJSON(workstations)
}

func handleListDevices(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_devices <site_id>")
	}
	siteID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid site ID: %v", err)
	}
	devices, err := apiClient.FetchDevicesBySite(siteID)
	if err != nil {
		log.Fatalf("Error fetching devices: %v", err)
	}
	outputJSON(devices)
}

func handleListDevicesAtClient(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_devices_at_client <client_id>")
	}
	clientID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid client ID: %v", err)
	}
	devices, err := apiClient.FetchDevices(clientID)
	if err != nil {
		log.Fatalf("Error fetching devices: %v", err)
	}
	outputJSON(devices)
}

func handleListAgentlessAssets(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_agentless_assets <site_id>")
	}
	siteID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid site ID: %v", err)
	}
	assets, err := apiClient.FetchAgentlessAssets(siteID)
	if err != nil {
		log.Fatalf("Error fetching agentless assets: %v", err)
	}
	outputJSON(assets)
}

func handleListFailingChecks(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_failing_checks")
	}
	checks, err := apiClient.FetchFailingChecks()
	if err != nil {
		log.Fatalf("Error fetching failing checks: %v", err)
	}
	outputJSON(checks)
}

func handleListChecks(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_checks <device_id | site_id>")
	}
	
	// Try as device ID first, then as site ID
	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid ID: %v", err)
	}
	
	// Try device checks first
	checks, err := apiClient.FetchChecks(id)
	if err != nil {
		// If device checks fail, try site checks
		checks, err = apiClient.FetchChecksBySite(id)
		if err != nil {
			log.Fatalf("Error fetching checks: %v", err)
		}
	}
	outputJSON(checks)
}

func handleListDeviceMonitoringDetails(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_device_monitoring_details <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	details, err := apiClient.FetchDeviceMonitoringDetails(deviceID)
	if err != nil {
		log.Fatalf("Error fetching device monitoring details: %v", err)
	}
	outputJSON(details)
}

func handleListCheckConfiguration(apiClient *nsight.ApiClient, serviceName string, args []string) {
	if len(args) < 1 {
		log.Fatalf("Usage: %s <device_id> [os]", serviceName)
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	
	var os string
	if strings.Contains(serviceName, "_windows") {
		os = "windows"
	} else if strings.Contains(serviceName, "_mac") {
		os = "mac"
	} else if strings.Contains(serviceName, "_linux") {
		os = "linux"
	} else if len(args) > 1 {
		os = args[1]
	}
	
	config, err := apiClient.FetchCheckConfiguration(deviceID, os)
	if err != nil {
		log.Fatalf("Error fetching check configuration: %v", err)
	}
	outputJSON(config)
}

func handleListOutages(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 3 {
		log.Fatalf("Usage: list_outages <site_id> <start_date> <end_date>")
	}
	siteID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid site ID: %v", err)
	}
	outages, err := apiClient.FetchOutages(siteID, args[1], args[2])
	if err != nil {
		log.Fatalf("Error fetching outages: %v", err)
	}
	outputJSON(outages)
}

func handleClearCheck(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: clear_check <check_id>")
	}
	checkID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid check ID: %v", err)
	}
	err = apiClient.ClearCheck(checkID)
	if err != nil {
		log.Fatalf("Error clearing check: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Check cleared\"}")
}

func handleAddCheckNote(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: add_check_note <check_id> \"<note>\"")
	}
	checkID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid check ID: %v", err)
	}
	err = apiClient.AddCheckNote(checkID, args[1])
	if err != nil {
		log.Fatalf("Error adding check note: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Note added to check\"}")
}

func handleListHardware(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_hardware <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	hardware, err := apiClient.FetchHardware(deviceID)
	if err != nil {
		log.Fatalf("Error fetching hardware: %v", err)
	}
	outputJSON(hardware)
}

func handleListSoftware(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_software <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	software, err := apiClient.FetchSoftware(deviceID)
	if err != nil {
		log.Fatalf("Error fetching software: %v", err)
	}
	outputJSON(software)
}

func handleListDeviceAssetDetails(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_device_asset_details <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	details, err := apiClient.FetchDeviceAssetDetails(deviceID)
	if err != nil {
		log.Fatalf("Error fetching device asset details: %v", err)
	}
	outputJSON(details)
}

func handleListLicenseGroups(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_license_groups")
	}
	groups, err := apiClient.FetchLicenseGroups()
	if err != nil {
		log.Fatalf("Error fetching license groups: %v", err)
	}
	outputJSON(groups)
}

func handleListPatches(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_patches <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	patches, err := apiClient.FetchPatches(deviceID)
	if err != nil {
		log.Fatalf("Error fetching patches: %v", err)
	}
	outputJSON(patches)
}

func handleApprovePatches(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: approve_patch <device_id> <patch_id1,patch_id2,...>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	
	patchIDStrings := strings.Split(args[1], ",")
	patchIDs := make([]int, len(patchIDStrings))
	for i, idStr := range patchIDStrings {
		patchIDs[i], err = strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			log.Fatalf("Invalid patch ID '%s': %v", idStr, err)
		}
	}
	
	err = apiClient.ApprovePatches(deviceID, patchIDs)
	if err != nil {
		log.Fatalf("Error approving patches: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Patches approved\"}")
}

func handleIgnorePatches(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: ignore_patch <device_id> <patch_id1,patch_id2,...>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	
	patchIDStrings := strings.Split(args[1], ",")
	patchIDs := make([]int, len(patchIDStrings))
	for i, idStr := range patchIDStrings {
		patchIDs[i], err = strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			log.Fatalf("Invalid patch ID '%s': %v", idStr, err)
		}
	}
	
	err = apiClient.IgnorePatches(deviceID, patchIDs)
	if err != nil {
		log.Fatalf("Error ignoring patches: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Patches ignored\"}")
}

func handleListAntivirusProducts(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_antivirus_products")
	}
	products, err := apiClient.FetchAntivirusProducts()
	if err != nil {
		log.Fatalf("Error fetching antivirus products: %v", err)
	}
	outputJSON(products)
}

func handleListAntivirusDefinitions(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_antivirus_definitions <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	definitions, err := apiClient.FetchAntivirusDefinitions(deviceID)
	if err != nil {
		log.Fatalf("Error fetching antivirus definitions: %v", err)
	}
	outputJSON(definitions)
}

func handleListQuarantine(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_quarantine <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	quarantine, err := apiClient.FetchQuarantineList(deviceID)
	if err != nil {
		log.Fatalf("Error fetching quarantine list: %v", err)
	}
	outputJSON(quarantine)
}

func handleStartAntivirusScan(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: start_scan <device_id> <scan_type>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	err = apiClient.StartAntivirusScan(deviceID, args[1])
	if err != nil {
		log.Fatalf("Error starting antivirus scan: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Antivirus scan started\"}")
}

func handleListPerformanceHistory(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 4 {
		log.Fatalf("Usage: list_performance_history <device_id> <check_id> <start_date> <end_date>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	checkID, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Invalid check ID: %v", err)
	}
	history, err := apiClient.FetchPerformanceHistory(deviceID, checkID, args[2], args[3])
	if err != nil {
		log.Fatalf("Error fetching performance history: %v", err)
	}
	outputJSON(history)
}

func handleListDriveSpaceHistory(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 3 {
		log.Fatalf("Usage: list_drive_space_history <device_id> <start_date> <end_date>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	history, err := apiClient.FetchDriveSpaceHistory(deviceID, args[1], args[2])
	if err != nil {
		log.Fatalf("Error fetching drive space history: %v", err)
	}
	outputJSON(history)
}

func handleListTemplates(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_templates")
	}
	templates, err := apiClient.FetchTemplates()
	if err != nil {
		log.Fatalf("Error fetching templates: %v", err)
	}
	outputJSON(templates)
}

func handleListBackupSessions(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_backup_sessions <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	sessions, err := apiClient.FetchBackupSessions(deviceID)
	if err != nil {
		log.Fatalf("Error fetching backup sessions: %v", err)
	}
	outputJSON(sessions)
}

func handleListWallChartSettings(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_wall_chart_settings")
	}
	settings, err := apiClient.FetchWallChartSettings()
	if err != nil {
		log.Fatalf("Error fetching wall chart settings: %v", err)
	}
	outputJSON(settings)
}

func handleListGeneralSettings(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 0 {
		log.Fatalf("Usage: list_general_settings")
	}
	settings, err := apiClient.FetchGeneralSettings()
	if err != nil {
		log.Fatalf("Error fetching general settings: %v", err)
	}
	outputJSON(settings)
}

func handleListActiveDirectoryUsers(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: list_active_directory_users <device_id>")
	}
	deviceID, err := resolveDeviceID(args[0])
	if err != nil {
		log.Fatalf("Invalid device ID: %v", err)
	}
	users, err := apiClient.FetchActiveDirectoryUsers(deviceID)
	if err != nil {
		log.Fatalf("Error fetching Active Directory users: %v", err)
	}
	outputJSON(users)
}

func handleRunTaskNow(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 1 {
		log.Fatalf("Usage: run_task_now <task_id>")
	}
	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid task ID: %v", err)
	}
	err = apiClient.RunTaskNow(taskID)
	if err != nil {
		log.Fatalf("Error running task: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Task started\"}")
}

func handleAddClient(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 3 {
		log.Fatalf("Usage: add_client \"<name>\" \"<contact_name>\" \"<contact_email>\"")
	}
	err := apiClient.AddClient(args[0], args[1], args[2])
	if err != nil {
		log.Fatalf("Error adding client: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Client added\"}")
}

func handleAddSite(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 4 {
		log.Fatalf("Usage: add_site <client_id> \"<name>\" \"<contact_name>\" \"<contact_email>\"")
	}
	clientID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid client ID: %v", err)
	}
	err = apiClient.AddSite(clientID, args[1], args[2], args[3])
	if err != nil {
		log.Fatalf("Error adding site: %v", err)
	}
	fmt.Println("{\"status\": \"success\", \"message\": \"Site added\"}")
}

func handleGetSiteInstallationPackage(apiClient *nsight.ApiClient, args []string) {
	if len(args) != 2 {
		log.Fatalf("Usage: get_site_installation_package <site_id> <package_type>")
	}
	siteID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Invalid site ID: %v", err)
	}
	packageData, err := apiClient.GetSiteInstallationPackage(siteID, args[1])
	if err != nil {
		log.Fatalf("Error getting installation package: %v", err)
	}
	// For binary data, we could base64 encode or save to file
	fmt.Printf("{\"status\": \"success\", \"package_size\": %d}\n", len(packageData))
}
