package tbc

import (
	"strconv"
)

const TicksPerSecond = 30

type Stats []float64

type Stat byte

const (
	StatInt Stat = iota
	StatStm
	StatSpellCrit
	StatSpellHit
	StatSpellDmg
	StatHaste
	StatMP5
	StatMana
	StatSpellPen
	StatSpirit

	StatLen
)

func (s Stat) StatName() string {
	switch s {
	case StatInt:
		return "StatInt"
	case StatStm:
		return "StatStm"
	case StatSpellCrit:
		return "StatSpellCrit"
	case StatSpellHit:
		return "StatSpellHit"
	case StatSpellDmg:
		return "StatSpellDmg"
	case StatHaste:
		return "StatHaste"
	case StatMP5:
		return "StatMP5"
	case StatMana:
		return "StatMana"
	case StatSpellPen:
		return "StatSpellPen"
	case StatSpirit:
		return "StatSpirit"
	}

	return "none"
}

func (st Stats) Clone() Stats {
	ns := make(Stats, StatLen)
	for i, v := range st {
		ns[i] = v
	}
	return ns
}

func (st Stats) Print(pretty bool) string {
	output := "{ "
	printed := false
	for k, v := range st {
		name := Stat(k).StatName()
		if name == "none" {
			continue
		}
		if printed {
			printed = false
			output += ","
			if pretty {
				output += "\n"
			}
		}
		if pretty {
			output += "\t"
		}
		if v < 50 {
			printed = true
			output += "\"" + name + "\": " + strconv.FormatFloat(v, 'f', 3, 64)
		} else {
			printed = true
			output += "\"" + name + "\": " + strconv.FormatFloat(v, 'f', 0, 64)
		}
	}
	output += " }"
	return output
}
