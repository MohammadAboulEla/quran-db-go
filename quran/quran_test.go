// test all functions of quran manager
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

package quran

import (
	"os"
	"strings"
	"testing"
)

var testQM *QuranManager

func TestMain(m *testing.M) {
	qm, err := NewQuranManager("quran.db")
	if err != nil {
		panic("could not open test db: " + err.Error())
	}
	testQM = qm
	code := m.Run()
	qm.Close()
	os.Exit(code)
}

func TestNewQuranManager(t *testing.T) {
	if testQM == nil {
		t.Fatal("test quran manager is nil")
	}
}

func TestGetQuranInfo(t *testing.T) {
	qf := testQM.GetQuranInfo()
	if qf.TotalVerses != 6236 || qf.TotalSuras != 114 || qf.TotalJuzs != 30 || qf.TotalHizbs != 60 {
		t.Errorf("incorrect quran info: %+v", qf)
	}
}

func TestGetSuraName(t *testing.T) {
	name, err := testQM.GetSuraName(1)
	if err != nil {
		t.Fatalf("could not get sura name: %v", err)
	}
	if name != "الفاتحة" {
		t.Errorf("incorrect sura name: %q", name)
	}
}

func TestGetSuraTotalVerses(t *testing.T) {
	total, err := testQM.GetSuraTotalVerses(1)
	if err != nil {
		t.Fatalf("could not get sura total verses: %v", err)
	}
	if total != 7 {
		t.Errorf("incorrect sura total verses: %d", total)
	}
}

func TestGetAyaByGlobal(t *testing.T) {
	cases := []struct {
		global   int
		sura     int
		local    int
		contains string
	}{
		{1, 1, 1, "بسم الله"},
		{7, 1, 7, ""},      // last aya of Al-Fatiha
		{8, 2, 1, ""},      // first aya of Al-Baqara
		{6236, 114, 6, ""}, // last aya of An-Nas
	}
	for _, c := range cases {
		a, err := testQM.GetAyaByGlobal(c.global)
		if err != nil {
			t.Fatalf("global %d: %v", c.global, err)
		}
		if a.SuraNumber != c.sura {
			t.Errorf("global %d: sura got %d want %d", c.global, a.SuraNumber, c.sura)
		}
		if a.LocalNumber != c.local {
			t.Errorf("global %d: local got %d want %d", c.global, a.LocalNumber, c.local)
		}
		if c.contains != "" && !strings.Contains(a.TextSimple, c.contains) {
			t.Errorf("global %d: text %q missing %q", c.global, a.TextSimple, c.contains)
		}
	}
}

func TestGetSuraText(t *testing.T) {
	text, err := testQM.GetSuraText(1, " ")
	if err != nil {
		t.Fatalf("GetSuraText(1): %v", err)
	}
	if text == "" {
		t.Fatal("empty sura text")
	}
	lines, err := testQM.GetSuraText(1, "\n")
	if err != nil {
		t.Fatalf("GetSuraText(1, \\n): %v", err)
	}
	if n := strings.Count(lines, "\n"); n != 6 {
		t.Errorf("Al-Fatiha joined by \\n should have 6 newlines (7 ayas), got %d", n)
	}
	if _, err := testQM.GetSuraText(0, " "); err == nil {
		t.Error("expected error for sura 0")
	}
}

func TestGetAyasByGlobalRange(t *testing.T) {
	ayas, err := testQM.GetAyasByGlobalRange(1, 7)
	if err != nil {
		t.Fatalf("GetAyasByGlobalRange(1,7): %v", err)
	}
	if len(ayas) != 7 {
		t.Fatalf("expected 7 ayas, got %d", len(ayas))
	}
	if ayas[0].GlobalNumber != 1 || ayas[6].GlobalNumber != 7 {
		t.Errorf("range bounds wrong: %d..%d", ayas[0].GlobalNumber, ayas[6].GlobalNumber)
	}
	one, err := testQM.GetAyasByGlobalRange(100, 100)
	if err != nil || len(one) != 1 {
		t.Errorf("single-element range failed: len=%d err=%v", len(one), err)
	}
	for _, bad := range [][2]int{{0, 5}, {5, 3}, {1, 6237}} {
		if _, err := testQM.GetAyasByGlobalRange(bad[0], bad[1]); err == nil {
			t.Errorf("expected error for range %v", bad)
		}
	}
}

