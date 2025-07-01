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
	"strconv"

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

// NewApiClientWithCredentials creates a new ApiClient with provided API key and server
func NewApiClientWithCredentials(apiKey, server string) (*ApiClient, error) {
	if apiKey == "" || server == "" {
		return nil, errors.New("apiKey and server must not be empty")
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

// FetchDeviceAssetDetails fetches asset details for a specific deviceID
func (c *ApiClient) FetchDeviceAssetDetails(deviceID int) (*AssetDetails, error) {
	params := map[string]string{
		"service":  "list_device_asset_details",
		"deviceid": strconv.Itoa(deviceID),
	}

	body, err := c.callAPI("list_device_asset_details", params)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request for device asset details (device %d): %w", deviceID, err)
	}

	var result AssetDetails
	// Use decodeXML to handle potential charset issues like ISO-8859-1
	if err := decodeXML(body, &result); err != nil {
		// Log the body for debugging if unmarshal fails
		log.Printf("Failed to decode device asset details XML for device %d. Body: %s", deviceID, string(body))
		return nil, fmt.Errorf("failed to decode device asset details XML for device %d: %w", deviceID, err)
	}

	// Basic check for empty results (though API might return OK with empty fields)
	if result.Manufacturer == "" && result.Model == "" && len(result.Hardware) == 0 { // Heuristic check
		log.Printf("Warning: Received potentially empty or incomplete asset details for device %d", deviceID)
		// Decide if this should be an error or just return the potentially empty struct
	}

	return &result, nil
}

// -- Check-related methods --

// FetchFailingChecks fetches all failing checks
func (c *ApiClient) FetchFailingChecks() ([]Check, error) {
	bodyBytes, err := c.callAPI("list_failing_checks", nil)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchChecks fetches all checks for a device or site
func (c *ApiClient) FetchChecks(deviceID int) ([]Check, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_checks", params)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchChecksBySite fetches all checks for a site
func (c *ApiClient) FetchChecksBySite(siteID int) ([]Check, error) {
	params := map[string]string{"siteid": fmt.Sprintf("%d", siteID)}
	bodyBytes, err := c.callAPI("list_checks", params)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ClearCheck clears a specific check
func (c *ApiClient) ClearCheck(checkID int) error {
	params := map[string]string{"checkid": fmt.Sprintf("%d", checkID)}
	_, err := c.callAPI("clear_check", params)
	return err
}

// AddCheckNote adds a note to a check
func (c *ApiClient) AddCheckNote(checkID int, note string) error {
	params := map[string]string{
		"checkid": fmt.Sprintf("%d", checkID),
		"note":    note,
	}
	_, err := c.callAPI("add_check_note", params)
	return err
}

// -- Device-related methods --

// FetchDevices fetches devices for a client
func (c *ApiClient) FetchDevices(clientID int) ([]Device, error) {
	params := map[string]string{"clientid": fmt.Sprintf("%d", clientID)}
	bodyBytes, err := c.callAPI("list_devices_at_client", params)
	if err != nil {
		return nil, err
	}
	var result DeviceResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchDevicesBySite fetches devices for a site
func (c *ApiClient) FetchDevicesBySite(siteID int) ([]Device, error) {
	params := map[string]string{"siteid": fmt.Sprintf("%d", siteID)}
	bodyBytes, err := c.callAPI("list_devices", params)
	if err != nil {
		return nil, err
	}
	var result DeviceResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchDeviceMonitoringDetails fetches monitoring details for a device
func (c *ApiClient) FetchDeviceMonitoringDetails(deviceID int) ([]Check, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_device_monitoring_details", params)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Agentless Asset methods --

// FetchAgentlessAssets fetches agentless assets for a site
func (c *ApiClient) FetchAgentlessAssets(siteID int) ([]AgentlessAsset, error) {
	params := map[string]string{"siteid": fmt.Sprintf("%d", siteID)}
	bodyBytes, err := c.callAPI("list_agentless_assets", params)
	if err != nil {
		return nil, err
	}
	var result AgentlessAssetResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Asset Tracking methods --

// FetchHardware fetches hardware information for a device
func (c *ApiClient) FetchHardware(deviceID int) ([]HardwareItem, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_hardware", params)
	if err != nil {
		return nil, err
	}
	var result struct {
		XMLName xml.Name       `xml:"result"`
		Items   []HardwareItem `xml:"items>item"`
	}
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchSoftware fetches software information for a device
func (c *ApiClient) FetchSoftware(deviceID int) ([]SoftwareItem, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_software", params)
	if err != nil {
		return nil, err
	}
	var result struct {
		XMLName xml.Name       `xml:"result"`
		Items   []SoftwareItem `xml:"items>item"`
	}
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchLicenseGroups fetches license groups
func (c *ApiClient) FetchLicenseGroups() ([]LicenseGroup, error) {
	bodyBytes, err := c.callAPI("list_license_groups", nil)
	if err != nil {
		return nil, err
	}
	var result LicenseGroupResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Patch Management methods --

// FetchPatches fetches all patches for a device
func (c *ApiClient) FetchPatches(deviceID int) ([]Patch, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_patches", params)
	if err != nil {
		return nil, err
	}
	var result PatchResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// ApprovePatches approves patches for a device
func (c *ApiClient) ApprovePatches(deviceID int, patchIDs []int) error {
	patchIDsStr := make([]string, len(patchIDs))
	for i, id := range patchIDs {
		patchIDsStr[i] = fmt.Sprintf("%d", id)
	}
	params := map[string]string{
		"deviceid": fmt.Sprintf("%d", deviceID),
		"patchids": fmt.Sprintf("[%s]", fmt.Sprintf("%s", patchIDsStr)),
	}
	_, err := c.callAPI("approve_patch", params)
	return err
}

// IgnorePatches ignores patches for a device
func (c *ApiClient) IgnorePatches(deviceID int, patchIDs []int) error {
	patchIDsStr := make([]string, len(patchIDs))
	for i, id := range patchIDs {
		patchIDsStr[i] = fmt.Sprintf("%d", id)
	}
	params := map[string]string{
		"deviceid": fmt.Sprintf("%d", deviceID),
		"patchids": fmt.Sprintf("[%s]", fmt.Sprintf("%s", patchIDsStr)),
	}
	_, err := c.callAPI("ignore_patch", params)
	return err
}

// -- Antivirus methods --

// FetchAntivirusProducts fetches supported antivirus products
func (c *ApiClient) FetchAntivirusProducts() ([]AntivirusProduct, error) {
	bodyBytes, err := c.callAPI("list_antivirus_products", nil)
	if err != nil {
		return nil, err
	}
	var result AntivirusProductResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchAntivirusDefinitions fetches antivirus definition information
func (c *ApiClient) FetchAntivirusDefinitions(deviceID int) ([]AntivirusDefinition, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_antivirus_definitions", params)
	if err != nil {
		return nil, err
	}
	var result AntivirusDefinitionResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchQuarantineList fetches quarantined items
func (c *ApiClient) FetchQuarantineList(deviceID int) ([]QuarantineItem, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_quarantine", params)
	if err != nil {
		return nil, err
	}
	var result QuarantineResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// StartAntivirusScan starts an antivirus scan
func (c *ApiClient) StartAntivirusScan(deviceID int, scanType string) error {
	params := map[string]string{
		"deviceid": fmt.Sprintf("%d", deviceID),
		"scantype": scanType,
	}
	_, err := c.callAPI("start_scan", params)
	return err
}

// -- Performance History methods --

// FetchPerformanceHistory fetches performance monitoring data
func (c *ApiClient) FetchPerformanceHistory(deviceID, checkID int, startDate, endDate string) ([]PerformanceData, error) {
	params := map[string]string{
		"deviceid":  fmt.Sprintf("%d", deviceID),
		"checkid":   fmt.Sprintf("%d", checkID),
		"startdate": startDate,
		"enddate":   endDate,
	}
	bodyBytes, err := c.callAPI("list_performance_history", params)
	if err != nil {
		return nil, err
	}
	var result PerformanceResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchDriveSpaceHistory fetches drive space history
func (c *ApiClient) FetchDriveSpaceHistory(deviceID int, startDate, endDate string) ([]PerformanceData, error) {
	params := map[string]string{
		"deviceid":  fmt.Sprintf("%d", deviceID),
		"startdate": startDate,
		"enddate":   endDate,
	}
	bodyBytes, err := c.callAPI("list_drive_space_history", params)
	if err != nil {
		return nil, err
	}
	var result PerformanceResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Template methods --

// FetchTemplates fetches monitoring templates
func (c *ApiClient) FetchTemplates() ([]Template, error) {
	bodyBytes, err := c.callAPI("list_templates", nil)
	if err != nil {
		return nil, err
	}
	var result TemplateResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Backup methods --

// FetchBackupSessions fetches backup & recovery sessions
func (c *ApiClient) FetchBackupSessions(deviceID int) ([]BackupSession, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_backup_sessions", params)
	if err != nil {
		return nil, err
	}
	var result BackupSessionResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Settings methods --

// FetchWallChartSettings fetches wall chart settings
func (c *ApiClient) FetchWallChartSettings() ([]Setting, error) {
	bodyBytes, err := c.callAPI("list_wall_chart_settings", nil)
	if err != nil {
		return nil, err
	}
	var result SettingResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// FetchGeneralSettings fetches general settings
func (c *ApiClient) FetchGeneralSettings() ([]Setting, error) {
	bodyBytes, err := c.callAPI("list_general_settings", nil)
	if err != nil {
		return nil, err
	}
	var result SettingResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Task methods --

// FetchActiveDirectoryUsers fetches Active Directory users
func (c *ApiClient) FetchActiveDirectoryUsers(deviceID int) ([]ADUser, error) {
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI("list_active_directory_users", params)
	if err != nil {
		return nil, err
	}
	var result ADUserResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// RunTaskNow runs a task immediately
func (c *ApiClient) RunTaskNow(taskID int) error {
	params := map[string]string{"taskid": fmt.Sprintf("%d", taskID)}
	_, err := c.callAPI("run_task_now", params)
	return err
}

// -- Site Management methods --

// AddClient adds a new client
func (c *ApiClient) AddClient(name, contactName, contactEmail string) error {
	params := map[string]string{
		"name":         name,
		"contactname":  contactName,
		"contactemail": contactEmail,
	}
	_, err := c.callAPI("add_client", params)
	return err
}

// AddSite adds a new site to a client
func (c *ApiClient) AddSite(clientID int, name, contactName, contactEmail string) error {
	params := map[string]string{
		"clientid":     fmt.Sprintf("%d", clientID),
		"name":         name,
		"contactname":  contactName,
		"contactemail": contactEmail,
	}
	_, err := c.callAPI("add_site", params)
	return err
}

// GetSiteInstallationPackage gets installation package for a site
func (c *ApiClient) GetSiteInstallationPackage(siteID int, packageType string) ([]byte, error) {
	params := map[string]string{
		"siteid":      fmt.Sprintf("%d", siteID),
		"packagetype": packageType,
	}
	return c.callAPI("get_site_installation_package", params)
}

// -- Outage methods --

// FetchOutages fetches system outages
func (c *ApiClient) FetchOutages(siteID int, startDate, endDate string) ([]Check, error) {
	params := map[string]string{
		"siteid":    fmt.Sprintf("%d", siteID),
		"startdate": startDate,
		"enddate":   endDate,
	}
	bodyBytes, err := c.callAPI("list_outages", params)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

// -- Check Configuration methods --

// FetchCheckConfiguration fetches check configuration for a device
func (c *ApiClient) FetchCheckConfiguration(deviceID int, os string) ([]Check, error) {
	serviceName := "list_check_configuration"
	if os != "" {
		serviceName = fmt.Sprintf("list_check_configuration_%s", os)
	}
	params := map[string]string{"deviceid": fmt.Sprintf("%d", deviceID)}
	bodyBytes, err := c.callAPI(serviceName, params)
	if err != nil {
		return nil, err
	}
	var result CheckResult
	if err := decodeXML(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}
