package quran

import "fmt"

// PrintTableSamples prints the first 5 rows of the suras and verses tables.
// Useful as a quick sanity check that the DB is reachable and schema matches.
func (qm *QuranManager) PrintTableSamples() error {
	fmt.Println("================================================")
	fmt.Println("SURAS TABLE (FIRST 5 ROWS)")
	fmt.Println("================================================")

	suraRows, err := qm.db.Query("SELECT id, name, type, total_verses FROM suras LIMIT 5")
	if err != nil {
		return fmt.Errorf("error querying suras: %w", err)
	}
	defer suraRows.Close()

	fmt.Printf("%-5s | %-15s | %-10s | %-12s\n", "ID", "Name", "Type", "Total Verses")
	fmt.Println("---------------")
	for suraRows.Next() {
		var id, totalVerses int
		var name, suraType string
		if err := suraRows.Scan(&id, &name, &suraType, &totalVerses); err != nil {
			return fmt.Errorf("error scanning sura row: %w", err)
		}
		fmt.Printf("%-5d | %-15s | %-10s | %-12d\n", id, name, suraType, totalVerses)
	}

	fmt.Println("\n==========================================")
	fmt.Println("VERSES TABLE (FIRST 5 ROWS)")
	fmt.Println("============================================")

	query := `SELECT id, text_simple, tafseer_muasr, hizb_text, juz_text, sura_num FROM verses LIMIT 5`
	verseRows, err := qm.db.Query(query)
	if err != nil {
		return fmt.Errorf("error querying verses: %w", err)
	}
	defer verseRows.Close()

	fmt.Printf("%-5s | %-5s | %-10s | %-10s | %-30s | %-30s\n", "ID", "Sura", "Juz", "Hizb", "Text Simple (Sample)", "Tafseer Muasr (Sample)")
	fmt.Println("----------------------------")
	for verseRows.Next() {
		var id, suraNum int
		var textSimple, tafseerMuasr, hizbText, juzText string
		if err := verseRows.Scan(&id, &textSimple, &tafseerMuasr, &hizbText, &juzText, &suraNum); err != nil {
			return fmt.Errorf("error scanning verse row: %w", err)
		}
		// Truncate long Arabic text so the terminal layout stays aligned.
		shortText := limitString(textSimple, 25)
		shortTafseer := limitString(tafseerMuasr, 25)
		fmt.Printf("%-5d | %-5d | %-10s | %-10s | %-30s | %-30s\n",
			id, suraNum, juzText, hizbText, shortText, shortTafseer)
	}

	return nil
}

// limitString truncates s to at most maxLen runes (not bytes), adding "..." if cut.
func limitString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
