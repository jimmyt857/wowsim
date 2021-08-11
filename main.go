package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/lologarithm/wowsim/tbc"
)

// Maps agentType flag values to actual enum types
var agentTypesMap = map[string]tbc.AgentType{
	"3LB1CL":   tbc.AGENT_TYPE_FIXED_3LB_1CL,
	"4LB1CL":   tbc.AGENT_TYPE_FIXED_4LB_1CL,
	"5LB1CL":   tbc.AGENT_TYPE_FIXED_5LB_1CL,
	"6LB1CL":   tbc.AGENT_TYPE_FIXED_6LB_1CL,
	"7LB1CL":   tbc.AGENT_TYPE_FIXED_7LB_1CL,
	"8LB1CL":   tbc.AGENT_TYPE_FIXED_8LB_1CL,
	"9LB1CL":   tbc.AGENT_TYPE_FIXED_9LB_1CL,
	"10LB1CL":  tbc.AGENT_TYPE_FIXED_10LB_1CL,
	"LB":       tbc.AGENT_TYPE_FIXED_LB_ONLY,
	"Adaptive": tbc.AGENT_TYPE_ADAPTIVE,
}

var DEFAULT_EQUIPMENT = tbc.EquipmentSpec{
	tbc.ItemSpec{ NameOrId: "Tidefury Helm" },
	tbc.ItemSpec{ NameOrId: "Charlotte's Ivy" },
	tbc.ItemSpec{ NameOrId: "Pauldrons of Wild Magic" },
	tbc.ItemSpec{ NameOrId: "Ogre Slayer's Cover" },
	tbc.ItemSpec{ NameOrId: "Tidefury Chestpiece" },
	tbc.ItemSpec{ NameOrId: "World's End Bracers" },
	tbc.ItemSpec{ NameOrId: "Earth Mantle Handwraps" },
	tbc.ItemSpec{ NameOrId: "Netherstrike Belt" },
	tbc.ItemSpec{ NameOrId: "Stormsong Kilt" },
	tbc.ItemSpec{ NameOrId: "Magma Plume Boots" },
	tbc.ItemSpec{ NameOrId: "Cobalt Band of Tyrigosa" },
	tbc.ItemSpec{ NameOrId: "Sparking Arcanite Ring" },
	tbc.ItemSpec{ NameOrId: "Mazthoril Honor Shield" },
	tbc.ItemSpec{ NameOrId: "Gavel of Unearthed Secrets" },
	tbc.ItemSpec{ NameOrId: "Natural Alignment Crystal" },
	tbc.ItemSpec{ NameOrId: "Icon of the Silver Crescent" },
	tbc.ItemSpec{ NameOrId: "Totem of the Void" },
}

var DEFAULT_OPTIONS = tbc.Options{
	NumBloodlust: 0,
	NumDrums:     0,
	Buffs: tbc.Buffs{
		ArcaneInt:                false,
		GiftOftheWild:            false,
		BlessingOfKings:          false,
		ImprovedBlessingOfWisdom: false,
		JudgementOfWisdom:        false,
		Moonkin:                  false,
		SpriestDPS:               0,
		WaterShield:              true,
		// Race:                     tbc.RaceBonusOrc,
		Custom: tbc.Stats{
			tbc.StatInt:       290,
			tbc.StatSpellDmg:  598 + 55,
			tbc.StatSpellHit:  24,
			tbc.StatSpellCrit: 120,
		},
	},
	Consumes: tbc.Consumes{
		// FlaskOfBlindingLight: true,
		// BrilliantWizardOil:   false,
		// MajorMageblood:       false,
		// BlackendBasilisk:     true,
		SuperManaPotion: false,
		// DarkRune:             false,
	},
	Talents: tbc.Talents{
		LightningOverload:  5,
		ElementalPrecision: 3,
		NaturesGuidance:    3,
		TidalMastery:       5,
		ElementalMastery:   true,
		UnrelentingStorm:   3,
		CallOfThunder:      5,
		Concussion:         5,
		Convection:         5,
	},
	Totems: tbc.Totems{
		TotemOfWrath: 1,
		WrathOfAir:   false,
		ManaStream:   true,
	},
}

