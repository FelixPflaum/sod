package item_sets

import (
	"time"

	"github.com/wowsims/sod/sim/core"
	"github.com/wowsims/sod/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                                 Cloth
///////////////////////////////////////////////////////////////////////////

var ItemSetMalevolentProphetsVestments = core.NewItemSet(core.ItemSet{
	Name: "Malevolent Prophet's Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAuras := c.NewEnemyAuraArray(func(target *core.Unit, _ int32) *core.Aura {
				return target.GetOrRegisterAura(core.Aura{
					Label:    "Malelovance Proc",
					ActionID: core.ActionID{SpellID: 449920},
					Duration: time.Second * 30,
					//MaxStacks: 30, // If the debuff has stacks uncomment everything related to stacks

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						//aura.SetStacks(sim, aura.MaxStacks)

						// TODO: Blocked by Bonus Damage Taken feature implementation
						// for si := stats.SchoolIndexArcane; si < stats.SchoolLen; si++ {
						// 	aura.Unit.PseudoStats.BonusDamageTaken[si] += 50
						// }
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						// TODO: Blocked by Bonus Damage Taken feature implementation
						// for si := stats.SchoolIndexArcane; si < stats.SchoolLen; si++ {
						// 	aura.Unit.PseudoStats.BonusDamageTaken[si] -= 50
						// }
					},
					// OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					// 	if result.Landed() && spell.ProcMask.Matches(core.ProcMaskDirect) {
					// 		aura.RemoveStack(sim)
					// 	}
					// },
				})
			})

			handler := func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAuras.Get(result.Target).Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 449919},
				Name:       "Malelovance",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskSpellDamage,
				Outcome:    core.OutcomeLanded,
				ProcChance: 0.2,
				Handler:    handler,
			})
		},
	},
})

var ItemSetKnightLieutenantsDreadweave = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Dreadweave",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetBloodGuardsDreadweave = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Dreadweave",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 18)
		},
	},
})

var ItemSetKnightLieutenantsSatin = core.NewItemSet(core.ItemSet{
	Name: "Knight Lieutenant's Satin",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetBloodGuardsSatin = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Satin",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 33)
		},
	},
})

var ItemSetEmeraldEnchantedVestments = core.NewItemSet(core.ItemSet{
	Name: "Emerald Enchanted Vestments",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.SpellPower, 12)
		},
	},
})

var ItemSetEmeraldWovenGarb = core.NewItemSet(core.ItemSet{
	Name: "Emerald Woven Garb",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.HealingPower, 22)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Leather
///////////////////////////////////////////////////////////////////////////

var ItemSetKnightLieutenantsLeather = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetBloodGuardsLeather = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Leather",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetEmeraldLeathers = core.NewItemSet(core.ItemSet{
	Name: "Emerald Leathers",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Mail
///////////////////////////////////////////////////////////////////////////

var ItemSetShunnedDevoteesChainmail = core.NewItemSet(core.ItemSet{
	Name: "Shunned Devotee's Chainmail",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			// Holy Spell Crit
			c.OnSpellRegistered(func(spell *core.Spell) {
				if spell.SpellSchool.Matches(core.SpellSchoolHoly) {
					spell.BonusCritRating += 3
				}
			})

			// Nature Bonus Proc
			procAura := c.NewTemporaryStatsAura("The Furious Storm Proc", core.ActionID{SpellID: 449934}, stats.Stats{stats.NaturePower: 60, stats.HealingPower: 60}, time.Second*10)

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 449935},
				Name:       "The Furious Storm",
				Callback:   core.CallbackOnCastComplete,
				ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
				ProcChance: 0.10,
				Handler:    handler,
			})
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Plate
///////////////////////////////////////////////////////////////////////////

var ItemSetWailingBerserkersPlateArmor = core.NewItemSet(core.ItemSet{
	Name: "Wailing Berserker's Plate Armor",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeCrit, 1)
			c.AddStat(stats.SpellCrit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				c.AutoAttacks.ExtraMHAttack(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 449970},
				Name:       "Extra Attack",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMelee,
				ProcChance: 0.05,
				Handler:    handler,
			})
		},
	},
})

var ItemSetBanishedMartyrsFullPlate = core.NewItemSet(core.ItemSet{
	Name: "Banished Martyr's Full Plate",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.MeleeHit, 1)
			c.AddStat(stats.SpellHit, 1)
		},
		3: func(agent core.Agent) {
			c := agent.GetCharacter()

			procAura := c.NewTemporaryStatsAura("Stalwart Block Proc", core.ActionID{SpellID: 449975}, stats.Stats{stats.BlockValue: 50}, time.Second*6)

			handler := func(sim *core.Simulation, spell *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&c.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 449974},
				Name:       "Stalwart Block",
				Callback:   core.CallbackOnSpellHitTaken,
				ProcMask:   core.ProcMaskMelee,
				Outcome:    core.OutcomeBlock,
				ProcChance: 1,
				Handler:    handler,
			})
		},
	},
})

var ItemSetKnightLieutenantsPlate = core.NewItemSet(core.ItemSet{
	Name: "Knight-Lieutenant's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetBloodGuardsPlate = core.NewItemSet(core.ItemSet{
	Name: "Blood Guard's Plate",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 15)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 30)
		},
	},
})

var ItemSetEmeraldDreamPlate = core.NewItemSet(core.ItemSet{
	Name: "Emerald Dream Plate",
	Bonuses: map[int32]core.ApplyEffect{
		3: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 10)
		},
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.AttackPower, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                                 Other
///////////////////////////////////////////////////////////////////////////

var ItemSetSerpentsAscension = core.NewItemSet(core.ItemSet{
	Name: "Serpent's Ascension",
	Bonuses: map[int32]core.ApplyEffect{
		2: func(agent core.Agent) {
			character := agent.GetCharacter()
			procAura := character.NewTemporaryStatsAura("Serpent's Ascension Proc", core.ActionID{SpellID: 446231}, stats.Stats{stats.AttackPower: 150}, time.Second*12)

			handler := func(sim *core.Simulation, _ *core.Spell, _ *core.SpellResult) {
				procAura.Activate(sim)
			}

			core.MakeProcTriggerAura(&character.Unit, core.ProcTrigger{
				ActionID:   core.ActionID{SpellID: 446233},
				Name:       "Serpent's Ascension",
				Callback:   core.CallbackOnSpellHitDealt,
				ProcMask:   core.ProcMaskMeleeOrRanged,
				ProcChance: 0.03,
				ICD:        time.Second * 120,
				Handler:    handler,
			})
		},
	},
})
