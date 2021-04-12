package tbc

import (
	"math"
)

type Aura struct {
	ID      int32
	Expires int // ticks aura will apply

	OnCast         AuraEffect
	OnCastComplete AuraEffect
	OnSpellHit     AuraEffect
	OnStruck       AuraEffect
	OnExpire       AuraEffect
}

func AuraName(a int32) string {
	switch a {
	case MagicIDUnknown:
		return "Unknown"
	case MagicIDLOTalent:
		return "Lightning Overload Talent"
	case MagicIDJoW:
		return "Judgement Of Wisdom Aura"
	case MagicIDEleFocus:
		return "Elemental Focus"
	case MagicIDEleMastery:
		return "Elemental Mastery"
	case MagicIDStormcaller:
		return "Stormcaller"
	case MagicIDBlessingSilverCrescent:
		return "Blessing of the Silver Crescent"
	case MagicIDQuagsEye:
		return "Quags Eye"
	case MagicIDFungalFrenzy:
		return "Fungal Frenzy"
	case MagicIDBloodlust:
		return "Bloodlust"
	case MagicIDSkycall:
		return "Skycall"
	case MagicIDEnergized:
		return "Energized"
	case MagicIDNAC:
		return "Nature Alignment Crystal"
	case MagicIDChaoticSkyfire:
		return "Chaotic Skyfire"
	case MagicIDInsightfulEarthstorm:
		return "Insightful Earthstorm"
	case MagicIDMysticSkyfire:
		return "Mystic Skyfire"
	case MagicIDMysticFocus:
		return "Mystic Focus"
	case MagicIDEmberSkyfire:
		return "Ember Skyfire"
	case MagicIDLB12:
		return "LB12"
	case MagicIDCL6:
		return "CL6"
	case MagicIDTLCLB:
		return "TLC-LB"
	case MagicIDISCTrink:
		return "Trink"
	case MagicIDNACTrink:
		return "NACTrink"
	case MagicIDPotion:
		return "Potion"
	case MagicIDRune:
		return "Rune"
	case MagicIDAllTrinket:
		return "AllTrinket"
	case MagicIDSpellPower:
		return "SpellPower"
	case MagicIDRubySerpent:
		return "RubySerpent"
	case MagicIDCallOfTheNexus:
		return "CallOfTheNexus"
	case MagicIDDCC:
		return "Darkmoon Card Crusade"
	case MagicIDDCCBonus:
		return "Aura of the Crusade"
	case MagicIDScryerTrink:
		return "Scryer Trinket"
	case MagicIDRubySerpentTrink:
		return "Ruby Serpent Trinket"
	case MagicIDXiriTrink:
		return "Xiri Trinket"
	case MagicIDDrums:
		return "Drums of Battle"
	case MagicIDDrum1:
		return "Drum #1"
	case MagicIDDrum2:
		return "Drum #2"
	case MagicIDDrum3:
		return "Drum #3"
	case MagicIDDrum4:
		return "Drum #4"
	case MagicIDNetherstrike:
		return "Netherstrike Set"
	case MagicIDTwinStars:
		return "Twin Stars Set"
	case MagicIDTidefury:
		return "Tidefury Set"
	case MagicIDSpellstrike:
		return "Spellstrike Set"
	case MagicIDSpellstrikeInfusion:
		return "Spellstrike Infusion"
	case MagicIDManaEtched:
		return "Mana-Etched Set"
	case MagicIDManaEtchedHit:
		return "Mana-EtchedHit"
	case MagicIDManaEtchedInsight:
		return "Mana-EtchedInsight"
	}

	return "<<TODO: Add Aura name to switch!!>>"
}

// AuraEffects will mutate a cast or simulation state.
type AuraEffect func(sim *Simulation, c *Cast)

