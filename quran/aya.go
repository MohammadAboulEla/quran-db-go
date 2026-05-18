package quran

import "fmt"

type Aya struct {
	GlobalNumber  int // from 1 to 6236
	LocalNumber   int // from 1 to N in sura
	Text          string
	TextSimple    string
	TafsserMuasar string
	TafseerSadi   string
	HizbText      string
	JuzText       string
	SuraNumber    int
}

func (a Aya) IsFirstInSura() bool { return a.LocalNumber == 1 }

func (a Aya) IsValid() bool {
	return a.GlobalNumber >= 1 && a.GlobalNumber <= 6236 &&
		a.SuraNumber >= 1 && a.SuraNumber <= 114 &&
		a.LocalNumber >= 1
}

func (a Aya) String() string {
	return fmt.Sprintf("[%d:%d] %s", a.SuraNumber, a.LocalNumber, a.TextSimple)
}
