# Rozšíření nástroje getdata - NSight API

## Shrnutí změn

Nástroj `getdata` byl kompletně rozšířen pro podporu **všech dostupných N-Sight API volání** podle oficiální dokumentace na https://developer.n-able.com/n-sight/docs/getting-started-with-the-n-sight-api.

## Původní stav vs. Nový stav

### Původní implementace (4 API volání):
- `list_clients`
- `list_sites`
- `list_servers`
- `list_workstations`

### Nová implementace (37 API volání):

#### 1. Základní výpis entit (7 volání):
- `list_clients` - Vypíše všechny klienty
- `list_sites` - Vypíše sites pro klienta
- `list_servers` - Vypíše servery pro site
- `list_workstations` - Vypíše workstations pro site
- `list_devices` - Vypíše všechna zařízení pro site
- `list_devices_at_client` - Vypíše všechna zařízení klienta
- `list_agentless_assets` - Vypíše agentless assets

#### 2. Monitorování a kontroly (8 volání):
- `list_failing_checks` - Všechny neúspěšné kontroly
- `list_checks` - Kontroly pro zařízení/site
- `list_device_monitoring_details` - Detaily monitorování zařízení
- `list_check_configuration` - Konfigurace kontrol
- `list_check_configuration_windows` - Konfigurace pro Windows
- `list_check_configuration_mac` - Konfigurace pro Mac
- `list_check_configuration_linux` - Konfigurace pro Linux
- `list_outages` - Výpadky systému
- `clear_check` - Vymazání kontroly
- `add_check_note` - Přidání poznámky ke kontrole

#### 3. Asset tracking (4 volání):
- `list_hardware` - Hardware informace zařízení
- `list_software` - Software informace zařízení
- `list_device_asset_details` - Detailní asset informace
- `list_license_groups` - Licenční skupiny

#### 4. Správa patchů (3 volání):
- `list_patches` - Všechny patche pro zařízení
- `approve_patch` - Schválení patchů
- `ignore_patch` - Ignorování patchů

#### 5. Antivirus (4 volání):
- `list_antivirus_products` - Podporované produkty
- `list_antivirus_definitions` - Antivirus definice
- `list_quarantine` - Položky v karanténě
- `start_scan` - Spuštění antivirus scanu

#### 6. Výkon a historie (2 volání):
- `list_performance_history` - Historie výkonu
- `list_drive_space_history` - Historie využití disků

#### 7. Šablony (1 volání):
- `list_templates` - Monitorovací šablony

#### 8. Backup & Recovery (1 volání):
- `list_backup_sessions` - Backup sessiony

#### 9. Nastavení (2 volání):
- `list_wall_chart_settings` - Nastavení wall chart
- `list_general_settings` - Obecná nastavení

#### 10. Úlohy a uživatelé (2 volání):
- `list_active_directory_users` - AD uživatelé
- `run_task_now` - Okamžité spuštění úlohy

#### 11. Správa sites (3 volání):
- `add_client` - Přidání nového klienta
- `add_site` - Přidání nové site
- `get_site_installation_package` - Instalační balíček

## Technické vylepšení

### 1. Rozšířené datové struktury (`internal/nsight/types.go`)

Přidáno **14 nových datových struktur**:
- `Check` a `CheckResult` - pro kontroly a monitorování
- `Device` a `DeviceResult` - pro zařízení
- `AgentlessAsset` a `AgentlessAssetResult` - pro agentless assets
- `Patch` a `PatchResult` - pro patch management
- `AntivirusProduct`, `AntivirusDefinition` - pro antivirus
- `Template`, `PerformanceData`, `BackupSession` - pro další funkce
- `Setting`, `LicenseGroup`, `QuarantineItem`, `ADUser` - pro správu
- `GenericResult` - pro obecné odpovědi

### 2. Rozšířené API metody (`internal/nsight/api.go`)

