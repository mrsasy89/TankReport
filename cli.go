//go:build !windows

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Uso: tank-report <cartella_riscaldi>")
	}

	dir := os.Args[1]

	recordsByMonth, err := ProcessDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for monthKey, records := range recordsByMonth {
		outputFile := fmt.Sprintf("tank_report_%s.csv", monthKey)
		if err := GenerateCSVForMonth(outputFile, records); err != nil {
			fmt.Printf("❌ Errore creazione %s: %v\n", outputFile, err)
			continue
		}
		fmt.Printf("✅ Generato %s: %d record (%s)\n",
			   outputFile, len(records), monthKey)
	}
}
