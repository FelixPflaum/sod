package druid

import (
	"time"

	"github.com/wowsims/sod/sim/core"
)

const CryingWind int32 = 45270

func (druid *Druid) registerInsectSwarmSpell() {
	missAuras := druid.NewEnemyAuraArray(core.InsectSwarmAura)
	idolSpellPower := core.TernaryFloat64(druid.Ranged().ID == CryingWind, 396, 0)

	impISMultiplier := 1 + 0.01*float64(druid.Talents.ImprovedInsectSwarm)

	druid.InsectSwarm = druid.RegisterSpell(Humanoid|Moonkin, core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48468},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagOmen | core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost:   0.08,
			Multiplier: 1,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		DamageMultiplier: 1 +
			0.01*float64(druid.Talents.Genesis) +
			core.TernaryFloat64(druid.HasSetBonus(ItemSetDreamwalkerGarb, 2), 0.1, 0),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Insect Swarm",
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					druid.Wrath.DamageMultiplier *= impISMultiplier
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					druid.Wrath.DamageMultiplier /= impISMultiplier
				},
			},
			NumberOfTicks: 6 + core.TernaryInt32(druid.Talents.NaturesSplendor, 1, 0),
			TickLength:    time.Second * 2,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = 215 + 0.2*(dot.Spell.SpellDamage()+idolSpellPower)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

				if druid.MoonkinT84PCAura != nil && sim.RandomFloat("Elune's Wrath proc") < 0.08 {
					druid.MoonkinT84PCAura.Activate(sim)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
				missAuras.Get(target).Activate(sim)
			}
			spell.DealOutcome(sim, result)
		},
	})

	druid.InsectSwarm.RelatedAuras = append(druid.InsectSwarm.RelatedAuras, missAuras)
}
