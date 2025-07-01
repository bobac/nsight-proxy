# N-Sight JSON Proxy

N-Sight JSON Proxy je HTTP server, který funguje jako proxy a převodník pro N-Sight API. Přijímá HTTP požadavky stejně jako originální N-Sight API, ale vrací data ve formátu JSON místo XML.

## Funkcionalita

- **HTTP Server**: Poslouchá na portu 80
- **API Proxy**: Předává volání na nakonfigurovaný N-Sight endpoint
- **XML→JSON Konverze**: Automaticky převádí XML odpovědi na JSON formát
- **Plná kompatibilita**: Podporuje všechny dostupné N-Sight API služby
- **CORS podpora**: Umožňuje cross-origin requests pro webové aplikace

## Konfigurace

Proxy server vyžaduje pouze konfiguraci N-Sight serveru. API klíč se předává v každém požadavku:

### Proměnné prostředí

```bash
NSIGHT_SERVER=your.nsight.server.com
```

### .env soubor

```env
NSIGHT_SERVER=your.nsight.server.com
```

**Poznámka**: API klíč se nepředává přes .env soubor, ale musí být součástí každého HTTP požadavku jako parametr `apikey`. Tím se zajistí bezpečnost - každý uživatel musí použít svůj vlastní API klíč.

## Spuštění

```bash
# Spustit proxy server
./nsight-proxy

# Nebo pomocí go run
go run cmd/nsight-proxy/main.go
```

Server se spustí na portu 80 a bude dostupný na:
- API endpoint: `http://localhost/api/`
- Health check: `http://localhost/health`
- Info endpoint: `http://localhost/`

## Použití

Formát volání je stejný jako originální N-Sight API, pouze s JSON výstupem:

```
http://localhost/api/?apikey=<your_api_key>&service=<service_name>&<parameters>
```

**Důležité**: Každý požadavek musí obsahovat parametr `apikey` s platným N-Sight API klíčem.

### Příklady volání

#### Získání seznamu klientů
```bash
curl "http://localhost/api/?apikey=YOUR_API_KEY&service=list_clients"
```

#### Získání serverů pro konkrétní site
```bash
curl "http://localhost/api/?apikey=YOUR_API_KEY&service=list_servers&siteid=123"
```

#### Získání informací o zařízení
```bash
curl "http://localhost/api/?apikey=YOUR_API_KEY&service=list_device_asset_details&deviceid=456"
```

## Podporované služby

Proxy server podporuje všechny služby dostupné v původním N-Sight API:

### Clients, Sites a Devices
- `list_clients` - Seznam všech klientů
- `list_sites` - Seznam site pro klienta (parametr: `clientid`)
- `list_servers` - Seznam serverů pro site (parametr: `siteid`)
- `list_workstations` - Seznam workstation pro site (parametr: `siteid`)
- `list_devices` - Seznam zařízení pro site (parametr: `siteid`)
- `list_devices_at_client` - Seznam zařízení pro klienta (parametr: `clientid`)
- `list_device_asset_details` - Detaily asset pro zařízení (parametr: `deviceid`)
- `list_device_monitoring_details` - Monitoring detaily zařízení (parametr: `deviceid`)
- `list_agentless_assets` - Agentless assets pro site (parametr: `siteid`)

### Checks a Results
- `list_failing_checks` - Seznam selhávajících kontrol
- `list_checks` - Seznam kontrol pro zařízení/site (parametry: `deviceid` nebo `siteid`)

### Asset Tracking
- `list_hardware` - Hardware informace (parametr: `deviceid`)
- `list_software` - Software informace (parametr: `deviceid`)
- `list_license_groups` - Skupiny licencí

### Patch Management
- `list_patches` - Seznam patchů pro zařízení (parametr: `deviceid`)

### Antivirus
- `list_antivirus_products` - Podporované antivirus produkty
- `list_antivirus_definitions` - Definice antiviru (parametr: `deviceid`)
- `list_quarantine` - Seznam karantény (parametr: `deviceid`)

### Performance
- `list_performance_history` - Historie výkonu (parametry: `deviceid`, `checkid`, `startdate`, `enddate`)
- `list_drive_space_history` - Historie místa na disku (parametry: `deviceid`, `startdate`, `enddate`)

### Templates
- `list_templates` - Seznam monitorovacích šablon

## Response Format

Všechny odpovědi jsou ve formátu JSON:

### Úspěšná odpověď
```json
[
  {
    "client_id": 123,
    "client_name": "Example Client",
    "contact_name": "John Doe",
    ...
  }
]
```

### Chybová odpověď
```json
{
  "error": "Error description"
}
```

## Health Check

Server poskytuje health check endpoint pro monitoring:

```bash
curl http://localhost/health
```

Odpověď:
```json
{
  "status": "ok",
  "service": "nsight-proxy"
}
```

## CORS podpora

Server automaticky přidává CORS hlavičky pro podporu webových aplikací:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type`

## Logování

Server loguje všechny příchozí požadavky a chyby do standardního výstupu:

```
2024/01/01 12:00:00 Starting N-Sight JSON Proxy Server...
2024/01/01 12:00:00 Server starting on port 80...
2024/01/01 12:00:01 Handling request for service: list_clients
```

## Bezpečnost

- **API klíč v URL**: Každý uživatel musí poskytnout svůj vlastní API klíč v každém požadavku
- **Žádné sdílené klíče**: Server neuchovává žádné API klíče, což eliminuje bezpečnostní rizika
- **Validace**: Server validuje přítomnost API klíče a serveru před zpracováním požadavku
- **Chybové zprávy**: Všechny chyby jsou bezpečně zpracovány bez odhalení citlivých informací
- **CORS**: Nakonfigurován permisivně pro vývojové účely - v produkci doporučujeme omezit domény

**Důležité pro produkční nasazení:**
- Používejte HTTPS pro ochranu API klíče v přenosu
- Omezte CORS politiky na konkrétní domény
- Implementujte rate limiting pro ochranu před zneužitím
- Monitorujte přístup k API endpointům

## Build

Nástroj se automaticky builduje pomocí build scriptu:

```bash
./build.sh
```

Výsledný binární soubor bude dostupný v `bin/` adresáři pro všechny podporované platformy.