// /script print(GetSpellBonusDamage(4))

func main() {

	// f, err := os.Create("profile2.cpu")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer f.Close() // error handling omitted for example
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	var isDebug = flag.Bool("debug", false, "Include --debug to spew the entire simulation log.")
	var noopt = flag.Bool("noopt", false, "If included it will disable optimization.")
	var agentTypeStr = flag.String("agentType", "", "Custom comma separated agent type to simulate.\n\tFor Example: --rotation=3LB1CL")
	var duration = flag.Float64("duration", 300, "Custom fight duration in seconds.")
	var iterations = flag.Int("iter", 10000, "Custom number of iterations for the sim to run.")
	var runWebUI = flag.Bool("web", false, "Use to run sim in web interface instead of in terminal")
	var configFile = flag.String("config", "", "Specify an input configuration.")

	flag.Parse()

	if *runWebUI {
		log.Printf("Closing: %s", http.ListenAndServe(":3333", nil))
	}

	simRequest := tbc.SimRequest{}
	if *configFile != "" {
		fileData, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to open config file(%s): %s", *configFile, err)
		}

		_ = json.Unmarshal([]byte(fileData), &simRequest)
	} else {
		simRequest.Options = DEFAULT_OPTIONS
		simRequest.Gear = DEFAULT_EQUIPMENT
	}

	if *isDebug {
		*iterations = 1
		simRequest.IncludeLogs = true
	}
	if agentTypeStr == nil {
		simRequest.Options.AgentType = tbc.AGENT_TYPE_ADAPTIVE
	} else {
		simRequest.Options.AgentType = agentTypesMap[*agentTypeStr]
	}
	simRequest.Options.Encounter.Duration = *duration
	simRequest.Options.RSeed = time.Now().Unix()
	simRequest.Iterations = *iterations

	runTBCSim(simRequest, *noopt)
}

func runTBCSim(simRequest tbc.SimRequest, noopt bool) {
	fmt.Printf(
			"\nSim Duration: %0.1f sec\nNum Simulations: %d\n",
			simRequest.Options.Encounter.Duration,
			simRequest.Iterations)

	equipment := tbc.NewEquipmentSet(simRequest.Gear)
	stats := tbc.CalculateTotalStats(simRequest.Options, equipment)
	fmt.Printf("\nFinal Stats: %s\n", stats.Print())

	if !noopt {
		weights := tbc.StatWeights(simRequest)
		fmt.Printf("Weights: [\n")
		for i, v := range weights {
			if tbc.Stat(i) == tbc.StatStm {
				continue
			}
			fmt.Printf("%s: %0.2f\t", tbc.Stat(i).StatName(), v)
		}
		fmt.Printf("\n]\n")
	}

	simResult := tbc.RunSimulation(simRequest)
	fmt.Printf("\n%s\n", simResultsToString(simRequest, simResult))
}

func simResultsToString(request tbc.SimRequest, result tbc.SimResult) string {
	output := ""
	output += fmt.Sprintf("Agent Type: %v\n", string(request.Options.AgentType))
	output += fmt.Sprintf("DPS:")
	output += fmt.Sprintf("\tMean: %0.1f +/- %0.1f\n", result.DpsAvg, result.DpsStDev)
	output += fmt.Sprintf("\tMax: %0.1f\n", result.DpsMax)
	output += fmt.Sprintf("Total Casts:\n")

	for castId, cast := range result.Casts {
		output += fmt.Sprintf("\t%s: %0.1f\n", tbc.AuraName(castId), float64(cast.Count) / float64(request.Iterations))
	}

	output += fmt.Sprintf("Went OOM: %d/%d sims\n", result.NumOom, request.Iterations)
	if result.NumOom > 0 {
		output += fmt.Sprintf("Avg OOM Time: %0.1f seconds\n", result.OomAtAvg)
		output += fmt.Sprintf("Avg DPS At OOM: %0.0f\n", result.DpsAtOomAvg)
	}
	output += fmt.Sprintf("Sim execution took %s", time.Duration(result.ExecutionDurationMs) * time.Millisecond)
	return output
}
