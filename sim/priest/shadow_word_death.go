package priest

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (priest *Priest) registerShadowWordDeathSpell() {
	if !priest.HasRune(proto.PriestRune_RuneHandsShadowWordDeath) {
		return
	}

	level := float64(priest.GetCharacter().Level)
	// 2024-02-22 In-game value is ~66% base damage after tuning
	baseCalc := (9.456667 + 0.635108*level + 0.039063*level*level) * .66
	baseLowDamage := baseCalc * 5.32 * priest.darknessDamageModifier()
	baseHighDamage := baseCalc * 6.2 * priest.darknessDamageModifier()
	spellCoeff := 0.429
	manaCost := .12
	cooldown := time.Second * 12

	priest.ShadowWordDeath = priest.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 401955},
		SpellSchool: core.SpellSchoolShadow,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagAPL,

		ManaCost: core.ManaCostOptions{
			BaseCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    priest.NewTimer(),
				Duration: cooldown,
			},
		},

		BonusHitRating:  priest.shadowHitModifier(),
		BonusCritRating: priest.forceOfWillCritRating(),

		DamageMultiplier: priest.forceOfWillDamageModifier(),
		ThreatMultiplier: priest.shadowThreatModifier(),

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseLowDamage, baseHighDamage) + spellCoeff*spell.SpellDamage()
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				priest.AddShadowWeavingStack(sim, target)
			}
			spell.DealDamage(sim, result)
		},
		ExpectedInitialDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			baseDamage := (baseLowDamage+baseHighDamage)/2 + spellCoeff*spell.SpellDamage()
			return spell.CalcDamage(sim, target, baseDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	})
}
