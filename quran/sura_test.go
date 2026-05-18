package quran

import "testing"

func TestSuraIsMeccanMedinan(t *testing.T) {
	m := Sura{Type: Meccan}
	if !m.IsMeccan() || m.IsMedinan() {
		t.Error("Meccan classification wrong")
	}
	d := Sura{Type: Medinan}
	if !d.IsMedinan() || d.IsMeccan() {
		t.Error("Medinan classification wrong")
	}
}

func TestShouldHaveBasmalah(t *testing.T) {
	if ShouldHaveBasmalah(9) {
		t.Error("sura 9 must not have basmalah")
	}
	if !ShouldHaveBasmalah(1) || !ShouldHaveBasmalah(114) {
		t.Error("suras 1 and 114 should have basmalah")
	}
	if ShouldHaveBasmalah(0) || ShouldHaveBasmalah(115) {
		t.Error("out-of-range suras should return false")
	}
}
