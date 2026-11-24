;;; Rule 121.1  A player draws a card by putting the top card of their library into their hand.

(defrule draw-card-from-library
    ?request <- (draw-request
        (player ?p))

    ?top-card <- (card
        (owner ?p)
        (zone library)
        (card-id ?id)
        (library-position ?pos))

    (not (card
        (owner ?p)
        (zone library)
        (library-position ?pos2&:(< ?pos2 ?pos))
    ))

    =>

    (modify ?top-card
        (zone hand)
        (library-position 0)
    )
    (printout t
        "Player " ?p
        " has drawn a card " ?id
        crlf
    )
    ;; assert che una carta è stata pescata in modo da poter triggherare altre regole
    (assert
        (card-drawn
            (player ?p)
            (card-id ?id)
        )
    )

    ;(assert
    ;    (action-result
    ;        (valid yes)
    ;        (reason "Card drawn.")
    ;    )
    ;)

    (retract ?request)

)