// List of all magic effects and spells and items and stuff that can go on CD or have an aura.
const (
	MagicIDUnknown int32 = iota
	//Spells
	MagicIDLB12
	MagicIDCL6
	MagicIDTLCLB

	// Auras
	MagicIDLOTalent
	MagicIDJoW
	MagicIDEleFocus
	MagicIDEleMastery
	MagicIDStormcaller
	MagicIDBlessingSilverCrescent
	MagicIDQuagsEye
	MagicIDFungalFrenzy
	MagicIDBloodlust
	MagicIDSkycall
	MagicIDEnergized
	MagicIDNAC
	MagicIDChaoticSkyfire
	MagicIDInsightfulEarthstorm
	MagicIDMysticSkyfire
	MagicIDMysticFocus
	MagicIDEmberSkyfire
	MagicIDSpellPower
	MagicIDRubySerpent
	MagicIDCallOfTheNexus
	MagicIDDCC
	MagicIDDCCBonus
	MagicIDDrums // drums effect
	MagicIDNetherstrike
	MagicIDTwinStars
	MagicIDTidefury
	MagicIDSpellstrike
	MagicIDSpellstrikeInfusion
	MagicIDManaEtched
	MagicIDManaEtchedHit
	MagicIDManaEtchedInsight
	MagicIDMisery

	//Items
	MagicIDISCTrink
	MagicIDNACTrink
	MagicIDPotion
	MagicIDRune
	MagicIDAllTrinket
	MagicIDScryerTrink
	MagicIDRubySerpentTrink
	MagicIDXiriTrink
	MagicIDDrum1 // Party drum item CDs
	MagicIDDrum2
	MagicIDDrum3
	MagicIDDrum4
)

func AuraJudgementOfWisdom() Aura {
	return Aura{
		ID:      MagicIDJoW,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			sim.debug(" +Judgement Of Wisdom: 74 mana\n")
			sim.CurrentMana += 74
		},
	}
}

// Currently coded into the 'cast' function because we need something to change 'final' damage.
// TODO: this could be coded as 'onspellhit'!
// func AuraMisery() Aura {
// 	return Aura{
// 		ID:      MagicIDMisery,
// 		Expires: math.MaxInt32,
// 		OnCastComplete: func(sim *Simulation, c *Cast) {
// 		},
// 	}
// }

func AuraLightningOverload(lvl int) Aura {
	chance := 0.04 * float64(lvl)
	return Aura{
		ID:      MagicIDLOTalent,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if c.Spell.ID != MagicIDLB12 && c.Spell.ID != MagicIDCL6 {
				return
			}
			if c.IsLO {
				return // can't proc LO on LO
			}
			if sim.rando.Float64() < chance {
				sim.debug(" +Lightning Overload\n")
				clone := &Cast{
					IsLO:    true,
					Spell:   c.Spell,
					Effects: []AuraEffect{func(sim *Simulation, c *Cast) { c.DidDmg /= 2 }},
				}
				sim.Cast(clone)
			}
		},
	}
}

func AuraElementalFocus(tick int) Aura {
	count := 2
	return Aura{
		ID:      MagicIDEleFocus,
		Expires: tick + (15 * TicksPerSecond),
		OnCast: func(sim *Simulation, c *Cast) {
			c.ManaCost *= .6 // reduced by 40%
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if c.ManaCost <= 0 {
				return // Don't consume charges from free spells.
			}
			count--
			if count == 0 {
				sim.removeAuraByID(MagicIDEleFocus)
			}
		},
	}
}

func AuraEleMastery() Aura {
	return Aura{
		ID:      MagicIDEleMastery,
		Expires: math.MaxInt32,
		OnCast: func(sim *Simulation, c *Cast) {
			c.ManaCost = 0
			sim.CDs[MagicIDEleMastery] = 180 * TicksPerSecond
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Crit += 1.01 // 101% chance of crit
			sim.removeAuraByID(MagicIDEleMastery)
		},
	}
}

