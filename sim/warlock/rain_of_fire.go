package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/proto"
)

func (warlock *Warlock) getRainOfFireBaseConfig(rank int) core.SpellConfig {
	spellId := [5]int32{0, 5740, 6219, 11677, 11678}[rank]
	spellCoeff := [5]float64{0, .083, .083, .083, .083}[rank]
	baseDamage := [5]float64{0, 42, 92, 155, 226}[rank]
	manaCost := [5]float64{0, 295, 605, 885, 1185}[rank]
	level := [5]int{0, 20, 34, 46, 58}[rank]
	hasRune := warlock.HasRune(proto.WarlockRune_RuneChestLakeOfFire)

	return core.SpellConfig{
		ActionID:      core.ActionID{SpellID: spellId},
		SpellSchool:   core.SpellSchoolFire,
		DefenseType:   core.DefenseTypeMagic,
		ProcMask:      core.ProcMaskSpellDamage,
		Flags:         core.SpellFlagChanneled | core.SpellFlagAPL | core.SpellFlagResetAttackSwing,
		RequiredLevel: level,
		Rank:          rank,

		BonusCritRating: float64(warlock.Talents.Devastation) * core.SpellCritRatingPerCritChance,

		CritDamageBonus: warlock.ruin(),

		DamageMultiplierAdditive: 1 + 0.02*float64(warlock.Talents.Emberstorm),
		DamageMultiplier:         1,
		ThreatMultiplier:         1,

		ManaCost: core.ManaCostOptions{
			FlatCost:   manaCost,
			Multiplier: 1 - float64(warlock.Talents.Cataclysm)*0.01,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "RainOfFire-" + warlock.Label + strconv.Itoa(rank),
			},
			NumberOfTicks:       4,
			TickLength:          time.Second * 2,
			AffectedByCastSpeed: false,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = baseDamage + spellCoeff*dot.Spell.SpellDamage()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex][dot.Spell.CastType])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					targetDamage := dot.SnapshotBaseDamage
					if warlock.LakeOfFireAuras != nil && warlock.LakeOfFireAuras.Get(aoeTarget).IsActive() {
						targetDamage *= warlock.getLakeOfFireMultiplier()
					}
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)

					if hasRune && dot.TickCount == dot.NumberOfTicks {
						warlock.LakeOfFireAuras.Get(aoeTarget).Activate(sim)
					}
				}

			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
		},
	}
}

func (warlock *Warlock) getLakeOfFireMultiplier() float64 {
	return 1.5
}

func (warlock *Warlock) registerRainOfFireSpell() {
	hasRune := warlock.HasRune(proto.WarlockRune_RuneChestLakeOfFire)
	if hasRune {
		warlock.LakeOfFireAuras = warlock.NewEnemyAuraArray(func(unit *core.Unit, level int32) *core.Aura {
			return unit.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 403650},
				Label:    "Lake of Fire",
				Duration: time.Second * 15,
			})
		})
	}

	maxRank := 4

	for i := 1; i <= maxRank; i++ {
		config := warlock.getRainOfFireBaseConfig(i)

		if config.RequiredLevel <= int(warlock.Level) {
			warlock.RainOfFire = warlock.GetOrRegisterSpell(config)
		}
	}
}
