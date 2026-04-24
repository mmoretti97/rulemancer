; rule 103.1 Starting the game
; rule 103.3 Deck shuffle

(defrule start-game
    ?ar <- (action-result (valid ?v) (reason ?r))
    ?gc <- (game-config (num-players ?np))
    ?gs <- (game-state
             (phase start-game-players))   
    ?ps <- (player-state(player-id ?player-id))
     =>

    (retract ?ar)
    (init-random-seed)
    ;(seed (+ 1001 (random 1 99999))) ; da modificare per randomizzare il seed
    (bind ?rnd (random 1 ?np))
    (bind ?player (sym-cat p ?rnd))
    (printout t
      "Starting player chosen: "
      ?player crlf)

    (do-for-all-facts
       ((?p player-state))
      TRUE
      (shuffle-library ?p:player-id)
    )

    (modify ?gs (phase initial-draw)(active-player ?player)(priority-player ?player))
    (do-for-all-facts
        ((?ps player-state))
        (eq ?ps:player-id ?player)
        (modify ?ps (has-priority yes)))
    (assert (action-result (valid yes) (reason "Starting player choosed.")))
)
