package shaman

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

const ShamanMoltenBlastResetChance = .10

func (shaman *Shaman) applyMoltenBlast() {
	if !shaman.HasRune(proto.ShamanRune_RuneHandsMoltenBlast) {
		return
	}

	level := float64(shaman.Level)
	baseCalc := 7.583798 + 0.471881*level + 0.036599*level*level
	baseDamageLow := baseCalc * .72
	baseDamageHigh := baseCalc * 1.08
	apCoef := .05
	cooldown := time.Second * 6
	manaCost := .18
	targetCount := 4

	shaman.MoltenBlast = shaman.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: int32(proto.ShamanRune_RuneHandsMoltenBlast)},
		SpellCode:   SpellCode_ShamanMoltenBlast,
		SpellSchool: core.SpellSchoolFire,
		ProcMask:    core.ProcMaskMeleeOHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: cooldown,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   shaman.ElementalCritMultiplier(0),
		ThreatMultiplier: 2,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for i, aoeTarget := range sim.Encounter.TargetUnits {
				if i < targetCount {
					baseDamage := sim.Roll(baseDamageLow, baseDamageHigh) + apCoef*spell.MeleeAttackPower() + spell.BonusWeaponDamage()
					spell.CalcAndDealDamage(sim, aoeTarget, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
				}
			}
		},
	})
}
