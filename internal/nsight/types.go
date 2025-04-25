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
