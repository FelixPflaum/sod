package priest

// func (priest *Priest) registerCircleOfHealingSpell() {
// 	targets := priest.Env.Raid.GetFirstNPlayersOrPets(5)

// 	priest.CircleOfHealing = priest.RegisterSpell(core.SpellConfig{
// 		ActionID:    core.ActionID{SpellID: 401946},
// 		SpellSchool: core.SpellSchoolHoly,
//      DefenseType: core.DefenseTypeMagic,
// 		ProcMask:    core.ProcMaskSpellHealing,
// 		Flags:       core.SpellFlagHelpful | core.SpellFlagAPL,

// 		ManaCost: core.ManaCostOptions{
// 			BaseCost:   0.21,
// 			Multiplier: 1,
// 		},
// 		Cast: core.CastConfig{
// 			DefaultCast: core.Cast{
// 				GCD: core.GCDDefault,
// 			},
// 			CD: core.Cooldown{
// 				Timer:    priest.NewTimer(),
// 				Duration: time.Second * 6,
// 			},
// 		},

// 		BonusCritRating:  float64(priest.Talents.HolySpecialization) * 1 * core.CritRatingPerCritChance,

// 		DamageMultiplier: 1 + .02*float64(priest.Talents.SpiritualHealing),
// 		ThreatMultiplier: 1 - []float64{0, .07, .14, .20}[priest.Talents.SilentResolve],

// 		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
// 			healFromSP := 0.4029 * spell.HealingPower(target)
// 			for _, aoeTarget := range targets {
// 				baseHealing := sim.Roll(958, 1058) + healFromSP
// 				spell.CalcAndDealHealing(sim, aoeTarget, baseHealing, spell.OutcomeHealingCrit)
// 			}
// 		},
// 	})
// }
