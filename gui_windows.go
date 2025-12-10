//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/tadvi/winc"
)

func main() {
	createMainWindow()
	winc.RunMainLoop()
}

func createMainWindow() {
	mainWindow := winc.NewForm(nil)
	mainWindow.SetText("Tank Report")
	mainWindow.SetSize(700, 200)
	mainWindow.Center()

	// Label percorso
	lblDir := winc.NewLabel(mainWindow)
	lblDir.SetPos(20, 20)
	lblDir.SetSize(120, 20)
	lblDir.SetText("Cartella riscaldi:")

	// TextBox percorso
	txtDir := winc.NewEdit(mainWindow)
	txtDir.SetPos(150, 18)
	txtDir.SetSize(320, 24)

	// Default: sottocartella "Riscaldi" accanto all'exe
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	defaultDir := filepath.Join(exeDir, "Riscaldi")
		txtDir.SetText(defaultDir)

		// Bottone "Sfoglia..."
		btnBrowse := winc.NewPushButton(mainWindow)
		btnBrowse.SetText("Sfoglia...")
		btnBrowse.SetPos(480, 16)
		btnBrowse.SetSize(80, 28)

		// Bottone "Genera report"
		btnGenerate := winc.NewPushButton(mainWindow)
		btnGenerate.SetText("Genera report")
		btnGenerate.SetPos(570, 16)
		btnGenerate.SetSize(100, 28)

		// Log (una sola riga di stato)
		txtLog := winc.NewEdit(mainWindow)
		txtLog.SetPos(20, 60)
		txtLog.SetSize(650, 70)
		txtLog.SetReadOnly(true)

		// Azione bottone "Sfoglia..."
		btnBrowse.OnClick().Bind(func(e *winc.Event) {
			folder, ok := winc.ShowBrowseFolderDlg(mainWindow, "Seleziona cartella riscaldi")
			if ok && folder != "" {
				txtDir.SetText(folder)
			}
		})

		// Azione bottone "Genera report"
		btnGenerate.OnClick().Bind(func(e *winc.Event) {
			selectedDir := txtDir.Text()
			if selectedDir == "" {
				txtLog.SetText("Nessuna cartella specificata.")
				return
			}

			txtLog.SetText("Avvio elaborazione su: " + selectedDir)

			recordsByMonth, err := ProcessDir(selectedDir)
			if err != nil {
				txtLog.SetText("Errore nella lettura cartella: " + err.Error())
				return
			}

			if len(recordsByMonth) == 0 {
				txtLog.SetText("Nessun record trovato nei file CSV.")
				return
			}

			exePath, _ := os.Executable()
			exeDir := filepath.Dir(exePath)

			for monthKey, records := range recordsByMonth {
				outputFile := fmt.Sprintf("tank_report_%s.csv", monthKey)
				outPath := filepath.Join(exeDir, outputFile)
				if err := GenerateCSVForMonth(outPath, records); err != nil {
					txtLog.SetText("Errore creazione " + outputFile + ": " + err.Error())
					return
				}
			}

			txtLog.SetText("Operazione completata con successo.")
		})

		mainWindow.OnClose().Bind(func(e *winc.Event) {
			winc.Exit()
		})

		mainWindow.Show()
}
