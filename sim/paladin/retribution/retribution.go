package retribution

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
	"github.com/wowsims/wotlk/sim/paladin"
)

func RegisterRetributionPaladin() {
	core.RegisterAgentFactory(
		proto.Player_RetributionPaladin{},
		proto.Spec_SpecRetributionPaladin,
		func(character core.Character, options proto.Player) core.Agent {
			return NewRetributionPaladin(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_RetributionPaladin) // I don't really understand this line
			if !ok {
				panic("Invalid spec value for Retribution Paladin!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewRetributionPaladin(character core.Character, options proto.Player) *RetributionPaladin {
	retOptions := options.GetRetributionPaladin()

	ret := &RetributionPaladin{
		Paladin:              paladin.NewPaladin(character, *retOptions.Talents),
		Rotation:             *retOptions.Rotation,
		Judgement:            retOptions.Options.Judgement,
		Seal:                 retOptions.Options.Seal,
		UseDivinePlea:        retOptions.Options.UseDivinePlea,
		DivinePleaPercentage: retOptions.Rotation.DivinePleaPercentage,
		ExoSlack:             retOptions.Rotation.ExoSlack,
		ConsSlack:            retOptions.Rotation.ConsSlack,
		HolyWrathThreshold:   retOptions.Rotation.HolyWrathThreshold,

		HasLightswornBattlegear2Pc: character.HasSetBonus(paladin.ItemSetLightswornBattlegear, 2),
	}
	ret.PaladinAura = retOptions.Options.Aura
	if retOptions.Rotation.CustomRotation != nil {
		ret.PriorityRotation = make([]int32, len(retOptions.Rotation.CustomRotation.Spells))
		for i, customSpellProto := range retOptions.Rotation.CustomRotation.Spells {
			ret.PriorityRotation[i] = customSpellProto.Spell
		}
	}

	if retOptions.Rotation.Type == proto.RetributionPaladin_Rotation_Standart {
		ret.SelectedRotation = ret.mainRotation
	} else if retOptions.Rotation.Type == proto.RetributionPaladin_Rotation_Custom {
		ret.SelectedRotation = ret.costumeRotation
	} else {
		ret.SelectedRotation = ret.mainRotation
	}

	// Convert DTPS option to bonus MP5
	spAtt := retOptions.Options.DamageTakenPerSecond * 5.0 / 10.0
	ret.AddStat(stats.MP5, spAtt)

	ret.EnableAutoAttacks(ret, core.AutoAttackOptions{
		MainHand:       ret.WeaponFromMainHand(0), // Set crit multiplier later when we have targets.
		AutoSwingMelee: true,
	})

	ret.EnableResumeAfterManaWait(ret.OnGCDReady)

	return ret
}

type RetributionPaladin struct {
	*paladin.Paladin

	Judgement            proto.PaladinJudgement
	Seal                 proto.PaladinSeal
	UseDivinePlea        bool
	DivinePleaPercentage float64
	ExoSlack             int32
	ConsSlack            int32
	HolyWrathThreshold   int32

	SealInitComplete       bool
	DivinePleaInitComplete bool

	HasLightswornBattlegear2Pc bool

	SelectedRotation func(*core.Simulation)
	PriorityRotation []int32

	Rotation proto.RetributionPaladin_Rotation
}

func (ret *RetributionPaladin) GetPaladin() *paladin.Paladin {
	return ret.Paladin
}

func (ret *RetributionPaladin) Initialize() {
	ret.Paladin.Initialize()
	ret.RegisterAvengingWrathCD()

	ret.DelayDPSCooldownsForArmorDebuffs()
}

func (ret *RetributionPaladin) Reset(sim *core.Simulation) {
	ret.Paladin.Reset(sim)
	ret.AutoAttacks.CancelAutoSwing(sim)
	ret.SealInitComplete = false
	ret.DivinePleaInitComplete = false
}
