package tbc

type Options struct {
	SpellOrder []string
	UseAI      bool // when set true, the AI will modulate the rotations to maximize DPS and mana.
	RSeed      int64
	ExitOnOOM  bool

	NumBloodlust int
	NumDrums     int

	Buffs    Buffs
	Consumes Consumes
	Talents  Talents
	Totems   Totems

	Debug bool // enables debug printing.
	// TODO: could change this to be a func/stream consumer could provide,
	// make it easier to integrate into different output systems.
}

func (o Options) StatTotal(e Equipment) Stats {
	gearStats := e.Stats()
	stats := o.BaseStats()
	for i := range stats {
		stats[i] += gearStats[i]
	}

	stats = o.Talents.AddStats(o.Buffs.AddStats(o.Consumes.AddStats(o.Totems.AddStats(stats))))

	if o.Buffs.BlessingOfKings {
		stats[StatInt] *= 1.1 // blessing of kings
	}
	if o.Buffs.ImprovedDivineSpirit {
		stats[StatSpellDmg] += stats[StatSpirit] * 0.1
	}

	// Final calculations
	stats[StatSpellCrit] += (stats[StatInt] / 80) * 22.08
	stats[StatMana] += stats[StatInt] * 15
	stats[StatMP5] += stats[StatInt] * (0.02 * float64(o.Talents.UnrelentingStorm))
	// fmt.Printf("\fFinal MP5: %f", (stats[StatMP5] + (stats[StatInt] * 0.06)))
	return stats
}

func (o Options) BaseStats() Stats {
	stats := Stats{
		StatInt:       104,    // Base int for troll
		StatMana:      2678,   // level 70 shaman
		StatSpirit:    135,    // lvl 70 shaman
		StatSpellCrit: 48.576, // base crit for 70 sham
		StatLen:       0,
	}
	return stats
}

type Totems struct {
	TotemOfWrath int
	WrathOfAir   bool
	ManaStream   bool
	Cyclone2PC   bool // Cyclone set 2pc bonus
}

func (tt Totems) AddStats(s Stats) Stats {
	s[StatSpellCrit] += 66.24 * float64(tt.TotemOfWrath)
	s[StatSpellHit] += 37.8 * float64(tt.TotemOfWrath)
	if tt.WrathOfAir {
		s[StatSpellDmg] += 101
		if tt.Cyclone2PC {
			s[StatSpellDmg] += 20
		}
	}
	if tt.ManaStream {
		s[StatMP5] += 50
	}
	return s
}

type Talents struct {
	LightninOverload   int
	ElementalPrecision int
	NaturesGuidance    int
	TidalMastery       int
	ElementalMastery   bool
	UnrelentingStorm   int
	CallOfThunder      int
	Convection         int
	Concussion         float64 // temp hack to speed up not converting this to a int on every spell cast
}

func (t Talents) AddStats(s Stats) Stats {
	s[StatSpellHit] += 25.2 * float64(t.ElementalPrecision)
	s[StatSpellHit] += 12.6 * float64(t.NaturesGuidance)
	s[StatSpellCrit] += 22.08 * float64(t.TidalMastery)
	s[StatSpellCrit] += 22.08 * float64(t.CallOfThunder)

	return s
}

type Buffs struct {
	// Raid buffs
	ArcaneInt                bool
	GiftOftheWild            bool
	BlessingOfKings          bool
	ImprovedBlessingOfWisdom bool
	ImprovedDivineSpirit     bool

	// Party Buffs
	Moonkin             bool
	MoonkinRavenGoddess bool // adds 20 spell crit to moonkin aura
	SpriestDPS          int  // adds Mp5 ~ 25% (dps*5%*5sec = 25%)
	EyeOfNight          bool // Eye of night bonus from party member (not you)
	TwilightOwl         bool // from party member

	// Self Buffs
	WaterShield    bool
	WaterShieldPPM int // how many procs per minute does watershield get? Every 3 requires a recast.
	Race           RaceBonusType

	// Target Debuff
	JudgementOfWisdom bool
	Misery            bool

	// Custom
	Custom Stats
}

type RaceBonusType byte

// These values are used directly in the dropdown in index.html
const (
	RaceBonusNone RaceBonusType = iota
	RaceBonusDraenei
	RaceBonusTroll10
	RaceBonusTroll30
	RaceBonusOrc
)

func (b Buffs) AddStats(s Stats) Stats {
	if b.ArcaneInt {
		s[StatInt] += 40
	}
	if b.GiftOftheWild {
		s[StatInt] += 18 // assumes improved gotw, rounded down to nearest int... not sure if that is accurate.
	}
	if b.ImprovedBlessingOfWisdom {
		s[StatMP5] += 42
	}
	if b.Moonkin {
		s[StatSpellCrit] += 110.4
		if b.MoonkinRavenGoddess {
			s[StatSpellCrit] += 20
		}
	}
	if b.TwilightOwl {
		s[StatSpellCrit] += 44.16
	}
	if b.EyeOfNight {
		s[StatSpellDmg] += 34
	}
	if b.WaterShield {
		s[StatMP5] += 50
	}
	if b.Race == RaceBonusDraenei {
		s[StatSpellHit] += 15.76 // 1% hit
	}
	s[StatMP5] += float64(b.SpriestDPS) * 0.25

	for k, v := range b.Custom {
		s[k] += v
	}
	return s
}

type Consumes struct {
	// Buffs
	BrilliantWizardOil       bool
	MajorMageblood           bool
	FlaskOfBlindingLight     bool
	FlaskOfMightyRestoration bool
	BlackendBasilisk         bool

	// Used in rotations
	SuperManaPotion bool
	DarkRune        bool
}

func (c Consumes) AddStats(s Stats) Stats {
	if c.BrilliantWizardOil {
		s[StatSpellCrit] += 14
		s[StatSpellDmg] += 36
	}
	if c.MajorMageblood {
		s[StatMP5] += 16.0
	}
	if c.FlaskOfBlindingLight {
		s[StatSpellDmg] += 80
	}
	if c.FlaskOfMightyRestoration {
		s[StatMP5] += 25
	}
	if c.BlackendBasilisk {
		s[StatSpellDmg] += 23
	}
	return s
}
