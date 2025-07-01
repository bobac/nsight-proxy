# NSight Proxy

Tento projekt poskytuje nástroje příkazové řádky a (v budoucnu) API server pro interakci s N-Able N-Sight Data Extraction API. Hlavním cílem je usnadnit získávání dat v JSON formátu a agregovat informace z více API volání.

## Požadavky

*   Nainstalované **Go** (verze 1.21 nebo vyšší doporučena). Ověřte pomocí `go version`.
*   Soubor `.env` v kořenovém adresáři projektu.

## Nastavení

1.  **Naklonujte repozitář** (pokud jste tak již neučinili):
    ```bash
    git clone <URL_REPOZITARE>
    cd nsight-proxy
    ```
2.  **Vytvořte soubor `.env`**: Zkopírujte `.env.example` a vyplňte své údaje:
    ```bash
    cp .env.example .env
    ```
    Upravte soubor `.env` a zadejte platný `NSIGHT_API_KEY` a `NSIGHT_SERVER` (hostname, např. `wwweurope1.systemmonitor.eu.com`).

## Dostupné Nástroje

Projekt obsahuje tři hlavní nástroje v adresáři `cmd/`:

### 1. `getdata` - Komplexní API nástroj

Tento nástroj nyní podporuje **všechna dostupná N-Sight API volání** a slouží k přímému volání specifických N-Sight API služeb. Vrací výsledek jako JSON na standardní výstup.

**Základní použití:**

```bash
go run cmd/getdata/main.go <název_služby> [parametry...]
```

#### Základní výpis entit:

*   **`list_clients`**: Vypíše všechny klienty.
    ```bash
    go run cmd/getdata/main.go list_clients
    ```

*   **`list_sites`**: Vypíše všechny sites pro daného klienta.
    ```bash
    go run cmd/getdata/main.go list_sites 123
    go run cmd/getdata/main.go list_sites "Jméno Klienta"
    ```

*   **`list_servers`**: Vypíše všechny servery pro danou site.
    ```bash
    go run cmd/getdata/main.go list_servers 456
    go run cmd/getdata/main.go list_servers "Jméno Site"
    ```

*   **`list_workstations`**: Vypíše všechny pracovní stanice pro danou site.
    ```bash
    go run cmd/getdata/main.go list_workstations 456
    go run cmd/getdata/main.go list_workstations "Jméno Site"
    ```

*   **`list_devices`**: Vypíše všechna zařízení pro danou site.
    ```bash
    go run cmd/getdata/main.go list_devices 456
    ```

*   **`list_devices_at_client`**: Vypíše všechna zařízení pro daného klienta.
    ```bash
    go run cmd/getdata/main.go list_devices_at_client 123
    ```

*   **`list_agentless_assets`**: Vypíše agentless assets pro danou site.
    ```bash
    go run cmd/getdata/main.go list_agentless_assets 456
    ```

#### Monitorování a kontroly:

*   **`list_failing_checks`**: Vypíše všechny neúspěšné kontroly.
    ```bash
    go run cmd/getdata/main.go list_failing_checks
    ```

*   **`list_checks`**: Vypíše kontroly pro zařízení nebo site.
    ```bash
    go run cmd/getdata/main.go list_checks 789
    ```

*   **`list_device_monitoring_details`**: Vypíše podrobnosti monitorování zařízení.
    ```bash
    go run cmd/getdata/main.go list_device_monitoring_details 789
    ```

*   **`list_check_configuration`**: Vypíše konfiguraci kontrol.
    ```bash
    go run cmd/getdata/main.go list_check_configuration 789
    go run cmd/getdata/main.go list_check_configuration_windows 789
    go run cmd/getdata/main.go list_check_configuration_mac 789
    go run cmd/getdata/main.go list_check_configuration_linux 789
    ```

*   **`list_outages`**: Vypíše výpadky systému.
    ```bash
    go run cmd/getdata/main.go list_outages 456 "2024-01-01" "2024-01-31"
    ```

*   **`clear_check`**: Vymaže specifickou kontrolu.
    ```bash
    go run cmd/getdata/main.go clear_check 12345
    ```

*   **`add_check_note`**: Přidá poznámku ke kontrole.
    ```bash
    go run cmd/getdata/main.go add_check_note 12345 "Poznámka k této kontrole"
    ```

#### Správa asset tracking:

*   **`list_hardware`**: Vypíše hardware informace pro zařízení.
    ```bash
    go run cmd/getdata/main.go list_hardware 789
    ```

*   **`list_software`**: Vypíše software informace pro zařízení.
    ```bash
    go run cmd/getdata/main.go list_software 789
    ```

