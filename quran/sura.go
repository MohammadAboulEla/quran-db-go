package quran

type Sura struct {
	Number           int
	Type             SuraType
	HasBasmalh       bool
	AyasList         []Aya
	TotalAyasNumbers int
}

type SuraType string

const (
	Meccan  SuraType = "مكية"
	Medinan SuraType = "مدنية"
)

func (s Sura) IsMeccan() bool  { return s.Type == Meccan }
func (s Sura) IsMedinan() bool { return s.Type == Medinan }

// At-Tawba (9) is the only sura that does not start with Basmalah.
func ShouldHaveBasmalah(suraNum int) bool {
	return suraNum != 9 && suraNum >= 1 && suraNum <= 114
}