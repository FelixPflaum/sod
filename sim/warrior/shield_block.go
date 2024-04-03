package warrior

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

// TODO: Classic Update
func (warrior *Warrior) RegisterShieldBlockCD() {
	actionID := core.ActionID{SpellID: 2565}
	cooldownDur := time.Second * 60

	warrior.ShieldBlockAura = warrior.RegisterAura(core.Aura{
		Label:    "Shield Block",
		ActionID: actionID,
		Duration: time.Second * 10,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, 100*core.BlockRatingPerBlockChance)
			warrior.PseudoStats.BlockValueMultiplier += 1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			warrior.AddStatDynamic(sim, stats.Block, -100*core.BlockRatingPerBlockChance)
			warrior.PseudoStats.BlockValueMultiplier -= 1
		},
	})

	warrior.ShieldBlock = warrior.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{},
			CD: core.Cooldown{
				Timer:    warrior.NewTimer(),
				Duration: cooldownDur,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return warrior.PseudoStats.CanBlock && (warrior.StanceMatches(DefensiveStance) || warrior.StanceMatches(GladiatorStance))
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			warrior.ShieldBlockAura.Activate(sim)
		},
	})

	warrior.AddMajorCooldown(core.MajorCooldown{
		Spell: warrior.ShieldBlock,
		Type:  core.CooldownTypeDPS,
	})
}
