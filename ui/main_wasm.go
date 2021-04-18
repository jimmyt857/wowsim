package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"syscall/js"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

func main() {
	c := make(chan struct{}, 0)

	simfunc := js.FuncOf(Simulate)
	statfunc := js.FuncOf(StatWeight)
	statsfunc := js.FuncOf(ComputeStats)
	gearlistfunc := js.FuncOf(GearList)

	js.Global().Set("simulate", simfunc)
	js.Global().Set("statweight", statfunc)
	js.Global().Set("computestats", statsfunc)
	js.Global().Set("gearlist", gearlistfunc)
	js.Global().Call("wasmready")
	<-c
}

// GearList reports all items of gear to the UI to display.
func GearList(this js.Value, args []js.Value) interface{} {
	slot := byte(128)

	if len(args) == 1 {
		slot = byte(args[0].Int())
	}
	gear := struct {
		Items    []tbc.Item
		Gems     []tbc.Gem
		Enchants []tbc.Enchant
	}{
		Items: make([]tbc.Item, 0, len(tbc.ItemLookup)),
	}
	for _, v := range tbc.ItemLookup {
		if slot != 128 && v.Slot != slot {
			continue
		}
		gear.Items = append(gear.Items, v)
	}
	gear.Gems = tbc.Gems
	gear.Enchants = tbc.Enchants

	output, err := json.Marshal(gear)
	if err != nil {
		// fmt.Printf("Failed to marshal gear list: %s", err)
		output = []byte(`{"error": ` + err.Error() + `}`)
	}
	// fmt.Printf("Item Output: %s", string(output))
	return string(output)
}

// GearStats takes a gear list and returns their total stats.
// This could power a simple 'current stats of all gear' UI.
func ComputeStats(this js.Value, args []js.Value) interface{} {
	gear := getGear(args[0])
	if len(args) != 2 {
		return `{"error": "incorrect args. expected computestats(gear, options)}`
	}
	if args[1].IsNull() {
		gearStats := gear.Stats()
		gearStats[tbc.StatSpellCrit] += (gearStats[tbc.StatInt] / 80) * 22.08
		gearStats[tbc.StatMana] += gearStats[tbc.StatInt] * 15
		return gearStats.Print(false)
	}
	opt := parseOptions(args[1])
	stats := opt.StatTotal(gear)
	opt.UseAI = true // stupid complaining sim...maybe I should just default AI on.
	fakesim := tbc.NewSim(stats, gear, opt)
	fakesim.ActivateSets()

	finalStats := stats
	for i, v := range fakesim.Buffs {
		finalStats[i] += v
	}
	return finalStats.Print(false)
}

// getGear converts js string array to a list of equipment items.
func getGear(val js.Value) tbc.Equipment {
	numGear := val.Length()
	gearSet := make([]tbc.Item, numGear)
	for i := range gearSet {
		v := val.Index(i)
		ic := tbc.ItemLookup[v.Get("Name").String()]
		gems := v.Get("Gems")
		if !(gems.IsUndefined() || gems.IsNull()) && gems.Length() > 0 {
			ic.Gems = make([]tbc.Gem, len(ic.GemSlots))
			for i := range ic.Gems {
				jsgem := gems.Index(i)
				if jsgem.IsNull() {
					continue
				}
				gv, ok := tbc.GemLookup[jsgem.String()]
				if !ok {
					continue // wasn't a valid gem
				}
				ic.Gems[i] = gv
			}
		}
		if !v.Get("Enchant").IsNull() && !v.Get("Enchant").IsUndefined() {
			ic.Enchant = tbc.EnchantLookup[v.Get("Enchant").String()]
		}
		gearSet[i] = ic
	}
	return tbc.Equipment(gearSet)
}

