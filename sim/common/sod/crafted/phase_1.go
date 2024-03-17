package crafted

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Cloth
	///////////////////////////////////////////////////////////////////////////

	// Extraplanar Spidersilk Boots
	core.NewItemEffect(210795, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 428489}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Planar Shift",
			ActionID: actionID,
			Duration: time.Second * 6,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.DamageDealtMultiplier *= .7
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.DamageDealtMultiplier /= .7
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeSurvival,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Leather
	///////////////////////////////////////////////////////////////////////////

	// Void-Touched Leather Gloves
	core.NewItemEffect(211423, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 429867}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Void Madness",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.1)
				character.PseudoStats.ThreatMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.MultiplyAttackSpeed(sim, 1.0/1.1)
				character.PseudoStats.ThreatMultiplier /= 1.2
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	// Void-Touched Leather Gauntlets
	core.NewItemEffect(211502, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 429868}

		buffAura := character.GetOrRegisterAura(core.Aura{
			Label:    "Void Madness",
			ActionID: actionID,
			Duration: time.Second * 10,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] *= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] *= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] *= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] *= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] *= 1.1
				character.PseudoStats.ThreatMultiplier *= 1.2
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexArcane] /= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFire] /= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexFrost] /= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexNature] /= 1.1
				character.PseudoStats.SchoolDamageDealtMultiplier[stats.SchoolIndexShadow] /= 1.1
				character.PseudoStats.ThreatMultiplier /= 1.2
			},
		})

		activationSpell := character.GetOrRegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
				buffAura.Activate(sim)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell:    activationSpell,
			Priority: core.CooldownPriorityLow,
			Type:     core.CooldownTypeDPS,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Mail
	///////////////////////////////////////////////////////////////////////////

	// Shifting Silver Breastplate
	// core.NewItemEffect(210794, func(agent core.Agent) {
	// 	// Nothing to do
	// })

	core.AddEffectsToTest = true
}
