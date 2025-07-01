package nsight

import "encoding/xml"

// -- Client Structures --

// Client represents the structure of a client in the XML response
type Client struct {
	XMLName  xml.Name `xml:"client"`
	ClientID int      `xml:"clientid"`
	Name     string   `xml:"name"`
	// Add other relevant fields based on actual API response if needed for caching/output
}

// ClientResult represents the top-level structure of the client list XML response
type ClientResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Client `xml:"items>client"`
}

// -- Site Structures --

// Site represents the structure of a site in the XML response
type Site struct {
	XMLName xml.Name `xml:"site"`
	SiteID  int      `xml:"siteid"`
	Name    string   `xml:"name"`
	// Add other relevant fields based on actual API response if needed
}

// SiteResult represents the top-level structure of the site list XML response
type SiteResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Site   `xml:"items>site"`
}

// -- Server Structures --

// Server represents the structure of a server in the XML response
type Server struct {
	XMLName      xml.Name `xml:"server"`
	ServerID     int      `xml:"serverid"`
	Name         string   `xml:"name"`
	Description  string   `xml:"description"`
	OS           string   `xml:"os"`
	IP           string   `xml:"ip"`
	Online       int      `xml:"online"`                   // 0 or 1
	User         string   `xml:"user,omitempty"`           // Added user
	Manufacturer string   `xml:"manufacturer,omitempty"`   // Added manufacturer
	Model        string   `xml:"model,omitempty"`          // Added model
	DeviceSerial string   `xml:"device_serial,omitempty"`  // Added device_serial
	LastBootTime string   `xml:"last_boot_time,omitempty"` // Added last_boot_time
	// Add other relevant fields based on actual API response if needed
}

// ServerResult represents the top-level structure of the server list XML response
type ServerResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Server `xml:"items>server"`
}

// -- Workstation Structures --

// Workstation represents the structure of a workstation in the XML response
type Workstation struct {
	XMLName       xml.Name `xml:"workstation"`
	WorkstationID int      `xml:"workstationid"`
	Name          string   `xml:"name"`
	Description   string   `xml:"description"`
	OS            string   `xml:"os"`
	IP            string   `xml:"ip"`
	Online        int      `xml:"online"`                   // 0 or 1
	User          string   `xml:"user,omitempty"`           // Added user
	Manufacturer  string   `xml:"manufacturer,omitempty"`   // Added manufacturer
	Model         string   `xml:"model,omitempty"`          // Added model
	DeviceSerial  string   `xml:"device_serial,omitempty"`  // Added device_serial
	LastBootTime  string   `xml:"last_boot_time,omitempty"` // Added last_boot_time
	// Add other relevant fields based on actual API response if needed
}

// WorkstationResult represents the top-level structure of the workstation list XML response
type WorkstationResult struct {
	XMLName xml.Name      `xml:"result"`
	Items   []Workstation `xml:"items>workstation"`
}

// -- Asset Detail Structures --

// HardwareItem represents a single hardware item in asset details
type HardwareItem struct {
	HardwareID   int    `xml:"hardwareid"`
	Name         string `xml:"name"`
	Type         int    `xml:"type"`
	Manufacturer string `xml:"manufacturer"`
	Details      string `xml:"details"`
	Status       string `xml:"status"`
	Deleted      int    `xml:"deleted"`
	Modified     int    `xml:"modified"`
}

// SoftwareItem represents a single software item in asset details
type SoftwareItem struct {
	Name        string `xml:"name"`
	SoftwareID  int    `xml:"softwareid"`
	Version     string `xml:"version"`
	InstallDate string `xml:"install_date"` // Keep as string, format later if needed
	Type        string `xml:"type"`
	Deleted     int    `xml:"deleted"`
	Modified    int    `xml:"modified"`
}

// AssetDetails represents the structure returned by list_device_asset_details
type AssetDetails struct {
	XMLName      xml.Name       `xml:"result"`
	Client       string         `xml:"client"`
	ChassisType  string         `xml:"chassistype"`
	IP           string         `xml:"ip"` // Note: Might differ from primary listing IP
	MAC1         string         `xml:"mac1"`
	MAC2         string         `xml:"mac2,omitempty"`
	MAC3         string         `xml:"mac3,omitempty"`
	User         string         `xml:"user"` // Note: Different meaning than primary listing user?
	Manufacturer string         `xml:"manufacturer"`
	Model        string         `xml:"model"`
	OS           string         `xml:"os"`
	SerialNumber string         `xml:"serialnumber"`
	ProductKey   string         `xml:"productkey,omitempty"`
	Role         string         `xml:"role"` // Keep as string, interpretation needed
	ServicePack  string         `xml:"servicepack"`
	RAM          int64          `xml:"ram"`      // Use int64 for potentially large numbers
	ScanTime     string         `xml:"scantime"` // Keep as string, format later if needed
	Custom1      CustomField    `xml:"custom1,omitempty"`
	Custom2      CustomField    `xml:"custom2,omitempty"`
	Custom3      CustomField    `xml:"custom3,omitempty"`
	Custom4      CustomField    `xml:"custom4,omitempty"`
	Custom5      CustomField    `xml:"custom5,omitempty"`
	Custom6      CustomField    `xml:"custom6,omitempty"`
	Custom7      CustomField    `xml:"custom7,omitempty"`
	Custom8      CustomField    `xml:"custom8,omitempty"`
	Custom9      CustomField    `xml:"custom9,omitempty"`
	Custom10     CustomField    `xml:"custom10,omitempty"`
	Hardware     []HardwareItem `xml:"hardware>item"`
	Software     []SoftwareItem `xml:"software>item"`
}

