package tbc

func (sim *Simulation) Run2(seconds int) SimMetrics {
	sim.reset()

	// Pop BL at start.
	for i := 0; i < sim.endTick; {
		if sim.CurrentMana < 0 {
			panic("you should never have negative mana.")
		}

		sim.CurrentTick = i
		advance := sim.Spellcasting(i)

		if sim.Options.ExitOnOOM && sim.metrics.OOMAt > 0 {
			return sim.metrics
		}

		sim.Advance(i, advance)
		i += advance

	}
	sim.metrics.ManaAtEnd = int(sim.CurrentMana)
	return sim.metrics
}

func (sim *Simulation) ActivateRacial() {
	switch v := sim.Options.Buffs.Race; v {
	case RaceBonusOrc:
		const spBonus = 143
		const dur = 15
		if sim.CDs[MagicIDOrcBloodFury] < 1 {
			sim.Buffs[StatSpellDmg] += spBonus
			sim.addAura(AuraStatRemoval(sim.CurrentTick, dur, spBonus, StatSpellDmg, MagicIDOrcBloodFury))
			sim.CDs[MagicIDOrcBloodFury] = 120 * TicksPerSecond
		}
	case RaceBonusTroll10, RaceBonusTroll30:
		hasteBonus := 157.6 // 10% haste
		const dur = 10
		if v == RaceBonusTroll30 {
			hasteBonus = 472.8 // 30% haste
		}
		if sim.CDs[MagicIDTrollBerserking] < 1 {
			sim.Buffs[StatHaste] += hasteBonus
			sim.addAura(AuraStatRemoval(sim.CurrentTick, dur, hasteBonus, StatHaste, MagicIDTrollBerserking))
			sim.CDs[MagicIDTrollBerserking] = 180 * TicksPerSecond
		}
	}
}

// Spellcasting will cast spells and calculate a new spell to cast.
//  Activates trinkets before spellcasting of off CD.
//  It will pop mana potions if needed.
func (sim *Simulation) Spellcasting(tickID int) int {
	// technically we dont really need this check with the new advancer.
	if sim.CastingSpell != nil && sim.CastingSpell.TicksUntilCast == 0 {
		sim.Cast(sim.CastingSpell)
	}

	if sim.CastingSpell == nil {
		if sim.Options.NumDrums > 0 && sim.CDs[MagicIDDrums] < 1 {
			// We have drums in the sim, and the drums aura isn't turned on.
			// Iterate our drum
			for i, v := range []int32{MagicIDDrum1, MagicIDDrum2, MagicIDDrum3, MagicIDDrum4} {
				if i == sim.Options.NumDrums {
					break
				}
				if sim.CDs[v] < 1 {
					sim.CDs[v] = 120 * TicksPerSecond // item goes on CD for 120s
					sim.addAura(ActivateDrums(sim))
					break
				}
			}
		}
		// Activate any specials
		if sim.Options.NumBloodlust > sim.bloodlustCasts && sim.CDs[MagicIDBloodlust] < 1 {
			sim.addAura(ActivateBloodlust(sim))
			sim.bloodlustCasts++ // TODO: will this break anything?
		}

		if sim.Options.Talents.ElementalMastery && sim.CDs[MagicIDEleMastery] < 1 {
			// Apply auras
			sim.addAura(AuraEleMastery())
		}

		sim.ActivateRacial()

		if sim.Options.Consumes.DestructionPotion && sim.CDs[MagicIDPotion] < 1 {
			// Only use dest potion if not using mana or if we haven't used it once.
			// If we are using mana, only use destruction potion on the pull.
			if !sim.Options.Consumes.SuperManaPotion || !sim.destructionPotion {
				sim.addAura(ActivateDestructionPotion(sim))
			}
		}

		didPot := false
		totalRegen := (sim.Stats[StatMP5] + sim.Buffs[StatMP5])
		// Pop potion before next cast if we have less than the mana provided by the potion minues 1mp5 tick.
		if sim.Options.Consumes.DarkRune && sim.Stats[StatMana]-sim.CurrentMana+totalRegen >= 1500 && sim.CDs[MagicIDRune] < 1 {
			// Restores 900 to 1500 mana. (2 Min Cooldown)
			sim.CurrentMana += 900 + (sim.rando.Float64() * 600)
			sim.CDs[MagicIDRune] = 120 * TicksPerSecond
			didPot = true
			if sim.Debug != nil {
				sim.Debug("Used Dark Rune\n")
			}
		}
		if sim.Options.Consumes.SuperManaPotion && sim.Stats[StatMana]-sim.CurrentMana+totalRegen >= 3000 && sim.CDs[MagicIDPotion] < 1 {
			// Restores 1800 to 3000 mana. (2 Min Cooldown)
			sim.CurrentMana += 1800 + (sim.rando.Float64() * 1200)
			sim.CDs[MagicIDPotion] = 120 * TicksPerSecond
			didPot = true
			if sim.Debug != nil {
				sim.Debug("Used Mana Potion\n")
			}
		}

		// Pop any on-use trinkets
		for _, item := range sim.activeEquip {
			if item.Activate == nil || item.ActivateCD == -1 { // ignore non-activatable, and always active items.
				continue
			}
			if sim.CDs[item.CoolID] > 0 {
				continue
			}
			if item.Slot == EquipTrinket && sim.CDs[MagicIDAllTrinket] > 0 {
				continue
			}
			sim.addAura(item.Activate(sim))
			sim.CDs[item.CoolID] = item.ActivateCD * TicksPerSecond
			if item.Slot == EquipTrinket {
				sim.CDs[MagicIDAllTrinket] = 30 * TicksPerSecond
			}
		}

		// Choose next spell
		ticks := sim.SpellChooser(sim, didPot)
		if sim.CastingSpell != nil {
			if sim.Debug != nil {
				sim.Debug("Start Casting %s Cast Time: %0.1fs\n", sim.CastingSpell.Spell.Name, float64(sim.CastingSpell.TicksUntilCast)/float64(TicksPerSecond))
			}
		}
		return ticks
	}

	return 1
}

// Advance moves time forward counting down auras, CDs, mana regen, etc
func (sim *Simulation) Advance(tickID int, ticks int) {

	if sim.CastingSpell != nil {
		sim.CastingSpell.TicksUntilCast -= ticks
	}

	// MP5 regen
	sim.CurrentMana += sim.manaRegen() * float64(ticks)

	if sim.CurrentMana > sim.Stats[StatMana] {
		sim.CurrentMana = sim.Stats[StatMana]
	}

	// CDS
	for k := range sim.CDs {
		sim.CDs[k] -= ticks
		if sim.CDs[k] < 1 {
			delete(sim.CDs, k)
		}
	}

	todel := []int{}
	for i := range sim.Auras {
		if sim.Auras[i].Expires <= (tickID + ticks) {
			todel = append(todel, i)
		}
	}
	for i := len(todel) - 1; i >= 0; i-- {
		sim.cleanAura(todel[i])
	}
}

func (sim *Simulation) manaRegen() float64 {
	return ((sim.Stats[StatMP5] + sim.Buffs[StatMP5]) / 5.0) / float64(TicksPerSecond)
}
