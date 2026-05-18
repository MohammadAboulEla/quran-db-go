package quran

/*
--- Suras Table Structure ---
  Column: id                   | Type: INTEGER    [PRIMARY KEY]
  Column: name                 | Type: TEXT
  Column: type                 | Type: TEXT
  Column: total_verses         | Type: INTEGER

--- Verses Table Structure ---
  Column: id                   | Type: INTEGER    [PRIMARY KEY]
  Column: text                 | Type: TEXT
  Column: text_simple          | Type: TEXT
  Column: tafseer_muasr        | Type: TEXT
  Column: tafseer_sadi         | Type: TEXT
  Column: hizb_text            | Type: TEXT
  Column: juz_text             | Type: TEXT
  Column: sura_num             | Type: INTEGER

==========================================
SURAS TABLE (FIRST 5 ROWS)
==========================================
ID| Name| Type| Total Verses
---------------
1|الفاتحة|مكية|7
2|البقرة|مدنية|286
3|آل عمران|مدنية|200
4|النساء|مدنية|176
5|المائدة|مدنية|120

==========================================
VERSES TABLE (FIRST 5 ROWS)
==========================================
ID| Sura| Juz| Hizb| Text Simple(Sample)| Tafseer Muasr(Sample)  
----------------------------
1|1|الجزء الأول|الحزب الـ 1|بسم الله الرحمن الرحيم|أبتدئ قراءة القرآن باس...  

2|1|الجزء الأول|الحزب الـ 1|الحمد لله رب العالمين|الثناء على الله بصفاته...  

3|1|الجزء الأول|الحزب الـ 1|الرحمن الرحيم|(ٱلرَّحۡمَٰنِ) ذي الرح...        

4|1|الجزء الأول|الحزب الـ 1|مالك يوم الدين|وهو سبحانه وحده مالك ي...  

5|1|الجزء الأول|الحزب الـ 1|إياك نعبد وإياك نستعين|إنا نخصك وحدك بالعبادة...  

*/

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

type QuranManager struct {
	db *sql.DB
	qf QuranInfo

	offsetsOnce sync.Once
	offsets     []int // offsets[suraNum] = sum of total_verses for suras with id < suraNum; index 0 unused
	offsetsErr  error
	totals      []int // totals[suraNum] = total_verses for that sura

	juzOnce sync.Once
	juzs    []string // distinct juz_text labels in order of first appearance
	juzErr  error
}

type QuranInfo struct {
	TotalVerses int
	TotalSuras  int
	TotalJuzs   int
	TotalHizbs  int
	TotalPages  int
	TotalRukus  int
}

func NewQuranManager(dbPath string) (*QuranManager, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	_, err = db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA cache_size = -64000;
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not set pragmas: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db is offline: %w", err)
	}

	qf := QuranInfo{
		TotalVerses: 6236,
		TotalSuras:  114,
		TotalJuzs:   30,
		TotalHizbs:  60,
		TotalPages:  604,
		TotalRukus:  558,
	}

	// Return an instance of our QuranManager holding the db connection
	return &QuranManager{db: db, qf: qf}, nil
}

func (qm *QuranManager) GetQuranInfo() QuranInfo {
	return qm.qf
}

func (qm *QuranManager) GetSuraName(suraNum int) (string, error) {
	var suraName string
	err := qm.db.QueryRow("SELECT name FROM suras WHERE id = ?", suraNum).Scan(&suraName)
	if err != nil {
		return "", fmt.Errorf("could not get sura name for %d: %w", suraNum, err)
	}
	return suraName, nil
}

func (qm *QuranManager) GetSuraTotalVerses(suraNum int) (int, error) {
	var totalVerses int
	err := qm.db.QueryRow("SELECT total_verses FROM suras WHERE id = ?", suraNum).Scan(&totalVerses)
	if err != nil {
		return 0, fmt.Errorf("could not get sura total verses for %d: %w", suraNum, err)
	}
	return totalVerses, nil
}

func (qm *QuranManager) GetSuraType(suraNum int) (string, error) {
	var suraType string
	err := qm.db.QueryRow("SELECT type FROM suras WHERE id = ?", suraNum).Scan(&suraType)
	if err != nil {
		return "", fmt.Errorf("could not get sura type for %d: %w", suraNum, err)
	}
	return suraType, nil
}

func (qm *QuranManager) Close() error {
	return qm.db.Close()
}

