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
