package quran

import "fmt"

// PrintTableSamples بيطبع أول 5 صفوف من جدولي السور والآيات
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
		err := verseRows.Scan(&id, &textSimple, &tafseerMuasr, &hizbText, &juzText, &suraNum)
		if err != nil {
			return fmt.Errorf("error scanning verse row: %w", err)
		}

		// قص النصوص الطويلة عشان شكل الطباعة في الـ Terminal ميبوظش
		shortText := limitString(textSimple, 25)
		shortTafseer := limitString(tafseerMuasr, 25)

		fmt.Printf("%-5d | %-5d | %-10s | %-10s | %-30s | %-30s\n", 
			id, suraNum, juzText, hizbText, shortText, shortTafseer)
	}

	return nil
}

// فانكشن مساعدة لقص النصوص الطويلة
func limitString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes [:maxLen-3]) + "..."
}