// GetSuraText returns the full Uthmani text of a sura with ayas joined by sep.
// Pass " " for a single-line read, "\n" for one aya per line.
func (qm *QuranManager) GetSuraText(suraNum int, sep string) (string, error) {
	if suraNum < 1 || suraNum > 114 {
		return "", fmt.Errorf("sura number %d out of range [1, 114]", suraNum)
	}
	rows, err := qm.db.Query("SELECT text FROM verses WHERE sura_num = ? ORDER BY id", suraNum)
	if err != nil {
		return "", fmt.Errorf("could not get text of sura %d: %w", suraNum, err)
	}
	defer rows.Close()
	var parts []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return "", err
		}
		parts = append(parts, t)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return strings.Join(parts, sep), nil
}

// GetAyasByGlobalRange returns ayas with global id in [from, to] inclusive.
func (qm *QuranManager) GetAyasByGlobalRange(from, to int) ([]Aya, error) {
	if from < 1 || to > 6236 || from > to {
		return nil, fmt.Errorf("invalid range [%d, %d]", from, to)
	}
	rows, err := qm.db.Query(
		"SELECT "+ayaSelectCols+" FROM verses WHERE id BETWEEN ? AND ? ORDER BY id",
		from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get ayas in range [%d, %d]: %w", from, to, err)
	}
	defer rows.Close()
	return qm.scanAyas(rows)
}

// GlobalToLocal converts a global aya number (1..6236) to (suraNum, localNumber).
func (qm *QuranManager) GlobalToLocal(global int) (int, int, error) {
	if global < 1 || global > 6236 {
		return 0, 0, fmt.Errorf("global number %d out of range [1, 6236]", global)
	}
	if err := qm.loadOffsets(); err != nil {
		return 0, 0, err
	}
	// Linear scan over 114 is trivial; binary search not worth the complexity.
	for s := 114; s >= 1; s-- {
		if qm.offsets[s] < global {
			return s, global - qm.offsets[s], nil
		}
	}
	return 0, 0, fmt.Errorf("could not locate sura for global %d", global)
}

// LocalToGlobal is the inverse of GlobalToLocal.
func (qm *QuranManager) LocalToGlobal(suraNum, local int) (int, error) {
	if suraNum < 1 || suraNum > 114 {
		return 0, fmt.Errorf("sura number %d out of range [1, 114]", suraNum)
	}
	if err := qm.loadOffsets(); err != nil {
		return 0, err
	}
	if local < 1 || local > qm.totals[suraNum] {
		return 0, fmt.Errorf("local %d out of range [1, %d] for sura %d", local, qm.totals[suraNum], suraNum)
	}
	return qm.offsets[suraNum] + local, nil
}

// GetSurasByType returns metadata for all suras of the given type.
func (qm *QuranManager) GetSurasByType(t SuraType) ([]Sura, error) {
	rows, err := qm.db.Query(
		"SELECT id, type, total_verses FROM suras WHERE type = ? ORDER BY id", string(t))
	if err != nil {
		return nil, fmt.Errorf("could not list suras of type %q: %w", t, err)
	}
	defer rows.Close()
	var out []Sura
	for rows.Next() {
		var id, total int
		var typ string
		if err := rows.Scan(&id, &typ, &total); err != nil {
			return nil, err
		}
		out = append(out, Sura{
			Number:           id,
			Type:             SuraType(typ),
			HasBasmalh:       ShouldHaveBasmalah(id),
			TotalAyasNumbers: total,
		})
	}
	return out, rows.Err()
}

