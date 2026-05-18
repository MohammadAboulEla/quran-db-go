# quran-db

A small Go library for querying a SQLite database of the Holy Quran (114 suras, 6236 ayas) with two Arabic tafseers (Muyassar and Sa'di), plus juz and hizb metadata.

Uses the pure-Go [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) driver — **no CGo, no system SQLite required**.

## Install

```bash
git clone https://github.com/<you>/quran-db
cd quran-db
go build ./...
```

The database file lives at [`quran/quran.db`](quran/quran.db) and ships with the repo.

## Quick start

```go
package main

import (
    "fmt"
    "log"

    . "quran-db/quran"
)

func main() {
    qm, err := NewQuranManager("quran/quran.db")
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
```

## API

### Construction

| Function | Description |
| --- | --- |
| `NewQuranManager(dbPath string)` | Open the SQLite DB (sets WAL + 64MB cache pragmas). |
| `Close()` | Close the underlying connection. |
| `GetQuranInfo()` | Constants: 6236 verses, 114 suras, 30 juzs, 60 hizbs. |

### Suras

| Function | Description |
| --- | --- |
| `GetSura(n)` | Metadata only (no ayas). |
| `GetSuraWithAyas(n)` | Metadata + full `AyasList`. |
| `GetAllSuras()` | All 114, metadata only. |
| `GetSurasByType(Meccan \| Medinan)` | Filter by revelation place. |
| `GetSuraName(n)` / `GetSuraType(n)` / `GetSuraTotalVerses(n)` | Single-column lookups. |
| `GetSuraNames()` | All 114 names in order. |
| `GetSuraText(n, sep)` | Full Uthmani text, ayas joined by `sep`. |

### Ayas

| Function | Description |
| --- | --- |
| `GetAyaByGlobal(n)` | By global id (1..6236). |
| `GetAyaByLocal(sura, local)` | By position within a sura. |
| `GetAyasByGlobalRange(from, to)` | Inclusive range. |
| `GetAyasByJuz(arabicLabel)` / `GetAyasByJuzNumber(1..30)` | All ayas in a juz. |
| `GetAyasByHizb(arabicLabel)` | All ayas in a hizb / quarter-hizb. |
| `SearchAyas(query)` | `LIKE %query%` over `text_simple`. |
| `GlobalToLocal(n)` / `LocalToGlobal(sura, local)` | Coordinate conversions. |

### Tafseer

| Function | Description |
| --- | --- |
| `GetTafseerMuasr(global)` | Modern (Muyassar) tafseer for an aya. |
| `GetTafseerSadi(global)` | Sa'di tafseer for an aya. |

## Schema

```
suras(id INTEGER PK, name TEXT, type TEXT, total_verses INTEGER)
verses(id INTEGER PK,
       text TEXT,           -- Uthmani script with diacritics
       text_simple TEXT,    -- simplified for search
       tafseer_muasr TEXT,
       tafseer_sadi TEXT,
       hizb_text TEXT,      -- Arabic label, e.g. "ربع الحزب الـ 1"
       juz_text TEXT,       -- Arabic label, e.g. "الجزء الأول"
       sura_num INTEGER)
```

- `verses.id` is the **global** ayah number (1..6236).
- `suras.type` is `"مكية"` (Meccan) or `"مدنية"` (Medinan).
- `juz_text` / `hizb_text` are pre-formatted Arabic labels, not numeric ids. Use `GetAyasByJuzNumber` for numeric (1..30) access.

## Tests

```bash
go test ./quran
```

## License

See repository root.

---

_Wait for more..._

Support me on PayPal

[![Donate](https://img.shields.io/badge/Support-PayPal-blue.svg)](https://paypal.me/mohammadmoustafa1)