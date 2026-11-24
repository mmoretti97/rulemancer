;;; Rule 121.2. Cards may only be drawn one at a time. If a player is instructed to draw multiple cards, that player performs that many individual card draws.

(defrule draw-multiple-cards
    ?request <- (draw-multiple-cards-request
        (player ?p)
        (counter ?n&:(> ?n 0))
    )

    =>

    (printout t
        "Player " ?p
        " is instructed to draw " ?n " cards."
        crlf
    )

    ;; indico di voler pescare una carta, che attiva la regola di pesca
    (assert
        (draw-request
            (player ?p)
        )
    )

)

(defrule decrement-draw-counter
    ?request <- (draw-multiple-cards-request
        (player ?p)
        (counter ?n&:(> ?n 0))
    )

    ?draw <- (card-drawn
        (player ?p)
    )

    =>
 
    (modify ?request
        (counter (- ?n 1))
    ) 

    (retract ?draw)
)

(defrule draw-multiple-cards-complete
    ?request <- (draw-multiple-cards-request
        (player ?p)
        (counter 0)
    )

    =>
    (printout t
        "Player " ?p
        " has completed drawing all requested cards."
        crlf
    )
    (retract ?request)
)