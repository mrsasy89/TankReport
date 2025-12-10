package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
	"strings"
)

type Record struct {
	IDRecord           string
	Cliente            string
	Operatore          string
	NumeroTank         string
	TargaRimorchio     string
	NumeroPista        int
	TempDaRaggiungere  float64
	DataInizio         time.Time
	OraInizio          string
	TempIniziale       float64
	DataFine           time.Time
	OraFine            string
	TempFinale         float64
	MinutiTotali       int
	MinutiValvola      int
}

// ProcessDir legge una cartella di file grezzi e restituisce i record raggruppati per mese (YYYY-MM)
func ProcessDir(dir string) (map[string][]Record, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	recordsByMonth := make(map[string][]Record)

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			filePath := filepath.Join(dir, file.Name())
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("❌ Errore apertura %s: %v\n", file.Name(), err)
				continue
			}

			reader := csv.NewReader(f)
			reader.Comma = ';'
			lines, err := reader.ReadAll()
			f.Close()

			if err != nil || len(lines) == 0 {
				fmt.Printf("❌ Errore lettura %s\n", file.Name())
				continue
			}

			rec, err := parseLine(lines[0])
			if err != nil {
				fmt.Printf("❌ Errore parsing %s: %v\n", file.Name(), err)
				continue
			}

			monthKey := rec.DataInizio.Format("2006-01")
			recordsByMonth[monthKey] = append(recordsByMonth[monthKey], rec)
		}
	}

	// ordina internamente per ogni mese
	for monthKey, records := range recordsByMonth {
		sort.Slice(records, func(i, j int) bool {
			return records[i].DataInizio.Before(records[j].DataInizio)
		})
		recordsByMonth[monthKey] = records
	}

	return recordsByMonth, nil
}

// GenerateCSVForMonth usa generateCSV esistente per scrivere il file finale
func GenerateCSVForMonth(filename string, records []Record) error {
	return generateCSV(filename, records)
}

func parseLine(line []string) (Record, error) {
	if len(line) < 19 {
		return Record{}, fmt.Errorf("riga troppo corta: %d campi", len(line))
	}

	pista, _ := strconv.Atoi(line[6])
	tempTarget, _ := strconv.ParseFloat(line[7], 64)
	tempInit, _ := strconv.ParseFloat(line[10], 64)

	// Cerca temp finale DOPO campo 13, saltando vuoti
	tempFinale := float64(0)
	for i := 13; i < len(line); i++ {
		if f, err := strconv.ParseFloat(line[i], 64); err == nil && f > 0 {
			tempFinale = f
			break
		}
	}

	minutiTotali, _ := strconv.Atoi(line[17])
	minutiValvola, _ := strconv.Atoi(line[18])

	// Data inizio (sempre presente)
	dataInizioStr := line[8] + " " + line[9]
	dataInizio, err := time.Parse("2006-01-02 15:04:05", dataInizioStr)
	if err != nil {
		return Record{}, fmt.Errorf("data inizio %s: %v", dataInizioStr, err)
	}

	// Data/Ora fine: cerca PRIMO formato valido DOPO campo 11
	oraFine := ""
	dataFine := dataInizio // Fallback
	for i := 11; i < len(line)-4; i++ {
		if line[i] != "" && strings.HasPrefix(line[i], "20") { // Data valida
			dataFineStr := line[i] + " " + line[i+1]
			if t, err := time.Parse("2006-01-02 15:04:05", dataFineStr); err == nil {
				dataFine = t
				oraFine = line[i+1]
				break
			}
		}
	}

	return Record{
		IDRecord:           line[0],
		Cliente:            line[2],
		Operatore:          line[3],
		NumeroTank:         line[4],
		TargaRimorchio:     line[5],
		NumeroPista:        pista,
		TempDaRaggiungere:  tempTarget,
		DataInizio:         dataInizio,
		OraInizio:          line[9],
		TempIniziale:       tempInit,
		DataFine:           dataFine,
		OraFine:            oraFine,
		TempFinale:         tempFinale,
		MinutiTotali:       minutiTotali,
		MinutiValvola:      minutiValvola,
	}, nil
}

func generateCSV(filename string, records []Record) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.Comma = ';'

	// Header SOLO con le colonne richieste
	header := []string{
		"Numero Tank",
		"Temperatura iniziale (°C)",
		"Data Inizio Riscaldo",
		"Ora Inizio Riscaldo",
		"Temperatura finale (°C)",
		"Data fine riscaldo",
		"Ora fine riscaldo",
		"Totale ore",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Somma totale minuti per la riga finale
	totalMinutes := 0

	for _, rec := range records {
		// calcolo ore singolo record dai minuti totali
		ore := rec.MinutiTotali / 60
		minuti := rec.MinutiTotali % 60
		totaleOreStr := fmt.Sprintf("%02d:%02d", ore, minuti)

		row := []string{
			rec.NumeroTank,
			fmt.Sprintf("%.0f", rec.TempIniziale),
			rec.DataInizio.Format("02/01/2006"), // come file operatore
			rec.OraInizio[:5],                  // HH:MM
			fmt.Sprintf("%.0f", rec.TempFinale),
			rec.DataFine.Format("02/01/2006"),
			func() string {
				if rec.OraFine == "" {
					return ""
				}
				if len(rec.OraFine) >= 5 {
					return rec.OraFine[:5]
				}
				return rec.OraFine
			}(),
			totaleOreStr,
		}
		if err := writer.Write(row); err != nil {
			return err
		}

		totalMinutes += rec.MinutiTotali
	}

	// Riga vuota
	if err := writer.Write([]string{}); err != nil {
		return err
	}

	// Riga TOTALE ORE
	totalHours := totalMinutes / 60
	totalMins := totalMinutes % 60
	totalStr := fmt.Sprintf("%02d:%02d", totalHours, totalMins)

	totalRow := []string{
		"", "", "", "", "", "", "Totale ore", totalStr,
	}
	if err := writer.Write(totalRow); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}
