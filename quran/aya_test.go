package quran

import (
	"strings"
	"testing"
)

func TestAyaIsFirstInSura(t *testing.T) {
	if !(Aya{LocalNumber: 1}).IsFirstInSura() {
		t.Error("LocalNumber 1 should be first")
	}
	if (Aya{LocalNumber: 2}).IsFirstInSura() {
		t.Error("LocalNumber 2 should not be first")
	}
}

func TestAyaIsValid(t *testing.T) {
	good := Aya{GlobalNumber: 1, LocalNumber: 1, SuraNumber: 1}
	if !good.IsValid() {
		t.Error("expected valid")
	}
	bad := []Aya{
		{GlobalNumber: 0, LocalNumber: 1, SuraNumber: 1},
		{GlobalNumber: 6237, LocalNumber: 1, SuraNumber: 1},
		{GlobalNumber: 1, LocalNumber: 0, SuraNumber: 1},
		{GlobalNumber: 1, LocalNumber: 1, SuraNumber: 115},
	}
	for i, a := range bad {
		if a.IsValid() {
			t.Errorf("case %d: expected invalid", i)
		}
	}
}

func TestAyaString(t *testing.T) {
	a := Aya{SuraNumber: 1, LocalNumber: 1, TextSimple: "بسم الله"}
	s := a.String()
	if !strings.Contains(s, "1:1") || !strings.Contains(s, "بسم الله") {
		t.Errorf("unexpected string: %q", s)
	}
}
