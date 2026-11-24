;; Rule 121.4. A player who attempts to draw a card from a library with no cards in it loses the game the next time a player would receive priority. 
; (This is a state-based action. See rule 704.)

(defrule draw-from-empty-library
    ?request <- (draw-request
        (player ?p)
        (counter ?n&:(> ?n 0))
    )

    (not (card
        (owner ?p)
        (zone library)
    ))

    =>

    (printout t
        "Player " ?p
        " cannot draw: library is empty."
        crlf
    )

    (retract ?request)
    ;; To-DO asserire una pesca da un mazzo vuoto che comporta la perdita della partita

)