package warlock

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
	"github.com/wowsims/sod/sim/core/stats"
)

func (warlock *Warlock) registerIncinerateSpell() {
	if !warlock.HasRune(proto.WarlockRune_RuneLegsIncinerate) {
		return
	}
	spellCoeff := 0.714

	level := float64(warlock.GetCharacter().Level)
	baseCalc := (6.568597 + 0.672028*level + 0.031721*level*level)
	baseLowDamage := baseCalc * 2.22
	baseHighDamage := baseCalc * 2.58

	warlock.IncinerateAura = warlock.RegisterAura(core.Aura{
		Label:    "Incinerate Aura",
		ActionID: core.ActionID{SpellID: 412758},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.25
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warlock.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.25
		},
	})

	warlock.Incinerate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 412758},
		SpellSchool:  core.SpellSchoolFire,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagAPL | core.SpellFlagResetAttackSwing | core.SpellFlagBinary,
		MissileSpeed: 24,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.14,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2250,
			},
		},

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplier:         1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplierAdditive: 1,
		ThreatMultiplier:         1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage = sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellDamage()

			if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(target).IsActive() {
				baseDamage *= 1.4
			}

			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			warlock.IncinerateAura.Activate(sim)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}