func AuraStormcaller(tick int) Aura {
	return Aura{
		ID:      MagicIDStormcaller,
		Expires: tick + (8 * TicksPerSecond),
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Spellpower += 50
		},
	}
}

// createSpellDmgActivate creates a function for trinket activation that uses +spellpower
//  This is so we don't need a separate function for every spell power trinket.
func createSpellDmgActivate(id int32, sp float64, durSeconds int) ItemActivation {
	return func(sim *Simulation) Aura {
		return Aura{
			ID:      id,
			Expires: sim.currentTick + durSeconds*TicksPerSecond,
			OnCastComplete: func(sim *Simulation, c *Cast) {
				c.Spellpower += sp
			},
		}
	}
}

func ActivateQuagsEye(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	internalCD := 45 * TicksPerSecond
	return Aura{
		ID:      MagicIDQuagsEye,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+internalCD < sim.currentTick && sim.rando.Float64() < 0.1 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.currentTick, 6.0, hasteBonus, StatHaste, MagicIDFungalFrenzy))
				lastActivation = sim.currentTick
			}
		},
	}
}

func ActivateNexusHorn(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	internalCD := 45 * TicksPerSecond
	const spellBonus = 225.0
	return Aura{
		ID:      MagicIDQuagsEye,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+internalCD < sim.currentTick && sim.rando.Float64() < 0.2 {
				sim.Buffs[StatSpellDmg] += spellBonus
				sim.addAura(AuraStatRemoval(sim.currentTick, 10.0, spellBonus, StatSpellDmg, MagicIDCallOfTheNexus))
				lastActivation = sim.currentTick
			}
		},
	}
}

func ActivateDCC(sim *Simulation) Aura {
	const spellBonus = 18.0
	stacks := 0
	return Aura{
		ID:      MagicIDDCC,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if stacks < 10 {
				stacks++
				sim.Buffs[StatSpellDmg] += spellBonus
			}
			// Removal aura will refresh with new total spellpower based on stacks.
			//  This will remove the old stack removal buff.
			sim.addAura(Aura{
				ID:      MagicIDDCCBonus,
				Expires: sim.currentTick + (10 * TicksPerSecond),
				OnExpire: func(sim *Simulation, c *Cast) {
					sim.Buffs[StatSpellDmg] -= spellBonus * float64(stacks)
					stacks = 0
				},
			})
		},
	}
}

// AuraStatRemoval creates a general aura for removing any buff stat on expiring.
// This is useful for activations / effects that give temp stats.
func AuraStatRemoval(tick int, seconds int, amount float64, stat Stat, id int32) Aura {
	return Aura{
		ID:      id,
		Expires: tick + (seconds * TicksPerSecond),
		OnExpire: func(sim *Simulation, c *Cast) {
			sim.debug(" -%0.0f %s from %s\n", amount, stat.StatName(), AuraName(id))
			sim.Buffs[stat] -= amount
		},
	}
}

func ActivateDrums(sim *Simulation) Aura {
	sim.Buffs[StatHaste] += 80
	sim.CDs[MagicIDDrums] = 30 * TicksPerSecond
	return AuraStatRemoval(sim.currentTick, 30, 80, StatHaste, MagicIDDrums)
}

func ActivateBloodlust(sim *Simulation) Aura {
	sim.Buffs[StatHaste] += 472.8
	sim.CDs[MagicIDBloodlust] = 40 * TicksPerSecond // assumes that multiple BLs are different shaman.
	return AuraStatRemoval(sim.currentTick, 40, 472.8, StatHaste, MagicIDBloodlust)
}

func ActivateSkycall(sim *Simulation) Aura {
	const hasteBonus = 101
	return Aura{
		ID:      MagicIDSkycall,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if c.Spell.ID == MagicIDLB12 && sim.rando.Float64() < 0.15 {
				sim.debug(" +Skycall Energized- \n")
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.currentTick, 10, hasteBonus, StatHaste, MagicIDEnergized))
			}
		},
	}
}

