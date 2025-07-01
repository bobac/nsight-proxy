# Použití Nástrojů NSight Proxy (bez zdrojového kódu)

Tento dokument popisuje, jak používat zkompilované nástroje `getdata`, `fetchall` a `nsight-proxy`.

## Požadavky

*   Spustitelné soubory `getdata`, `fetchall` a `nsight-proxy` pro váš operační systém.
*   Soubor `.env` umístěný ve stejném adresáři jako spustitelné soubory.

## Nastavení

1.  **Umístěte spustitelné soubory**: Zkopírujte soubory `getdata`, `fetchall` a `nsight-proxy` do adresáře, kde chcete projekt používat.
2.  **Vytvořte soubor `.env`**: Ve stejném adresáři vytvořte soubor s názvem `.env` (můžete zkopírovat `.env.example`, pokud ho máte k dispozici).
3.  **Vyplňte `.env`**: Otevřete soubor `.env` a zadejte platné hodnoty pro proměnné:
    *   `NSIGHT_API_KEY=VAS_API_KLIC`
    *   `NSIGHT_SERVER=HOSTNAME_SERVERU` (např. `wwweurope1.systemmonitor.eu.com`, bez `https://`)

## Dostupné Nástroje

Následující příkazy předpokládají, že jste v terminálu/příkazovém řádku ve stejném adresáři, kde jsou umístěny spustitelné soubory a soubor `.env`.

### 1. `getdata`

Získá specifická data z N-Sight API a vypíše je jako JSON na standardní výstup.

**Syntaxe:**

```bash
./getdata <název_služby> [parametry...]
```
*(Na Windows použijte `.\getdata.exe`)*

**Podporované služby:**

*   **`list_clients`**: Vypíše všechny klienty.
    ```bash
    ./getdata list_clients
    ```

*   **`list_sites`**: Vypíše sites pro klienta (podle ID nebo jména).
    ```bash
    ./getdata list_sites 123
    ./getdata list_sites "Jméno Klienta"
    ```

*   **`list_servers`**: Vypíše servery pro site (podle ID nebo jména).
    ```bash
    ./getdata list_servers 456
    ./getdata list_servers "Jméno Site"
    ```

*   **`list_workstations`**: Vypíše pracovní stanice pro site (podle ID nebo jména).
    ```bash
    ./getdata list_workstations 456
    ./getdata list_workstations "Jméno Site"
    ```

### 2. `fetchall`

Stáhne data o všech klientech, sites a zařízeních. Vytvoří/aktualizuje CSV soubory v podadresáři `data/` a vypíše kompletní vnořenou strukturu jako JSON.

**Syntaxe:**

```bash
./fetchall [-cache] [vystupni_soubor.json]
```
*(Na Windows použijte `.\fetchall.exe`)*

**Argumenty:**

*   `-cache` (volitelný): Načte data z existujících CSV souborů v adresáři `data/` místo volání API. Pokud adresář `data/` nebo potřebné CSV soubory neexistují, skončí chybou.
*   `[vystupni_soubor.json]` (volitelný): Zapíše výsledný JSON do tohoto souboru místo výpisu na obrazovku.

**Příklady:**

*   **Načíst z API, vytvořit/aktualizovat CSV, vypsat JSON:**
    ```bash
    ./fetchall
    ```
*   **Načíst z API, vytvořit/aktualizovat CSV, zapsat JSON do souboru:**
    ```bash
    ./fetchall komplet.json
    ```
*   **Načíst z existující CSV cache, vypsat JSON:**
    ```bash
    ./fetchall -cache
    ```
*   **Načíst z existující CSV cache, zapsat JSON do souboru:**
    ```bash
    ./fetchall -cache cache_data.json
    ```

### 3. `nsight-proxy`

Spustí HTTP proxy server, který převádí N-Sight XML API na JSON formát.

**Syntaxe:**

```bash
./nsight-proxy
```
*(Na Windows použijte `.\nsight-proxy.exe`)*

Server běží na popředí a naslouchá na portu 80. Ukončíte ho stiskem `Ctrl+C`. Server potřebuje přístup k souboru `.env` pro konfiguraci N-Sight serveru. Více informací o použití proxy serveru najdete v dokumentaci. 