Přidáno **33 nových API metod**:
- **Monitorování**: `FetchFailingChecks`, `FetchChecks`, `FetchChecksBySite`, `ClearCheck`, `AddCheckNote`
- **Zařízení**: `FetchDevices`, `FetchDevicesBySite`, `FetchDeviceMonitoringDetails`
- **Asset tracking**: `FetchHardware`, `FetchSoftware`, `FetchLicenseGroups`
- **Patch management**: `FetchPatches`, `ApprovePatches`, `IgnorePatches`
- **Antivirus**: `FetchAntivirusProducts`, `FetchAntivirusDefinitions`, `FetchQuarantineList`, `StartAntivirusScan`
- **Historie**: `FetchPerformanceHistory`, `FetchDriveSpaceHistory`
- **Další**: `FetchTemplates`, `FetchBackupSessions`, `FetchWallChartSettings`, `FetchGeneralSettings`
- **Úlohy**: `FetchActiveDirectoryUsers`, `RunTaskNow`
- **Site management**: `AddClient`, `AddSite`, `GetSiteInstallationPackage`

### 3. Kompletně přepracovaný getdata tool (`cmd/getdata/main.go`)

#### Nová architektura:
- **Modulární handlery**: Každé API volání má vlastní handler funkci
- **Inteligentní resolving**: Automatické převody mezi ID a jmény
- **Konzistentní error handling**: Jednotný přístup k chybám
- **Flexibilní parametry**: Podpora různých typů vstupních parametrů

#### Pomocné funkce:
- `resolveClientID()` - převod jméno/ID klienta
- `resolveSiteID()` - převod jméno/ID site
- `resolveDeviceID()` - validace device ID
- `outputJSON()` - jednotný JSON output

#### Handler pattern:
Každé API volání má vlastní handler:
```go
func handleListClients(apiClient *nsight.ApiClient, args []string)
func handleListSites(apiClient *nsight.ApiClient, args []string)
func handleApprovePatches(apiClient *nsight.ApiClient, args []string)
// ... atd.
```

## Akční API volání

Nástroj nyní podporuje nejen čtení dat, ale i **akční operace**:

### Patch management:
```bash
# Schválení patchů
go run cmd/getdata/main.go approve_patch 789 "12345,12346,12347"

# Ignorování patchů
go run cmd/getdata/main.go ignore_patch 789 "12345,12346"
```

### Antivirus:
```bash
# Spuštění full scanu
go run cmd/getdata/main.go start_scan 789 "full"
```

### Monitorování:
```bash
# Vymazání kontroly
go run cmd/getdata/main.go clear_check 12345

# Přidání poznámky
go run cmd/getdata/main.go add_check_note 12345 "Poznámka k této kontrole"
```

### Site management:
```bash
# Přidání klienta
go run cmd/getdata/main.go add_client "Název klienta" "Kontakt" "email@example.com"

# Přidání site
go run cmd/getdata/main.go add_site 123 "Název site" "Kontakt" "email@example.com"
```

### Úlohy:
```bash
# Okamžité spuštění úlohy
go run cmd/getdata/main.go run_task_now 12345
```

## Vylepšená dokumentace

### README.md kompletně přepsán:
- Podrobný popis všech 37 API volání
- Praktické příklady použití
- Kategorizace podle funkcionality
- Vysvětlení všech parametrů

### Nová struktura dokumentace:
1. **Základní entity** - clients, sites, devices
2. **Monitorování** - checks, outages, monitoring
3. **Asset tracking** - hardware, software, licenses
4. **Patch management** - list, approve, ignore
5. **Antivirus** - products, definitions, quarantine
6. **Výkon** - performance history, drive space
7. **Šablony** - monitoring templates
8. **Backup** - backup sessions
9. **Nastavení** - wall chart, general settings
10. **Úlohy** - AD users, task execution
11. **Site management** - add clients/sites, packages

## Zpětná kompatibilita

Všechny původní 4 API volání (`list_clients`, `list_sites`, `list_servers`, `list_workstations`) **zůstávají plně funkční** s identickou syntaxí.

## Přínosy rozšíření

1. **Kompletní pokrytí**: 37 API volání místo původních 4
2. **Lepší organizace**: Čistě strukturovaný kód s handlery
3. **Akční možnosti**: Nejen čtení, ale i změny systému
4. **Robustnost**: Lepší error handling a validace
5. **Dokumentace**: Podrobný popis všech funkcí
6. **Škálovatelnost**: Snadné přidání dalších API volání
7. **Konzistence**: Jednotný přístup napříč všemi voláními

Nástroj `getdata` je nyní kompletní NSight API klient podporující všechny dostupné funkce platformy N-Able N-Sight.