func ActivateNAC(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDNAC,
		Expires: sim.currentTick + 300*TicksPerSecond,
		OnCast: func(sim *Simulation, c *Cast) {
			c.ManaCost *= 1.2
		},
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.Spellpower += 250
		},
	}
}

func ActivateCSD(sim *Simulation) Aura {
	return Aura{
		ID:      MagicIDChaoticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			c.CritBonus *= 1.03
		},
	}
}

func ActivateIED(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const icd = 15 * TicksPerSecond
	return Aura{
		ID:      MagicIDInsightfulEarthstorm,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+icd < sim.currentTick && sim.rando.Float64() < 0.04 {
				lastActivation = sim.currentTick
				sim.debug(" *Insightful Earthstorm Mana Restore - 300\n")
				sim.CurrentMana += 300
			}
		},
	}
}

func ActivateMSD(sim *Simulation) Aura {
	lastActivation := math.MinInt32
	const hasteBonus = 320.0
	const icd = 35 * TicksPerSecond
	return Aura{
		ID:      MagicIDMysticSkyfire,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if lastActivation+icd < sim.currentTick && sim.rando.Float64() < 0.15 {
				sim.Buffs[StatHaste] += hasteBonus
				sim.addAura(AuraStatRemoval(sim.currentTick, 4.0, hasteBonus, StatHaste, MagicIDMysticFocus))
				lastActivation = sim.currentTick
			}
		},
	}
}

func ActivateESD(sim *Simulation) Aura {
	sim.Buffs[StatInt] += (sim.Stats[StatInt] + sim.Buffs[StatInt]) * 0.02
	return Aura{
		ID:      MagicIDEmberSkyfire,
		Expires: math.MaxInt32,
	}
}

func ActivateSpellstrike(sim *Simulation) Aura {
	const spellBonus = 92.0
	const duration = 10.0
	return Aura{
		ID:      MagicIDSpellstrike,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.05 { // TODO: validate
				sim.addAura(Aura{
					ID:      MagicIDSpellstrikeInfusion,
					Expires: sim.currentTick + (duration * TicksPerSecond),
					OnCastComplete: func(sim *Simulation, c *Cast) {
						c.Spellpower += spellBonus
					},
				})
			}
		},
	}
}

func ActivateManaEtched(sim *Simulation) Aura {
	const spellBonus = 110.0
	const duration = 15.0
	return Aura{
		ID:      MagicIDManaEtched,
		Expires: math.MaxInt32,
		OnCastComplete: func(sim *Simulation, c *Cast) {
			if sim.rando.Float64() < 0.02 { // TODO: validate
				sim.addAura(Aura{
					ID:      MagicIDManaEtchedInsight,
					Expires: sim.currentTick + (duration * TicksPerSecond),
					OnCastComplete: func(sim *Simulation, c *Cast) {
						c.Spellpower += spellBonus
					},
				})
			}
		},
	}
}

func ActivateTLC(sim *Simulation) Aura {
	const spellBonus = 110.0
	const duration = 15.0

	tlcspell := spellmap[MagicIDTLCLB]
	const icd = 2.5 * TicksPerSecond

	charges := 0
	lastActivation := 0
	return Aura{
		ID:      MagicIDManaEtched,
		Expires: math.MaxInt32,
		OnSpellHit: func(sim *Simulation, c *Cast) {
			if lastActivation+icd >= sim.currentTick {
				return
			}
			if !c.DidCrit {
				return
			}
			lastActivation = sim.currentTick
			charges++
			sim.debug(" Lightning Capacitor Charges: %d\n", charges)
			if charges >= 3 {
				sim.debug(" Lightning Capacitor Triggered!\n")
				clone := &Cast{
					Spell: tlcspell,
				}
				sim.Cast(clone)
				charges = 0
			}
		},
	}
}