*   **`list_device_asset_details`**: Vypíše podrobné asset informace zařízení.
    ```bash
    go run cmd/getdata/main.go list_device_asset_details 789
    ```

*   **`list_license_groups`**: Vypíše licenční skupiny.
    ```bash
    go run cmd/getdata/main.go list_license_groups
    ```

#### Správa patchů:

*   **`list_patches`**: Vypíše všechny patche pro zařízení.
    ```bash
    go run cmd/getdata/main.go list_patches 789
    ```

*   **`approve_patch`**: Schválí patche pro zařízení.
    ```bash
    go run cmd/getdata/main.go approve_patch 789 "12345,12346,12347"
    ```

*   **`ignore_patch`**: Ignoruje patche pro zařízení.
    ```bash
    go run cmd/getdata/main.go ignore_patch 789 "12345,12346"
    ```

#### Antivirus:

*   **`list_antivirus_products`**: Vypíše podporované antivirus produkty.
    ```bash
    go run cmd/getdata/main.go list_antivirus_products
    ```

*   **`list_antivirus_definitions`**: Vypíše antivirus definice.
    ```bash
    go run cmd/getdata/main.go list_antivirus_definitions 789
    ```

*   **`list_quarantine`**: Vypíše položky v karanténě.
    ```bash
    go run cmd/getdata/main.go list_quarantine 789
    ```

*   **`start_scan`**: Spustí antivirus scan.
    ```bash
    go run cmd/getdata/main.go start_scan 789 "full"
    ```

#### Výkon a historie:

*   **`list_performance_history`**: Vypíše historii výkonu.
    ```bash
    go run cmd/getdata/main.go list_performance_history 789 12345 "2024-01-01" "2024-01-31"
    ```

*   **`list_drive_space_history`**: Vypíše historii využití disků.
    ```bash
    go run cmd/getdata/main.go list_drive_space_history 789 "2024-01-01" "2024-01-31"
    ```

#### Šablony:

*   **`list_templates`**: Vypíše monitorovací šablony.
    ```bash
    go run cmd/getdata/main.go list_templates
    ```

#### Backup & Recovery:

*   **`list_backup_sessions`**: Vypíše backup session.
    ```bash
    go run cmd/getdata/main.go list_backup_sessions 789
    ```

#### Nastavení:

*   **`list_wall_chart_settings`**: Vypíše nastavení wall chart.
    ```bash
    go run cmd/getdata/main.go list_wall_chart_settings
    ```

*   **`list_general_settings`**: Vypíše obecná nastavení.
    ```bash
    go run cmd/getdata/main.go list_general_settings
    ```

#### Úlohy a uživatelé:

*   **`list_active_directory_users`**: Vypíše Active Directory uživatele.
    ```bash
    go run cmd/getdata/main.go list_active_directory_users 789
    ```

*   **`run_task_now`**: Spustí úlohu okamžitě.
    ```bash
    go run cmd/getdata/main.go run_task_now 12345
    ```

#### Správa site:

*   **`add_client`**: Přidá nového klienta.
    ```bash
    go run cmd/getdata/main.go add_client "Název klienta" "Kontakt" "email@example.com"
    ```

*   **`add_site`**: Přidá novou site ke klientovi.
    ```bash
    go run cmd/getdata/main.go add_site 123 "Název site" "Kontakt" "email@example.com"
    ```

*   **`get_site_installation_package`**: Získá instalační balíček pro site.
    ```bash
    go run cmd/getdata/main.go get_site_installation_package 456 "windows"
    ```

### 2. `fetchall`

Tento nástroj stáhne komplexní data o všech klientech, jejich sites a zařízeních (servery, stanice). Data uloží do CSV souborů v adresáři `data/` (slouží jako cache) a zároveň vypíše kompletní vnořenou strukturu jako JSON.

**Použití:**

```bash
go run cmd/fetchall/main.go [-cache] [vystupni_soubor.json]
```

**Argumenty:**

*   `-cache` (volitelný): Pokud je tento příznak uveden, nástroj **nevolá N-Sight API**, ale místo toho načte data z existujících CSV souborů v adresáři `data/` a sestaví z nich JSON výstup. Vyžaduje, aby CSV soubory již existovaly (tj. aby byl `fetchall` spuštěn alespoň jednou bez `-cache`).
*   `[vystupni_soubor.json]` (volitelný): Pokud je zadán název souboru, výsledný JSON se zapíše do tohoto souboru. Pokud není zadán, JSON se vypíše na standardní výstup.

**Příklady:**

*   **Načíst z API, zapsat CSV, vypsat JSON na obrazovku:**
    ```bash
    go run cmd/fetchall/main.go
    ```
*   **Načíst z API, zapsat CSV, zapsat JSON do `komplet.json`:**
    ```bash
    go run cmd/fetchall/main.go komplet.json
    ```
