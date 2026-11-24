(defrule initial-shuffle

  ?gc <- (game-config (num-players ?np))
  ?gs <- (game-state (phase start-game-shuffle)(active-player ?ap)(priority-player ?pp))
  ?ar <- (action-result (valid ?v) (reason ?r))
  => 
  (retract ?ar)
  (shuffle-library ?pp)
  (printout t "Player " ?pp " has shuffled his library." crlf)
  (bind ?next-player (next-player ?pp ?np))
  (if (eq ?next-player ?ap)
     then
        (printout t "All players have shuffled their libraries." crlf)
        (modify ?gs (phase initial-draw))
        (modify ?gs (priority-player ?ap))
     else
        (printout t "Next player to shuffle: " ?next-player crlf)
        (modify ?gs (priority-player ?next-player)))

    (assert
        (action-result
            (valid yes)
            (reason "Player has shuffled library.")
        )
    )
)