func TestGlobalLocalRoundTrip(t *testing.T) {
	cases := []struct{ global, sura, local int }{
		{1, 1, 1},
		{7, 1, 7},
		{8, 2, 1},
		{293, 2, 286},
		{6236, 114, 6},
	}
	for _, c := range cases {
		s, l, err := testQM.GlobalToLocal(c.global)
		if err != nil {
			t.Fatalf("GlobalToLocal(%d): %v", c.global, err)
		}
		if s != c.sura || l != c.local {
			t.Errorf("GlobalToLocal(%d) = (%d,%d), want (%d,%d)", c.global, s, l, c.sura, c.local)
		}
		g, err := testQM.LocalToGlobal(c.sura, c.local)
		if err != nil {
			t.Fatalf("LocalToGlobal(%d,%d): %v", c.sura, c.local, err)
		}
		if g != c.global {
			t.Errorf("LocalToGlobal(%d,%d) = %d, want %d", c.sura, c.local, g, c.global)
		}
	}
	if _, _, err := testQM.GlobalToLocal(0); err == nil {
		t.Error("expected error for global 0")
	}
	if _, err := testQM.LocalToGlobal(1, 99); err == nil {
		t.Error("expected error for sura 1 local 99")
	}
}

func TestGetSurasByType(t *testing.T) {
	meccan, err := testQM.GetSurasByType(Meccan)
	if err != nil {
		t.Fatalf("GetSurasByType(Meccan): %v", err)
	}
	medinan, err := testQM.GetSurasByType(Medinan)
	if err != nil {
		t.Fatalf("GetSurasByType(Medinan): %v", err)
	}
	if len(meccan)+len(medinan) != 114 {
		t.Errorf("Meccan(%d) + Medinan(%d) != 114", len(meccan), len(medinan))
	}
	for _, s := range meccan {
		if s.Type != Meccan {
			t.Errorf("sura %d in meccan list has type %q", s.Number, s.Type)
		}
	}
}

func TestGetSuraNames(t *testing.T) {
	names, err := testQM.GetSuraNames()
	if err != nil {
		t.Fatalf("GetSuraNames: %v", err)
	}
	if len(names) != 114 {
		t.Fatalf("expected 114 names, got %d", len(names))
	}
	if names[0] != "الفاتحة" {
		t.Errorf("first name = %q, want الفاتحة", names[0])
	}
}

func TestGetTafseer(t *testing.T) {
	m, err := testQM.GetTafseerMuasr(1)
	if err != nil || m == "" {
		t.Errorf("GetTafseerMuasr(1): %q err=%v", m, err)
	}
	s, err := testQM.GetTafseerSadi(1)
	if err != nil || s == "" {
		t.Errorf("GetTafseerSadi(1): %q err=%v", s, err)
	}
	if _, err := testQM.GetTafseerMuasr(0); err == nil {
		t.Error("expected error for global 0")
	}
}

func TestGetAyasByJuzNumber(t *testing.T) {
	j1, err := testQM.GetAyasByJuzNumber(1)
	if err != nil {
		t.Fatalf("GetAyasByJuzNumber(1): %v", err)
	}
	if len(j1) == 0 || j1[0].GlobalNumber != 1 {
		t.Errorf("juz 1 should start at global 1, got len=%d first=%d", len(j1), j1[0].GlobalNumber)
	}
	j30, err := testQM.GetAyasByJuzNumber(30)
	if err != nil {
		t.Fatalf("GetAyasByJuzNumber(30): %v", err)
	}
	if j30[len(j30)-1].GlobalNumber != 6236 {
		t.Errorf("juz 30 should end at global 6236, got %d", j30[len(j30)-1].GlobalNumber)
	}
	total := 0
	for n := 1; n <= 30; n++ {
		ayas, err := testQM.GetAyasByJuzNumber(n)
		if err != nil {
			t.Fatalf("juz %d: %v", n, err)
		}
		total += len(ayas)
	}
	if total != 6236 {
		t.Errorf("sum across 30 juzs = %d, want 6236", total)
	}
	if _, err := testQM.GetAyasByJuzNumber(0); err == nil {
		t.Error("expected error for juz 0")
	}
	if _, err := testQM.GetAyasByJuzNumber(31); err == nil {
		t.Error("expected error for juz 31")
	}
}

func TestGetAyasByJuz(t *testing.T) {
	ayas, err := testQM.GetAyasByJuz("الجزء الأول")
	if err != nil {
		t.Fatalf("GetAyasByJuz: %v", err)
	}
	if len(ayas) == 0 {
		t.Fatal("expected ayas in juz 1")
	}
	if ayas[0].GlobalNumber != 1 {
		t.Errorf("juz 1 should start at global 1, got %d", ayas[0].GlobalNumber)
	}
	for _, a := range ayas {
		if a.JuzText != "الجزء الأول" {
			t.Fatalf("aya %d has wrong juz: %q", a.GlobalNumber, a.JuzText)
		}
	}

	empty, err := testQM.GetAyasByJuz("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("expected empty result, got %d", len(empty))
	}
}

func TestGetAyasByHizb(t *testing.T) {
	ayas, err := testQM.GetAyasByHizb("الحزب الـ 1")
	if err != nil {
		t.Fatalf("GetAyasByHizb: %v", err)
	}
	if len(ayas) == 0 {
		t.Fatal("expected ayas in hizb 1")
	}
	for _, a := range ayas {
		if a.HizbText != "الحزب الـ 1" {
			t.Fatalf("aya %d has wrong hizb: %q", a.GlobalNumber, a.HizbText)
		}
	}
}

