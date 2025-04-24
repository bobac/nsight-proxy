# NSight Proxy

## Účel projektu
Účelem projektu je vyrobit zástupné API pro N-Able N-Sight. N-Sight používá tzv. Data extraction API (https://developer.n-able.com/n-sight/docs/getting-started-with-the-n-sight-api), které:

1. Vrací XML místo JSON, který potřebuji. Prvním cílem je tedy implementovat obdobné API endpointy, ale tak, aby z nich data chodila jako JSON. Navíc by konvence měly odpovídat současnému REST API.
2. Většinu dat není možné získat jinak, než postupným voláním několika endpointů, například pokud chci dostat všechen nainstalovaný software, musím postupně iterovat přes LIST_CLIENTS, LIST_SITES, LIST_SERVERS a LIST_WORKSTATIONS a pak pro každé zařízení volat LIST_DEVICE_ASSET_DETAILS. Idea je vytvořit API volání, které postupně načte všechna data do cache (stačí rozsáhlý CSV soubor) a pak volat souhrnné požadavky, které vezmou date z této cache a vrátí je jako strukturovaný json. Tedy např. výše uvedený příklad by vrátil strukturu (pseudo-json) jako:
{
    [client1,
        [site1,
            [server1,
            {software,
            [sw1,
            sw2]]
            }
        site2,]
    slient2,
    ...]
}

## Implementace
Implementovat v jazyce GO, tak, aby projekt byl zároveň modulem, který lze volat z jiných projektů. v `cmd/server/main.go` implementovat API endpoint, v `cmd/<tool>/main.go` implemntovat nástroje. Například bych rád implementoval `cmd/getdata/main.go`, který by z příkazového řádku zavolal N-Sight Data Extraction API a na standardní výstup vrátil odpovídající reprezentaci dat v JSON. Další tool je `cmd/fetchall/main.go`, který do adresáře data/ stáhne všechna relevantní data jako velké CVS, které bude sloužit jako cache pro agregované volání z bodu 2 výše.

API key a server pro vytvoření URL pro data extraction api bude v .env souboru. Obecný formát volání N-Sight API je `https://{server}/api/?apikey={yourAPIkey}&service={service_name}&parameter={parameter-value}`

## Umístění dat
Funkce, které mapují stávající (dle bodu 1) API 1:1 data nikam neukládají - při volání REST API pouze zavolají odpovídající Data Extarction API N-Sight. Agregované funkce (dle bodu 2) používají jako zdroj dat CSV soubor v adresáři data/