*   **Načíst z CSV cache, vypsat JSON na obrazovku:**
    ```bash
    go run cmd/fetchall/main.go -cache
    ```
*   **Načíst z CSV cache, zapsat JSON do `cache_data.json`:**
    ```bash
    go run cmd/fetchall/main.go -cache cache_data.json
    ```

**CSV Cache:**

*   Nástroj `fetchall` (v režimu bez `-cache`) vytváří následující CSV soubory v adresáři `data/`:
    *   `clients.csv`
    *   `sites.csv`
    *   `servers.csv`
    *   `workstations.csv`
*   Tento adresář je zahrnut v `.gitignore`, takže cache soubory nebudou součástí Gitu.

**JSON Výstup (`fetchall`):**

Výstupem je pole objektů, kde každý objekt reprezentuje klienta a obsahuje vnořený seznam jeho sites, a každá site obsahuje vnořené seznamy serverů a pracovních stanic.

```json
[
  {
    "client_id": 123,
    "client_name": "Klient A",
    "sites": [
      {
        "site_id": 456,
        "site_name": "Site X",
        "servers": [
          {
            "server_id": 789,
            "server_name": "Server Y",
            "online": true,
            "os": "...",
            "ip": "..."
          }
        ],
        "workstations": [
          {
            "workstation_id": 101,
            "workstation_name": "Stanice Z",
            "online": false,
            "os": "...",
            "ip": "..."
          }
        ]
      }
      // ... dalsi sites pro Klienta A
    ]
  }
  // ... dalsi klienti
]
```

### 3. `nsight-proxy` - JSON API Proxy Server

HTTP proxy server, který převádí N-Sight XML API na JSON formát. Poslouchá na portu 80 a poskytuje stejné API volání jako originální N-Sight, ale s JSON výstupem místo XML.

**Spuštění:**

```bash
go run cmd/nsight-proxy/main.go
```

**Použití:**

```bash
# Získání seznamu klientů ve formátu JSON
curl "http://localhost/api/?service=list_clients"

# Získání serverů pro site 123
curl "http://localhost/api/?service=list_servers&siteid=123"

# Health check
curl "http://localhost/health"
```

Proxy server podporuje všechna API volání stejně jako nástroj `getdata`, ale poskytuje je přes HTTP rozhraní s JSON výstupem. Více informací v [dokumentaci proxy serveru](cmd/nsight-proxy/README.md).

## API Server (`cmd/server`)

Adresář `cmd/server/main.go` obsahuje základ pro budoucí REST API server. V současné době pouze spustí jednoduchý HTTP server na portu 8080 a na kořenové cestě `/` vrací uvítací zprávu.

**Spuštění:**

```bash
go run cmd/server/main.go
```

Tento server bude v budoucnu rozšířen o endpointy, které budou využívat buď přímá volání API (přes `internal/nsight`) nebo data načtená z CSV cache (vytvořené nástrojem `fetchall`).

## Podporovaná API volání

Nástroj `getdata` nyní podporuje všechna dostupná N-Sight API volání podle oficiální dokumentace na https://developer.n-able.com/n-sight/docs/getting-started-with-the-n-sight-api, včetně:

- **Základní entity**: klienti, sites, servery, workstations, zařízení, agentless assets
- **Monitorování**: failing checks, device monitoring, check configuration, outages
- **Asset tracking**: hardware, software, license groups, device asset details
- **Patch management**: list, approve, ignore patches
- **Antivirus**: produkty, definice, karanténa, spuštění scanů
- **Výkon a historie**: performance history, drive space history
- **Šablony**: monitoring templates
- **Backup & Recovery**: backup sessions
- **Nastavení**: wall chart settings, general settings
- **Úlohy a AD**: Active Directory users, task execution
- **Site management**: přidání klientů/sites, installation packages

Celkem je podporováno **37 různých API volání**, což pokrývá kompletní funkcionalitu N-Sight Data Extraction API.

## Vylepšení oproti původní verzi

1. **Kompletní pokrytí API**: Podporuje všechna dostupná NSight API volání místo pouze 4 základních.
2. **Lepší organizace kódu**: Čistě oddělené handlery pro každé API volání.
3. **Rozšířené typy dat**: Podporuje všechny datové struktury podle API dokumentace.
4. **Flexibilní parametry**: Inteligentní resolving identifikátorů (ID vs. jména).
5. **Konzistentní error handling**: Jednotný přístup k chybám napříč všemi voláními.
6. **Akční API volání**: Kromě čtení dat podporuje i akce jako schvalování patchů, spouštění scanů atd.
7. **Podrobná nápověda**: Kompletní usage informace pro všechna API volání. 