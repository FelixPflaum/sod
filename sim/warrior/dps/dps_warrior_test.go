package dps

import (
	"testing"

	_ "github.com/wowsims/sod/sim/common" // imported to get item effects included.
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func init() {
	RegisterDpsWarrior()
}

func TestFury(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P2FuryTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_dw"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_2_fury"),
			Buffs:       core.FullBuffsPhase2,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Fury", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func TestArms(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Level:      40,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:     P2ArmsTalents,
			GearSet:     core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_2h"),
			Rotation:    core.GetAplRotation("../../../ui/warrior/apls", "phase_2_arms"),
			Buffs:       core.FullBuffsPhase1,
			Consumes:    Phase2Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arms", SpecOptions: PlayerOptionsArms},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

func BenchmarkSimulate(b *testing.B) {
	core.Each([]*proto.Player{
		{
			Race:               proto.Race_RaceOrc,
			Class:              proto.Class_ClassWarrior,
			Level:              40,
			Equipment:          core.GetGearSet("../../../ui/warrior/gear_sets", "phase_2_2h").GearSet,
			Rotation:           core.GetAplRotation("../../../ui/warrior/apls", "phase_2_arms").Rotation,
			Consumes:           Phase2Consumes.Consumes,
			Spec:               PlayerOptionsFury,
			TalentsString:      P2FuryTalents,
			Buffs:              core.FullIndividualBuffsPhase2,
			DistanceFromTarget: 0,
		},
	}, func(player *proto.Player) { core.SpecBenchmark(b, player) })
}

var P2ArmsTalents = "303050213525100001"

var P2FuryTalents = "-05050005405010051"

var Phase1Consumes = core.ConsumesCombo{
	Label: "Phase 1 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir: proto.AgilityElixir_ElixirOfLesserAgility,
		MainHandImbue: proto.WeaponImbue_WildStrikes,
		OffHandImbue:  proto.WeaponImbue_BlackfathomSharpeningStone,
		StrengthBuff:  proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var Phase2Consumes = core.ConsumesCombo{
	Label: "Phase 2 Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:     proto.AgilityElixir_ElixirOfAgility,
		DragonBreathChili: true,
		Food:              proto.Food_FoodSagefishDelight,
		MainHandImbue:     proto.WeaponImbue_WildStrikes,
		OffHandImbue:      proto.WeaponImbue_SolidSharpeningStone,
		StrengthBuff:      proto.StrengthBuff_ElixirOfOgresStrength,
	},
}

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage:    50,
	UseRecklessness: true,
	Shout:           proto.WarriorShout_WarriorShoutBattle,
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHit,
}
