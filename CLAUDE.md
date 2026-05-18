# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- Build: `go build`
- Run: `go run .` (expects `quran/quran.db` to be reachable from the working directory)
- Test package: `go test ./quran`
- Run a single test: `go test ./quran -run TestName`

## Architecture

A Go library + thin `main.go` driver wrapping a read-only SQLite database (`quran/quran.db`) of the Quran. Uses the pure-Go `modernc.org/sqlite` driver (no CGo).

- `main.go` — entry point; constructs a `QuranManager` and exercises a few demo lookups.
- `quran/` — the library package.
  - `quran.go` — `QuranManager` owns the `*sql.DB`. `NewQuranManager` opens the DB and sets WAL + `synchronous=NORMAL` + 64MB cache PRAGMAs. Also holds `QuranInfo` with hardcoded totals (6236 verses, 114 suras, 30 juzs, 60 hizbs). Page/ruku counts are *not* available — the DB has no such columns.
  - `aya.go`, `sura.go` — domain structs + pure helpers (`IsMeccan`, `IsValid`, `ShouldHaveBasmalah`, etc.).
  - `helpers.go` — `PrintTableSamples` debug printer.
  - `quran.db` — the SQLite data file (lives inside the package so consumers importing the module get one canonical path).

### Database schema (`quran/quran.db`, SQLite)

```
suras(id INTEGER PK, name TEXT, type TEXT, total_verses INTEGER)
verses(id INTEGER PK, text TEXT, text_simple TEXT,
       tafseer_muasr TEXT, tafseer_sadi TEXT,
       hizb_text TEXT, juz_text TEXT, sura_num INTEGER)
```

`verses.id` is the global ayah number (1..6236). `hizb_text` / `juz_text` are pre-formatted Arabic labels (e.g. `"الجزء الأول"`, `"الحزب الـ 1"`), not numeric IDs — use `GetAyasByJuzNumber(1..30)` if you need numeric access.

### Caching

`QuranManager` lazily builds two `sync.Once`-guarded caches on first use:
- **offsets/totals** — cumulative `total_verses` per sura, used to convert between global and local aya numbers without a SQL query.
- **juzs** — distinct `juz_text` labels in canonical order, used by `GetAyasByJuzNumber`.

### Conventions

- `SuraType` constants (`Meccan` = `"مكية"`, `Medinan` = `"مدنية"`) must equal the exact strings stored in `suras.type` — don't translate.
- Library methods return wrapped errors (`fmt.Errorf("...: %w", err)`). Do not call `log.Fatal` from inside the library.
- All new public methods on `QuranManager` should have a matching test in `quran_test.go`. Tests share one `QuranManager` via `TestMain`.