func parseOptions(val js.Value) tbc.Options {
	var custom = val.Get("custom")
	opt := tbc.Options{
		ExitOnOOM:    val.Get("exitoom").Truthy(),
		NumBloodlust: val.Get("buffbl").Int(),
		NumDrums:     val.Get("buffdrums").Int(),
		UseAI:        val.Get("useai").Truthy(),
		Buffs: tbc.Buffs{
			ArcaneInt:                val.Get("buffai").Truthy(),
			GiftOftheWild:            val.Get("buffgotw").Truthy(),
			BlessingOfKings:          val.Get("buffbk").Truthy(),
			ImprovedBlessingOfWisdom: val.Get("buffibow").Truthy(),
			ImprovedDivineSpirit:     val.Get("buffids").Truthy(),
			JudgementOfWisdom:        val.Get("debuffjow").Truthy(),
			ImpSealofCrusader:        val.Get("debuffisoc").Truthy(),
			Misery:                   val.Get("debuffmis").Truthy(),
			Moonkin:                  val.Get("buffmoon").Truthy(),
			MoonkinRavenGoddess:      val.Get("buffmoonrg").Truthy(),
			SpriestDPS:               val.Get("buffspriest").Int(),
			WaterShield:              val.Get("sbufws").Truthy(),
			EyeOfNight:               val.Get("buffeyenight").Truthy(),
			TwilightOwl:              val.Get("bufftwilightowl").Truthy(),
			Race:                     tbc.RaceBonusType(val.Get("sbufrace").Int()),
			Custom: tbc.Stats{
				tbc.StatInt:       custom.Get("custint").Float(),
				tbc.StatSpellCrit: custom.Get("custsc").Float(),
				tbc.StatSpellHit:  custom.Get("custsh").Float(),
				tbc.StatSpellDmg:  custom.Get("custsp").Float(),
				tbc.StatHaste:     custom.Get("custha").Float(),
				tbc.StatMP5:       custom.Get("custmp5").Float(),
				tbc.StatMana:      custom.Get("custmana").Float(),
			},
		},
		Consumes: tbc.Consumes{
			FlaskOfBlindingLight:     val.Get("confbl").Truthy(),
			FlaskOfMightyRestoration: val.Get("confmr").Truthy(),
			BrilliantWizardOil:       val.Get("conbwo").Truthy(),
			MajorMageblood:           val.Get("conmm").Truthy(),
			BlackendBasilisk:         val.Get("conbb").Truthy(),
			DestructionPotion:        val.Get("condp").Truthy(),
			SuperManaPotion:          val.Get("consmp").Truthy(),
			DarkRune:                 val.Get("condr").Truthy(),
		},
		Talents: tbc.Talents{
			LightninOverload:   5,
			ElementalPrecision: 3,
			NaturesGuidance:    3,
			TidalMastery:       5,
			ElementalMastery:   true,
			UnrelentingStorm:   3,
			CallOfThunder:      5,
		},
		Totems: tbc.Totems{
			TotemOfWrath: val.Get("totwr").Int(),
			WrathOfAir:   val.Get("totwoa").Truthy(),
			Cyclone2PC:   val.Get("totcycl2p").Truthy(),
			ManaStream:   val.Get("totms").Truthy(),
		},
	}

	return opt
}