// CustomField handles the custom fields which have a name attribute and text value
type CustomField struct {
	Name  string `xml:"customname,attr"`
	Value string `xml:",chardata"`
}

// -- Check Structures --

// Check represents a check result
type Check struct {
	XMLName     xml.Name `xml:"check"`
	CheckID     int      `xml:"checkid"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	DeviceID    int      `xml:"deviceid"`
	DeviceName  string   `xml:"devicename"`
	State       int      `xml:"state"`
	Severity    int      `xml:"severity"`
	Message     string   `xml:"message"`
	LastCheck   string   `xml:"lastcheck"`
	NextCheck   string   `xml:"nextcheck"`
}

// CheckResult represents the top-level structure for check list responses
type CheckResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Check  `xml:"items>check"`
}

// -- Device Structures --

// Device represents a device in the system
type Device struct {
	XMLName    xml.Name `xml:"device"`
	DeviceID   int      `xml:"deviceid"`
	DeviceName string   `xml:"devicename"`
	DeviceType string   `xml:"devicetype"`
	ClientID   int      `xml:"clientid"`
	ClientName string   `xml:"clientname"`
	SiteID     int      `xml:"siteid"`
	SiteName   string   `xml:"sitename"`
	Online     int      `xml:"online"`
	OS         string   `xml:"os"`
	IP         string   `xml:"ip"`
}

// DeviceResult represents the top-level structure for device list responses
type DeviceResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Device `xml:"items>device"`
}

// -- Agentless Asset Structures --

// AgentlessAsset represents an agentless asset
type AgentlessAsset struct {
	XMLName   xml.Name `xml:"asset"`
	AssetID   int      `xml:"assetid"`
	Name      string   `xml:"name"`
	Type      string   `xml:"type"`
	IP        string   `xml:"ip"`
	MAC       string   `xml:"mac"`
	Vendor    string   `xml:"vendor"`
	SiteID    int      `xml:"siteid"`
	SiteName  string   `xml:"sitename"`
	Discovered string  `xml:"discovered"`
}

// AgentlessAssetResult represents the response for agentless assets
type AgentlessAssetResult struct {
	XMLName xml.Name         `xml:"result"`
	Items   []AgentlessAsset `xml:"items>asset"`
}

// -- Patch Structures --

// Patch represents a patch
type Patch struct {
	XMLName     xml.Name `xml:"patch"`
	PatchID     int      `xml:"patchid"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Severity    string   `xml:"severity"`
	Status      string   `xml:"status"`
	DeviceID    int      `xml:"deviceid"`
	DeviceName  string   `xml:"devicename"`
	Released    string   `xml:"released"`
	Installed   string   `xml:"installed"`
}

// PatchResult represents the response for patch lists
type PatchResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Patch  `xml:"items>patch"`
}

// -- Antivirus Structures --

// AntivirusProduct represents an antivirus product
type AntivirusProduct struct {
	XMLName   xml.Name `xml:"product"`
	ProductID int      `xml:"productid"`
	Name      string   `xml:"name"`
	Vendor    string   `xml:"vendor"`
	Supported int      `xml:"supported"`
}

// AntivirusProductResult represents the response for antivirus products
type AntivirusProductResult struct {
	XMLName xml.Name           `xml:"result"`
	Items   []AntivirusProduct `xml:"items>product"`
}

// AntivirusDefinition represents antivirus definition information
type AntivirusDefinition struct {
	XMLName     xml.Name `xml:"definition"`
	ProductID   int      `xml:"productid"`
	ProductName string   `xml:"productname"`
	Version     string   `xml:"version"`
	ReleaseDate string   `xml:"releasedate"`
	DeviceID    int      `xml:"deviceid"`
	DeviceName  string   `xml:"devicename"`
}

// AntivirusDefinitionResult represents the response for antivirus definitions
type AntivirusDefinitionResult struct {
	XMLName xml.Name              `xml:"result"`
	Items   []AntivirusDefinition `xml:"items>definition"`
}

// -- Template Structures --

