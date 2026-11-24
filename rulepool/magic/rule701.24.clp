;;; RULE 701.24 
;701.24a. To shuffle a library or a face-down pile of cards, randomize the cards within it so that no player knows their order.


(defrule shuffle
    ?request <- (shuffle-request
        (player-id ?p))
    =>
    (printout t "Player " ?p " is going to shuffle his library." crlf)
    (shuffle-library ?p)
    (retract ?request)
    
)