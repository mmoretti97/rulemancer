(deffacts initial-game-state
  (game-state 
    (phase start-game-players)
    (turn-number 0)
    (active-player none)
    (priority-player none))
  
  (player-state 
    (player-id p1)
    (life 20)
    (max-mana 0)
    (current-mana 0)
    (lands-played 0)
    (has-priority no))
  
  (player-state 
    (player-id p2)
    (life 20)
    (max-mana 0)
    (current-mana 0)
    (lands-played 0)
    (has-priority no))
  
  (action-result (valid none) (reason "Game started."))
  
  ;(card (card-id p1-land1) (name "Forest") (type land) (owner p1) (zone hand) (mana-cost 0) (power 0) (toughness 0) (damage 0))
  ;(card (card-id p1-land2) (name "Forest") (type land) (owner p1) (zone hand) (mana-cost 0) (power 0) (toughness 0) (damage 0))
  ;(card (card-id p1-creature1) (name "Grizzly Bears") (type creature) (owner p1) (zone hand) (mana-cost 2) (power 2) (toughness 2) (damage 0))
  ;(card (card-id p1-creature2) (name "Giant Spider") (type creature) (owner p1) (zone hand) (mana-cost 4) (power 2) (toughness 4) (damage 0))
  
  ;(card (card-id p2-land1) (name "Mountain") (type land) (owner p2) (zone hand) (mana-cost 0) (power 0) (toughness 0) (damage 0))
  ;(card (card-id p2-land2) (name "Mountain") (type land) (owner p2) (zone hand) (mana-cost 0) (power 0) (toughness 0) (damage 0))
  ;(card (card-id p2-creature1) (name "Goblin Warrior") (type creature) (owner p2) (zone hand) (mana-cost 2) (power 2) (toughness 1) (damage 0))
  ;(card (card-id p2-creature2) (name "Fire Dragon") (type creature) (owner p2) (zone hand) (mana-cost 5) (power 5) (toughness 5) (damage 0))
  )