func parseRotation(val js.Value) [][]string {

	out := [][]string{}

	for i := 0; i < val.Length(); i++ {
		rot := []string{}
		jsrot := val.Index(i)
		for j := 0; j < jsrot.Length(); j++ {
			rot = append(rot, jsrot.Index(j).String())
		}
		out = append(out, rot)
	}

	return out
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func StatWeight(this js.Value, args []js.Value) interface{} {
	numSims := args[0].Int()
	if numSims == 1 {
		tbc.IsDebug = true
	} else {
		tbc.IsDebug = false
	}
	seconds := args[1].Int()
	gear := getGear(args[2])
	opts := parseOptions(args[3])
	stat := args[4].Int()
	statModVal := args[5].Float()

	opts.Buffs.Custom = tbc.Stats{tbc.StatLen: 0}
	opts.Buffs.Custom[tbc.Stat(stat)] += statModVal
	opts.UseAI = true // use AI optimal rotation.

	simdmg := 0.0
	simmet := make([]tbc.SimMetrics, 0, numSims)

	opts.RSeed = time.Now().Unix()

	oomcount := 0
	sim := tbc.NewSim(opts.StatTotal(gear), gear, opts)
	for ns := 0; ns < numSims; ns++ {
		metrics := sim.Run(seconds)
		simdmg += metrics.TotalDamage
		simmet = append(simmet, metrics)
		if metrics.OOMAt > 0 && metrics.OOMAt < seconds-5 {
			oomcount++
		}
	}

	if float64(oomcount)/float64(numSims) > 0.25 {
		return -1
	}
	return simdmg / float64(numSims) / float64(seconds)
}

// Simulate takes in number of iterations, duration, a gear list, and simulation options.
// (iterations, duration, gearlist, options, <optional, custom rotation)
func Simulate(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 {
		print("Expected 4 min arguments:  (#iterations, duration, gearlist, options)")
		return `{"error": "invalid arguments supplied"}`
	}

	customRotation := [][]string{}
	customHaste := 0.0
	if len(args) >= 6 {
		if args[4].Truthy() {
			customRotation = parseRotation(args[4])
		}
		if args[5].Truthy() {
			customHaste = args[5].Float()
		}
	}
	gear := getGear(args[2])
	opt := parseOptions(args[3])
	stats := opt.StatTotal(gear)
	if customHaste != 0 {
		stats[tbc.StatHaste] = customHaste
	}

	simi := args[0].Int()
	if simi == 1 {
		tbc.IsDebug = true
	} else {
		tbc.IsDebug = false
	}
	dur := args[1].Int()
	fullLogs := false
	if len(args) > 6 {
		fullLogs = args[6].Truthy()
		fmt.Printf("Building Full Log:%v\n", fullLogs)
	}

	results := runTBCSim(opt, stats, gear, dur, simi, customRotation, fullLogs)
	st := time.Now()
	output, err := json.Marshal(results)
	if err != nil {
		print("Failed to json marshal results: ", err.Error())
	}
	fmt.Printf("Took %s to json marshal response.\n", time.Now().Sub(st))
	return string(output)
}

func jsonmarshal(results []SimResult) string {

	val := "["
	for i, v := range results {
		js, _ := json.Marshal(v)
		val += string(js)
		if i != len(results)-1 {
			val += ","
		}
	}
	val += "]"
	return val
}

type SimResult struct {
	Rotation     []string
	SimSeconds   int
	RealDuration float64
	Logs         string
	DPSAvg       float64              `json:"dps"`
	DPSDev       float64              `json:"dev"`
	MaxDPS       float64              `json:"max"`
	OOMAt        float64              `json:"oomat"`
	NumOOM       int                  `json:"numOOM"`
	DPSAtOOM     float64              `json:"dpsAtOOM"`
	Casts        map[int32]CastMetric `json:"casts"`
	DPSHist      map[int]int          `json:"dpsHist"` // rounded DPS to count
}

type CastMetric struct {
	Count int
	Dmg   float64
	Crits int
}

func runTBCSim(opts tbc.Options, stats tbc.Stats, equip tbc.Equipment, seconds int, numSims int, customRotation [][]string, fullLogs bool) []SimResult {
	print("\nSim Duration:", seconds)
	print("\nNum Simulations: ", numSims)
	print("\n")

	spellOrders := [][]string{}
	doingCustom := false
	if len(customRotation) > 0 {
		doingCustom = true
		spellOrders = customRotation
	}
	results := []SimResult{}
	logsBuffer := &strings.Builder{}

	dosim := func(spells []string, simsec int) {
		simMetrics := SimResult{
			DPSHist:  map[int]int{},
			Casts:    map[int32]CastMetric{},
			Rotation: spells,
		}
		if opts.UseAI {
			simMetrics.Rotation = []string{"AI Optimized"}
		}
		st := time.Now()
		rseed := time.Now().Unix()
		optNow := opts
		optNow.SpellOrder = spells
		optNow.RSeed = rseed
		sim := tbc.NewSim(stats, equip, optNow)

		var totalSq float64
		for ns := 0; ns < numSims; ns++ {
			if fullLogs {
				sim.Debug = func(s string, vals ...interface{}) {
					logsBuffer.WriteString(fmt.Sprintf("[%0.1f] "+s, append([]interface{}{(float64(sim.CurrentTick) / float64(tbc.TicksPerSecond))}, vals...)...))
				}
			}
			metrics := sim.Run(simsec)
			dps := metrics.TotalDamage / float64(simsec)
			totalSq += dps * dps
			simMetrics.DPSAvg += dps
			dpsRounded := int(math.Round(dps/10) * 10)
			simMetrics.DPSHist[dpsRounded] += 1
			if dps > simMetrics.MaxDPS {
				fmt.Printf("New max dps: %0.0f\n", dps)
				simMetrics.MaxDPS = dps
			}
			if (metrics.OOMAt) > 0 {
				simMetrics.OOMAt += float64(metrics.OOMAt)
				simMetrics.DPSAtOOM += float64(metrics.DamageAtOOM) / float64(metrics.OOMAt)
				simMetrics.NumOOM++
			}
		}

		meanSq := totalSq / float64(numSims)
		mean := simMetrics.DPSAvg / float64(numSims)
		stdev := math.Sqrt(meanSq - mean*mean)

		simMetrics.DPSDev = stdev
		simMetrics.DPSAvg /= float64(numSims)
		if simMetrics.NumOOM > 0 {
			simMetrics.OOMAt /= float64(simMetrics.NumOOM)
			simMetrics.DPSAtOOM /= float64(simMetrics.NumOOM)
		}

		simMetrics.Logs = logsBuffer.String()
		simMetrics.SimSeconds = simsec
		simMetrics.RealDuration = time.Now().Sub(st).Seconds()
		results = append(results, simMetrics)
	}

	if !doingCustom && opts.UseAI {
		dosim([]string{"AI Optimized"}, seconds) // Let AI determine best possible DPS
	} else {
		for _, spells := range spellOrders {
			dosim(spells, seconds)
		}
	}
	return results
}

// var castStats = {
// 	1: {count: 0, dmg: 0, crits: 0},
// 	2: {count: 0, dmg: 0, crits: 0},
// 	3: {count: 0, dmg: 0, crits: 0},
// 	999: {count: 0, dmg: 0, crits: 0},
// 	998: {count: 0, dmg: 0, crits: 0},
// }
// out.Casts.forEach((casts)=>{
// 	casts.forEach((cast)=>{
// 		var id = cast.ID
// 		if (cast.IsLO)  {
// 			id = 1000-cast.ID;
// 		}
// 		castStats[id].count += 1;
// 		castStats[id].dmg += cast.Dmg;
// 		if (cast.Crit) {
// 			castStats[id].crits += 1;
// 		}
// 	});
// });
