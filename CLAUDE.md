# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- Build: `go build`
- Run: `go run .` (expects `quran.db` in the working directory)
- Test package: `go test ./quran`
- Run a single test: `go test ./quran -run TestName`

## Architecture

A Go library + thin `main.go` driver wrapping a read-only SQLite database (`quran.db`) of the Quran. Uses the pure-Go `modernc.org/sqlite` driver (no CGo).

- `main.go` — entry point; constructs a `QuranManager` and calls demo helpers.
- `quran/` — the library package (imported as `. "quran-db/quran"` via dot-import in main).
  - `quran.go` — `QuranManager` owns the `*sql.DB`. `NewQuranManager` opens the DB and sets WAL + `synchronous=NORMAL` + 64MB cache PRAGMAs. Also holds `QuranInfo` with hardcoded totals (6236 verses, 114 suras, 30 juzs, 60 hizbs, 604 pages, 558 rukus).
  - `aya.go`, `sura.go` — domain structs. `SuraType` is the Arabic string `"مكية"` (Meccan) or `"مدنية"` (Medinan), matching the value stored in the `suras.type` column.
  - `helpers.go` — debug/sample printers.

### Database schema (`quran.db`, SQLite)

```
suras(id INTEGER PK, name TEXT, type TEXT, total_verses INTEGER)
verses(id INTEGER PK, text TEXT, text_simple TEXT,
       tafseer_muasr TEXT, tafseer_sadi TEXT,
       hizb_text TEXT, juz_text TEXT, sura_num INTEGER)
```

`verses.id` is the global ayah number (1..6236). `hizb_text` / `juz_text` are pre-formatted Arabic labels (e.g. `"الجزء الأول"`, `"الحزب الـ 1"`), not numeric IDs.

### Conventions

- Arabic text and comments are used throughout — preserve UTF-8 and don't translate identifiers/strings that map to DB values (e.g. `SuraType` constants must equal the exact column values).
- Existing code calls `log.Fatal` inside library methods on DB errors. This is a known smell — when adding new methods, prefer returning the error rather than copying the `log.Fatal` pattern.