// Template represents a monitoring template
type Template struct {
	XMLName      xml.Name `xml:"template"`
	TemplateID   int      `xml:"templateid"`
	Name         string   `xml:"name"`
	Description  string   `xml:"description"`
	OS           string   `xml:"os"`
	DeviceType   string   `xml:"devicetype"`
	CheckCount   int      `xml:"checkcount"`
	Created      string   `xml:"created"`
	Modified     string   `xml:"modified"`
}

// TemplateResult represents the response for template lists
type TemplateResult struct {
	XMLName xml.Name   `xml:"result"`
	Items   []Template `xml:"items>template"`
}

// -- Performance History Structures --

// PerformanceData represents performance monitoring data
type PerformanceData struct {
	XMLName   xml.Name `xml:"data"`
	CheckID   int      `xml:"checkid"`
	DeviceID  int      `xml:"deviceid"`
	Timestamp string   `xml:"timestamp"`
	Value     float64  `xml:"value"`
	Unit      string   `xml:"unit"`
}

// PerformanceResult represents the response for performance data
type PerformanceResult struct {
	XMLName xml.Name          `xml:"result"`
	Items   []PerformanceData `xml:"items>data"`
}

// -- Backup Structures --

// BackupSession represents a backup session
type BackupSession struct {
	XMLName       xml.Name `xml:"session"`
	SessionID     int      `xml:"sessionid"`
	DeviceID      int      `xml:"deviceid"`
	DeviceName    string   `xml:"devicename"`
	Type          string   `xml:"type"`
	Status        string   `xml:"status"`
	StartTime     string   `xml:"starttime"`
	EndTime       string   `xml:"endtime"`
	BytesTotal    int64    `xml:"bytestotal"`
	BytesBackedUp int64    `xml:"bytesbackedup"`
}

// BackupSessionResult represents the response for backup sessions
type BackupSessionResult struct {
	XMLName xml.Name        `xml:"result"`
	Items   []BackupSession `xml:"items>session"`
}

// -- Task Structures --

// Task represents an automated task
type Task struct {
	XMLName     xml.Name `xml:"task"`
	TaskID      int      `xml:"taskid"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Type        string   `xml:"type"`
	Status      string   `xml:"status"`
	DeviceID    int      `xml:"deviceid"`
	DeviceName  string   `xml:"devicename"`
	Scheduled   string   `xml:"scheduled"`
	LastRun     string   `xml:"lastrun"`
}

// TaskResult represents the response for task lists
type TaskResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []Task   `xml:"items>task"`
}

// -- Settings Structures --

// Setting represents a configuration setting
type Setting struct {
	XMLName xml.Name `xml:"setting"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
	Type    string   `xml:"type"`
	Section string   `xml:"section"`
}

// SettingResult represents the response for settings
type SettingResult struct {
	XMLName xml.Name  `xml:"result"`
	Items   []Setting `xml:"items>setting"`
}

// -- License Structures --

// LicenseGroup represents a license group
type LicenseGroup struct {
	XMLName   xml.Name `xml:"group"`
	GroupID   int      `xml:"groupid"`
	Name      string   `xml:"name"`
	Publisher string   `xml:"publisher"`
	Version   string   `xml:"version"`
	Count     int      `xml:"count"`
}

// LicenseGroupResult represents the response for license groups
type LicenseGroupResult struct {
	XMLName xml.Name       `xml:"result"`
	Items   []LicenseGroup `xml:"items>group"`
}

// -- Quarantine Structures --

// QuarantineItem represents a quarantined item
type QuarantineItem struct {
	XMLName      xml.Name `xml:"item"`
	ItemID       int      `xml:"itemid"`
	DeviceID     int      `xml:"deviceid"`
	DeviceName   string   `xml:"devicename"`
	ThreatName   string   `xml:"threatname"`
	FilePath     string   `xml:"filepath"`
	Quarantined  string   `xml:"quarantined"`
	Size         int64    `xml:"size"`
	ProductName  string   `xml:"productname"`
}

// QuarantineResult represents the response for quarantine items
type QuarantineResult struct {
	XMLName xml.Name         `xml:"result"`
	Items   []QuarantineItem `xml:"items>item"`
}

// -- Active Directory Structures --

// ADUser represents an Active Directory user
type ADUser struct {
	XMLName     xml.Name `xml:"user"`
	UserID      string   `xml:"userid"`
	Name        string   `xml:"name"`
	DisplayName string   `xml:"displayname"`
	Email       string   `xml:"email"`
	Domain      string   `xml:"domain"`
	LastLogon   string   `xml:"lastlogon"`
	Enabled     int      `xml:"enabled"`
}

// ADUserResult represents the response for AD users
type ADUserResult struct {
	XMLName xml.Name `xml:"result"`
	Items   []ADUser `xml:"items>user"`
}

// -- General Response Structure --

// GenericResult represents a generic API response structure
type GenericResult struct {
	XMLName xml.Name `xml:"result"`
	Status  string   `xml:"status,attr"`
	Message string   `xml:"message,omitempty"`
}
