# Tank Report

**Generatore automatico di report mensili dai file CSV grezzi estratti da un PCL**

Analizza i log CSV della cartella "Riscaldi" e genera report aggregati per mese (tank_report_YYYY-MM.csv).

## Funzionalità

- **Multi-piattaforma**: Linux e Windows nativo (no console su Windows)
- **GUI intuitiva** (solo Windows): 
  - Cartella predefinita `./Riscaldi`
  - Bottone "Sfoglia..." per scegliere altra cartella
  - Log di stato in tempo reale
- **Elaborazione veloce**: legge tutti i CSV, raggruppa per mese, genera report
- **Output**: CSV per ogni mese con statistiche aggregate

## Utilizzo

### Windows (con GUI)
1. Avvia `TankReport.exe`
2. Cartella predefinita: `./Riscaldi` (creala se non esiste)
3. Clicca **"Sfoglia..."** per scegliere altra cartella
4. Clicca **"Genera report"**
5. Report salvati come `tank_report_YYYY-MM.csv` accanto all'exe

### Linux (CLI)
`./tank-report ./Riscaldi` 

## Struttura cartella "Riscaldi"
Riscaldi/
├── 2025-01-15_log.csv
├── 2025-01-16_log.csv
├── 2025-02-01_log.csv
└── ...


## Compilazione
### Linux
`go build -o tank-report .`

### Windows (cross-compile da Linux)
`GOOS=windows GOARCH=amd64 go build -o TankReport.exe -ldflags="-H windowsgui"`

## Dipendenze

go mod tidy

**Solo winc per GUI Windows** (nessuna dipendenza su Linux).

## Personalizzazioni

- **Cartella output**: modifica `exeDir` in `gui_windows.go`
- **Formato data**: adatta `parseDate` in `process.go`
- **Campi CSV**: personalizza `ProcessDir` e `GenerateCSVForMonth`
