; "103.3. After the starting player has been determined and any additional steps performed, each player shuffles their deck so that the cards are in a random order."


(defrule initial-shuffle
    ?gs <- (game-state
        (phase start-game-shuffle)
    )
    ;?ar <- (action-result (valid ?v) (reason ?r))
    ?ish <- (initial-shuffle
        (player-id ?p)
        (has-shuffled no)
    )
    =>
    ;(retract ?ar)
    (assert (shuffle-request (player-id ?p)))
    
)

(defrule player-has-shuffled
    ?gs <- (game-state
        (phase start-game-shuffle)
    )
    ;?ar <- (action-result (valid ?v) (reason ?r))
    ?ish <- (initial-shuffle
        (player-id ?p)
        (has-shuffled no)
    )
    ?hs <- (shuffled
        (player-id ?p)
    )
    =>
    ;(retract ?ar)
    (printout t "Player " ?p " has shuffled his library." crlf)
    (modify ?ish (has-shuffled yes))
    (retract ?hs)
)

(defrule all-players-shuffled

    ?gs <- (game-state
        (phase start-game-shuffle))

    (not
        (initial-shuffle
            (has-shuffled no)
        )
    )
    ?ar <- (action-result (valid ?v) (reason ?r))
    =>
    (retract ?ar)
    (printout t "All players have shuffled their libraries." crlf)
    (modify ?gs
        (phase initial-draw)
    )
    (assert
        (action-result
            (valid yes)
            (reason "All players have shuffled their libraries.")
        )
    )
)