// GetSuraNames returns the names of all 114 suras in order.
func (qm *QuranManager) GetSuraNames() ([]string, error) {
	rows, err := qm.db.Query("SELECT name FROM suras ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("could not list sura names: %w", err)
	}
	defer rows.Close()
	out := make([]string, 0, 114)
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// GetTafseerMuasr returns the contemporary tafseer for a given global aya number.
func (qm *QuranManager) GetTafseerMuasr(global int) (string, error) {
	if global < 1 || global > 6236 {
		return "", fmt.Errorf("global %d out of range", global)
	}
	var s string
	if err := qm.db.QueryRow("SELECT tafseer_muasr FROM verses WHERE id = ?", global).Scan(&s); err != nil {
		return "", fmt.Errorf("could not get tafseer_muasr for %d: %w", global, err)
	}
	return s, nil
}

// GetTafseerSadi returns Sa'di's tafseer for a given global aya number.
func (qm *QuranManager) GetTafseerSadi(global int) (string, error) {
	if global < 1 || global > 6236 {
		return "", fmt.Errorf("global %d out of range", global)
	}
	var s string
	if err := qm.db.QueryRow("SELECT tafseer_sadi FROM verses WHERE id = ?", global).Scan(&s); err != nil {
		return "", fmt.Errorf("could not get tafseer_sadi for %d: %w", global, err)
	}
	return s, nil
}

// juzLabelsOnce caches the 30 distinct juz_text labels in order of first appearance.
func (qm *QuranManager) juzLabels() ([]string, error) {
	qm.juzOnce.Do(func() {
		rows, err := qm.db.Query(
			"SELECT juz_text FROM verses GROUP BY juz_text ORDER BY MIN(id)")
		if err != nil {
			qm.juzErr = err
			return
		}
		defer rows.Close()
		for rows.Next() {
			var s string
			if err := rows.Scan(&s); err != nil {
				qm.juzErr = err
				return
			}
			qm.juzs = append(qm.juzs, s)
		}
		qm.juzErr = rows.Err()
	})
	return qm.juzs, qm.juzErr
}

// GetAyasByJuzNumber returns all ayas in juz n (1..30).
func (qm *QuranManager) GetAyasByJuzNumber(n int) ([]Aya, error) {
	if n < 1 || n > 30 {
		return nil, fmt.Errorf("juz number %d out of range [1, 30]", n)
	}
	labels, err := qm.juzLabels()
	if err != nil {
		return nil, err
	}
	if len(labels) < n {
		return nil, fmt.Errorf("expected 30 juz labels, db has %d", len(labels))
	}
	return qm.GetAyasByJuz(labels[n-1])
}

// scanAyas reads rows produced by a "SELECT id, text, text_simple, ..." query
// and fills LocalNumber using the cached offsets.
func (qm *QuranManager) scanAyas(rows *sql.Rows) ([]Aya, error) {
	if err := qm.loadOffsets(); err != nil {
		return nil, err
	}
	var out []Aya
	for rows.Next() {
		var a Aya
		if err := rows.Scan(
			&a.GlobalNumber, &a.Text, &a.TextSimple,
			&a.TafsserMuasar, &a.TafseerSadi,
			&a.HizbText, &a.JuzText, &a.SuraNumber,
		); err != nil {
			return nil, err
		}
		a.LocalNumber = a.GlobalNumber - qm.offsets[a.SuraNumber]
		out = append(out, a)
	}
	return out, rows.Err()
}

const ayaSelectCols = "id, text, text_simple, tafseer_muasr, tafseer_sadi, hizb_text, juz_text, sura_num"

// GetAyasByJuz returns all ayas whose juz_text equals the given Arabic label
// (e.g. "الجزء الأول"). The verses table stores juz as a text label, not a number.
func (qm *QuranManager) GetAyasByJuz(juzText string) ([]Aya, error) {
	rows, err := qm.db.Query(
		"SELECT "+ayaSelectCols+" FROM verses WHERE juz_text = ? ORDER BY id",
		juzText,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get ayas for juz %q: %w", juzText, err)
	}
	defer rows.Close()
	return qm.scanAyas(rows)
}

// GetAyasByHizb returns all ayas whose hizb_text matches. The DB stores per-quarter
// labels ("الحزب الـ 1", "ربع الحزب الـ 1", ...); pass the exact label you want.
func (qm *QuranManager) GetAyasByHizb(hizbText string) ([]Aya, error) {
	rows, err := qm.db.Query(
		"SELECT "+ayaSelectCols+" FROM verses WHERE hizb_text = ? ORDER BY id",
		hizbText,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get ayas for hizb %q: %w", hizbText, err)
	}
	defer rows.Close()
	return qm.scanAyas(rows)
}

// SearchAyas does a substring match on text_simple (LIKE %query%).
func (qm *QuranManager) SearchAyas(query string) ([]Aya, error) {
	if query == "" {
		return nil, fmt.Errorf("search query is empty")
	}
	rows, err := qm.db.Query(
		"SELECT "+ayaSelectCols+" FROM verses WHERE text_simple LIKE ? ORDER BY id",
		"%"+query+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()
	return qm.scanAyas(rows)
}

// GetSura returns sura metadata only (no ayas list populated).
func (qm *QuranManager) GetSura(suraNum int) (Sura, error) {
	if suraNum < 1 || suraNum > 114 {
		return Sura{}, fmt.Errorf("sura number %d out of range [1, 114]", suraNum)
	}
	var typ string
	var total int
	err := qm.db.QueryRow("SELECT type, total_verses FROM suras WHERE id = ?", suraNum).
		Scan(&typ, &total)
	if err != nil {
		return Sura{}, fmt.Errorf("could not get sura %d: %w", suraNum, err)
	}
	return Sura{
		Number:           suraNum,
		Type:             SuraType(typ),
		HasBasmalh:       ShouldHaveBasmalah(suraNum),
		TotalAyasNumbers: total,
	}, nil
}

// GetAllSuras returns metadata for all 114 suras (no ayas).
func (qm *QuranManager) GetAllSuras() ([]Sura, error) {
	rows, err := qm.db.Query("SELECT id, type, total_verses FROM suras ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("could not list suras: %w", err)
	}
	defer rows.Close()
	out := make([]Sura, 0, 114)
	for rows.Next() {
		var id, total int
		var typ string
		if err := rows.Scan(&id, &typ, &total); err != nil {
			return nil, err
		}
		out = append(out, Sura{
			Number:           id,
			Type:             SuraType(typ),
			HasBasmalh:       ShouldHaveBasmalah(id),
			TotalAyasNumbers: total,
		})
	}
	return out, rows.Err()
}

// GetSuraWithAyas returns the sura metadata along with its full AyasList.
func (qm *QuranManager) GetSuraWithAyas(suraNum int) (Sura, error) {
	s, err := qm.GetSura(suraNum)
	if err != nil {
		return Sura{}, err
	}
	if err := qm.loadOffsets(); err != nil {
		return Sura{}, err
	}
	rows, err := qm.db.Query(`
		SELECT id, text, text_simple, tafseer_muasr, tafseer_sadi, hizb_text, juz_text, sura_num
		FROM verses WHERE sura_num = ? ORDER BY id`, suraNum)
	if err != nil {
		return Sura{}, fmt.Errorf("could not get ayas of sura %d: %w", suraNum, err)
	}
	defer rows.Close()
	offset := qm.offsets[suraNum]
	ayas := make([]Aya, 0, s.TotalAyasNumbers)
	for rows.Next() {
		var a Aya
		if err := rows.Scan(
			&a.GlobalNumber, &a.Text, &a.TextSimple,
			&a.TafsserMuasar, &a.TafseerSadi,
			&a.HizbText, &a.JuzText, &a.SuraNumber,
		); err != nil {
			return Sura{}, err
		}
		a.LocalNumber = a.GlobalNumber - offset
		ayas = append(ayas, a)
	}
	if err := rows.Err(); err != nil {
		return Sura{}, err
	}
	s.AyasList = ayas
	return s, nil
}

// loadOffsets reads all sura totals once and builds a cumulative offset table.
// offsets[s] = number of ayas in suras 1..s-1, so global = offsets[s] + local.
func (qm *QuranManager) loadOffsets() error {
	qm.offsetsOnce.Do(func() {
		rows, err := qm.db.Query("SELECT id, total_verses FROM suras ORDER BY id")
		if err != nil {
			qm.offsetsErr = fmt.Errorf("could not load sura offsets: %w", err)
			return
		}
		defer rows.Close()
		totals := make([]int, 115)
		for rows.Next() {
			var id, tv int
			if err := rows.Scan(&id, &tv); err != nil {
				qm.offsetsErr = err
				return
			}
			if id >= 1 && id <= 114 {
				totals[id] = tv
			}
		}
		offsets := make([]int, 115)
		sum := 0
		for i := 1; i <= 114; i++ {
			offsets[i] = sum
			sum += totals[i]
		}
		qm.totals = totals
		qm.offsets = offsets
	})
	return qm.offsetsErr
}

// GetAyaByLocal returns the aya at position local (1..N) within suraNum.
func (qm *QuranManager) GetAyaByLocal(suraNum, local int) (Aya, error) {
	if suraNum < 1 || suraNum > 114 {
		return Aya{}, fmt.Errorf("sura number %d out of range [1, 114]", suraNum)
	}
	if err := qm.loadOffsets(); err != nil {
		return Aya{}, err
	}
	if local < 1 || local > qm.totals[suraNum] {
		return Aya{}, fmt.Errorf("local number %d out of range [1, %d] for sura %d", local, qm.totals[suraNum], suraNum)
	}
	return qm.GetAyaByGlobal(qm.offsets[suraNum] + local)
}

// GetAyaByGlobal returns the aya whose verses.id == global (1..6236).
// LocalNumber is computed from sura_num and global offset of the sura.
func (qm *QuranManager) GetAyaByGlobal(global int) (Aya, error) {
	if global < 1 || global > 6236 {
		return Aya{}, fmt.Errorf("global number %d out of range [1, 6236]", global)
	}
	var a Aya
	err := qm.db.QueryRow(`
		SELECT id, text, text_simple, tafseer_muasr, tafseer_sadi, hizb_text, juz_text, sura_num
		FROM verses WHERE id = ?`, global).Scan(
		&a.GlobalNumber, &a.Text, &a.TextSimple,
		&a.TafsserMuasar, &a.TafseerSadi,
		&a.HizbText, &a.JuzText, &a.SuraNumber,
	)
	if err != nil {
		return Aya{}, fmt.Errorf("could not get aya %d: %w", global, err)
	}

	if err := qm.loadOffsets(); err != nil {
		return Aya{}, err
	}
	a.LocalNumber = a.GlobalNumber - qm.offsets[a.SuraNumber]
	return a, nil
}
