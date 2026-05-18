package main

import (
	"fmt"
	"log"

	"github.com/MohammadAboulEla/quran-db-go/quran"
)

func main() {
	qm, err := quran.NewQuranManager("quran/quran.db")
	if err != nil {
		log.Fatal(err)
	}
	defer qm.Close()

	// Get a single aya by global number (1..6236)
	aya, _ := qm.GetAyaByGlobal(1)
	fmt.Println(aya) // [1:1] بسم الله الرحمن الرحيم

	// Get an aya by sura + local position
	aya, _ = qm.GetAyaByLocal(2, 255) // Ayat al-Kursi
	fmt.Println(aya.Text)

	// Whole sura with all ayas
	fatiha, _ := qm.GetSuraWithAyas(1)
	fmt.Printf("%s — %d ayas\n", fatiha.Type, fatiha.TotalAyasNumbers)

	// Search
	hits, _ := qm.SearchAyas("الرحمن")
	fmt.Printf("%d matches\n", len(hits))
}
