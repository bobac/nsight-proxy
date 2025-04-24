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

Projekt obsahuje dva hlavní nástroje v adresáři `cmd/`:

### 1. `getdata`

Tento nástroj slouží k přímému volání specifických N-Sight API služeb a vrací výsledek jako JSON na standardní výstup. Je užitečný pro rychlé získání konkrétních informací.

**Použití:**

```bash
go run cmd/getdata/main.go <název_služby> [parametry...]
```

**Podporované služby:**

*   **`list_clients`**: Vypíše všechny klienty.
    ```bash
    go run cmd/getdata/main.go list_clients
    ```
    *Výstup:* `[{"client_id": 123, "client_name": "Klient A"}, ...]`

*   **`list_sites`**: Vypíše všechny sites pro daného klienta. Klienta lze specifikovat pomocí ID nebo jména (v uvozovkách, pokud obsahuje mezery).
    ```bash
    # Podle ID klienta
    go run cmd/getdata/main.go list_sites 123

    # Podle jména klienta
    go run cmd/getdata/main.go list_sites "Jméno Klienta"
    ```
    *Výstup:* `[{"site_id": 456, "site_name": "Site X"}, ...]`

*   **`list_servers`**: Vypíše všechny servery pro danou site. Site lze specifikovat pomocí ID nebo jména. Pokud je zadáno jméno, nástroj prohledá všechny klienty, aby našel odpovídající site ID.
    ```bash
    # Podle ID site
    go run cmd/getdata/main.go list_servers 456

    # Podle jména site (může trvat déle)
    go run cmd/getdata/main.go list_servers "Jméno Site"
    ```
    *Výstup:* `[{"server_id": 789, "server_name": "Server Y", "online": true, ...}, ...]`

*   **`list_workstations`**: Vypíše všechny pracovní stanice pro danou site. Site lze specifikovat pomocí ID nebo jména (s vyhledáváním napříč klienty jako u `list_servers`).
    ```bash
    # Podle ID site
    go run cmd/getdata/main.go list_workstations 456

    # Podle jména site (může trvat déle)
    go run cmd/getdata/main.go list_workstations "Jméno Site"
    ```
    *Výstup:* `[{"workstation_id": 101, "workstation_name": "Stanice Z", "online": false, ...}, ...]`

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

## API Server (`cmd/server`)

Adresář `cmd/server/main.go` obsahuje základ pro budoucí REST API server. V současné době pouze spustí jednoduchý HTTP server na portu 8080 a na kořenové cestě `/` vrací uvítací zprávu.

**Spuštění:**

```bash
go run cmd/server/main.go
```

Tento server bude v budoucnu rozšířen o endpointy, které budou využívat buď přímá volání API (přes `internal/nsight`) nebo data načtená z CSV cache (vytvořené nástrojem `fetchall`). 