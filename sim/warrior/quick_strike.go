package warrior

import (
	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warrior *Warrior) registerQuickStrike() {
	if !warrior.HasRune(proto.WarriorRune_RuneQuickStrike) || warrior.GetMHWeapon().HandType != proto.HandType_HandTypeTwoHand {
		return
	}

	warrior.QuickStrike = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 429765},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagAPL | SpellFlagBloodSurge,

		RageCost: core.RageCostOptions{
			Cost:   20 - float64(warrior.Talents.ImprovedHeroicStrike) - warrior.FocusedRageDiscount,
			Refund: 0.8,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   warrior.critMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(0.10*spell.MeleeAttackPower(), 0.20*spell.MeleeAttackPower())

			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeWeaponSpecialHitAndCrit)

			if !result.Landed() {
				spell.IssueRefund(sim)
			}
		},
	})
}