func TestSearchAyas(t *testing.T) {
	results, err := testQM.SearchAyas("الرحمن")
	if err != nil {
		t.Fatalf("SearchAyas: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected matches for الرحمن")
	}
	for _, a := range results {
		if !strings.Contains(a.TextSimple, "الرحمن") {
			t.Errorf("aya %d does not contain query: %q", a.GlobalNumber, a.TextSimple)
		}
	}
	if _, err := testQM.SearchAyas(""); err == nil {
		t.Error("expected error for empty query")
	}
}

func TestGetSura(t *testing.T) {
	s, err := testQM.GetSura(1)
	if err != nil {
		t.Fatalf("GetSura(1): %v", err)
	}
	if s.Number != 1 || s.Type != Meccan || s.TotalAyasNumbers != 7 || !s.HasBasmalh {
		t.Errorf("unexpected sura 1: %+v", s)
	}
	s9, err := testQM.GetSura(9)
	if err != nil {
		t.Fatalf("GetSura(9): %v", err)
	}
	if s9.HasBasmalh {
		t.Error("sura 9 must not have basmalah")
	}
	if _, err := testQM.GetSura(0); err == nil {
		t.Error("expected error for sura 0")
	}
}

func TestGetAllSuras(t *testing.T) {
	all, err := testQM.GetAllSuras()
	if err != nil {
		t.Fatalf("GetAllSuras: %v", err)
	}
	if len(all) != 114 {
		t.Fatalf("expected 114 suras, got %d", len(all))
	}
	total := 0
	for _, s := range all {
		total += s.TotalAyasNumbers
	}
	if total != 6236 {
		t.Errorf("total ayas across all suras = %d, want 6236", total)
	}
	if all[0].Number != 1 || all[113].Number != 114 {
		t.Error("suras not ordered by id")
	}
}

func TestGetSuraWithAyas(t *testing.T) {
	s, err := testQM.GetSuraWithAyas(1)
	if err != nil {
		t.Fatalf("GetSuraWithAyas(1): %v", err)
	}
	if len(s.AyasList) != 7 {
		t.Fatalf("Al-Fatiha should have 7 ayas, got %d", len(s.AyasList))
	}
	if s.AyasList[0].LocalNumber != 1 || s.AyasList[6].LocalNumber != 7 {
		t.Error("local numbers not 1..7")
	}
	if s.AyasList[0].GlobalNumber != 1 {
		t.Errorf("first aya global = %d, want 1", s.AyasList[0].GlobalNumber)
	}

	s2, err := testQM.GetSuraWithAyas(2)
	if err != nil {
		t.Fatalf("GetSuraWithAyas(2): %v", err)
	}
	if len(s2.AyasList) != 286 || s2.AyasList[0].GlobalNumber != 8 {
		t.Errorf("Al-Baqara: len=%d firstGlobal=%d", len(s2.AyasList), s2.AyasList[0].GlobalNumber)
	}
}

func TestGetAyaByLocal(t *testing.T) {
	cases := []struct {
		sura, local, wantGlobal int
	}{
		{1, 1, 1},
		{1, 7, 7},
		{2, 1, 8},
		{2, 286, 293},
		{114, 6, 6236},
	}
	for _, c := range cases {
		a, err := testQM.GetAyaByLocal(c.sura, c.local)
		if err != nil {
			t.Fatalf("sura %d local %d: %v", c.sura, c.local, err)
		}
		if a.GlobalNumber != c.wantGlobal {
			t.Errorf("sura %d local %d: global got %d want %d", c.sura, c.local, a.GlobalNumber, c.wantGlobal)
		}
		if a.SuraNumber != c.sura || a.LocalNumber != c.local {
			t.Errorf("sura %d local %d: round-trip mismatch (%d:%d)", c.sura, c.local, a.SuraNumber, a.LocalNumber)
		}
	}
}

func TestGetAyaByLocalInvalid(t *testing.T) {
	cases := []struct{ sura, local int }{
		{0, 1},
		{115, 1},
		{1, 0},
		{1, 8},   // Al-Fatiha has 7
		{2, 287}, // Al-Baqara has 286
	}
	for _, c := range cases {
		if _, err := testQM.GetAyaByLocal(c.sura, c.local); err == nil {
			t.Errorf("expected error for sura %d local %d", c.sura, c.local)
		}
	}
}

func TestGetAyaByGlobalInvalid(t *testing.T) {
	for _, n := range []int{0, -1, 6237, 99999} {
		if _, err := testQM.GetAyaByGlobal(n); err == nil {
			t.Errorf("expected error for global %d", n)
		}
	}
}

func TestGetSuraType(t *testing.T) {
	suraType, err := testQM.GetSuraType(1)
	if err != nil {
		t.Fatalf("could not get sura type: %v", err)
	}
	if suraType != "مكية" {
		t.Errorf("incorrect sura type: %q", suraType)
	}
}
