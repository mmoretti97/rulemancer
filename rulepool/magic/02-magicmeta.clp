(deffacts magic-config
  (game-config
    (game-name magic)
    (description "Magic simplified for 2 players.")
    (num-players 2)))

(deffacts magic-interface
  (assertable
    (name mulligan-decision)
    (relations mulligan-decision))
  (assertable
    (name mulligan-cards-back-on-library)
    (relations mulligan-cards-back-on-library))
  (assertable
    (name play-land)
    (relations play-land))
  (assertable
    (name cast-creature)
    (relations cast-creature))
  (assertable
    (name declare-attacker)
    (relations declare-attacker))
  (assertable
    (name declare-blocker)
    (relations declare-blocker))
  (assertable
    (name pass-priority)
    (relations pass-priority))
  
  (results 
    (name mulligan-decision)
    (relations action-result))
  (results 
    (name mulligan-cards-back-on-library)
    (relations action-result))
  (results 
    (name cast-creature)
    (relations action-result))
  (results 
    (name declare-attacker)
    (relations action-result))
  (results 
    (name declare-blocker)
    (relations action-result))
  (results 
    (name pass-priority)
    (relations action-result))
  

  (queryable
    (name only-game-and-player-state)
    (relations game-state player-state))
  (queryable
    (name game-state)
    (relations game-state player-state permanent card))
  (queryable
    (name winner)
    (relations winner))
  (queryable
    (name mulligan-counter)
    (relations mulligan-decision))
  (queryable
    (name mulligan-yes-counter)
    (relations mulligan-yes-counter)))

