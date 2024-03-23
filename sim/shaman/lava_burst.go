package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (shaman *Shaman) applyLavaBurst() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsLavaBurst) {
		return
	}

	shaman.LavaBurst = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(false))

	if shaman.HasRune(proto.ShamanRune_RuneChestOverload) {
		shaman.LavaBurstOverload = shaman.RegisterSpell(shaman.newLavaBurstSpellConfig(true))
	}
}

func (shaman *Shaman) newLavaBurstSpellConfig(isOverload bool) core.SpellConfig {
	level := float64(shaman.Level)
	baseCalc := 7.583798 + 0.471881*level + 0.036599*level*level
	baseDamageLow := baseCalc * 4.69
	baseDamageHigh := baseCalc * 6.05
	spellCoeff := .571
	castTime := time.Second * 2
	cooldown := time.Second * 8
	manaCost := .10

	flags := SpellFlagFocusable
	if !isOverload {
		flags |= core.SpellFlagAPL
	}

	canOverload := !isOverload && shaman.HasRune(proto.ShamanRune_RuneChestOverload)

	spell := core.SpellConfig{
		ActionID:     core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsLavaBurst)},
		SpellCode:    SpellCode_ShamanLavaBurst,
		SpellSchool:  core.SpellSchoolFire,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        flags,
		MissileSpeed: 20,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
			// Convection does not currently apply to Lava Burst in SoD
			// Multiplier: 1 - 0.02*float64(shaman.Talents.Convection),
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				// Lightning Mastery does not currently apply to Lava Burst in SoD
				// CastTime: time.Second*2 - time.Millisecond*200*time.Duration(shaman.Talents.LightningMastery),
				CastTime: castTime,
				GCD:      core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		CritDamageBonus: shaman.elementalFury(),

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)

				if canOverload && result.Landed() && sim.RandomFloat("LvB Overload") < ShamanOverloadChance {
					shaman.LavaBurstOverload.Cast(sim, target)
				}
			})
		},
	}

	if isOverload {
		shaman.applyOverloadModifiers(&spell)
	}

	return spell
}
