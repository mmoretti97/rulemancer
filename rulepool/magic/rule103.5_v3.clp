
;; 103.5. Each player draws a number of cards equal to their starting hand size, which is normally seven.
(defrule initial-draw
    ?gs <- (game-state
             (phase initial-draw)
             (active-player ?ap)
             (starting-hand-size ?shs))   
   ?ps <- (player-state (player-id ?p))
   ?progress <- (starting-draw (player-id ?p) (has-drawn no))
   
   =>

    (assert (draw-multiple-cards-request (player ?p) (counter ?shs)))
    (printout t "Cards drawing for player: " ?p crlf)
    (modify ?progress (has-drawn yes))
)

(defrule initial-draw-finished
    ?gs <- (game-state
        (phase initial-draw)
        (starting-hand-size ?shs)
        (active-player ?ap)
    )
       (not
        (starting-draw
            (player-id ?p)
            (has-drawn no)
        )
    )
    
    =>
    (modify ?gs
        (phase mulligan)
        (priority-player ?ap)
    )
    (printout t "Initial draw phase completed. Proceeding to mulligan phase." crlf)
)

;;  A player who is dissatisfied with their initial hand may take a mulligan. First, the starting player declares whether they will take a mulligan.
(defrule mulligan-decision-yes
(declare (salience 5))
    ?ar <- (action-result (valid ?v) (reason ?r))
    ?gs <- (game-state (phase mulligan) (priority-player ?p))
    ?ms <- (mulligan-state (player ?p) (state pending) (yes-counter ?yc))
    ?md <- (mulligan-decision (player ?p) (decision yes))
    =>
    (retract ?md)
    (retract ?ar)
    (modify ?ms (state to-shuffle)(yes-counter (+ ?yc 1)))
    (assert (action-result (valid yes)
        (reason "Mulligan choice accepted")))
)

(defrule mulligan-decision-no
    ?ar <- (action-result (valid ?v) (reason ?r))
    ?gs <- (game-state (phase mulligan) (priority-player ?p))
    ?ms <- (mulligan-state (player ?p) (state pending))
    ?md <- (mulligan-decision (player ?p) (decision no))

    =>
    (retract ?ar)
    (retract ?md)
    (modify ?ms (state end))
    (assert (action-result (valid yes)
        (reason "Mulligan choice declined")))
)

;; A player can take mulligans until their opening hand would be zero cards, after which they may not take further mulligans
(defrule mulligan-limit-reached
    (declare (salience 10))
   (game-state (phase mulligan) (priority-player ?pp)
               (starting-hand-size ?shs))
   ?md <- (mulligan-decision (player ?pp) (decision yes))
   (mulligan-state (player ?pp)(yes-counter ?c&:(>= ?c (- ?shs 1))))
   ?ar <- (action-result (valid ?v) (reason ?r))
   =>
   (retract ?ar ?md)
   (assert (action-result (valid no)
      (reason "Illegal Mulligan: the next hand would be empty.")))
)

;; Once each player has made a declaration, all players who decided to take mulligans do so at the same time.
(defrule mulligan-declarations-complete
   ?gs <- (game-state (phase mulligan)(priority-player ?pp)(active-player ?ap))
   (not (mulligan-state (state pending)))
   =>
   (modify ?gs (phase mulligan-shuffle))
   (modify ?gs (priority-player ?ap))
)

;; To take a mulligan, a player shuffles the cards in their hand back into their library, draws a new hand of cards equal to their starting hand size
(defrule mulligan-shuffle-back
    ?gs <- (game-state (phase mulligan-shuffle)(starting-hand-size ?shs))
    ?ms <- (mulligan-state (player ?p) (state to-shuffle))
    =>
    (put-all-cards-from-hand-to-library ?p)
    (printout t "Player " ?p " is shuffling back their hand into the library." crlf)
    (shuffle-library ?p)
    ;(bind ?iv 0)
    (assert (draw-multiple-cards-request (player ?p)(counter ?shs)))
    (modify ?ms (state to-put-back))
)

(defrule mulligan-shuffle-to-put-back
    ?gs <- (game-state (phase mulligan-shuffle))
    (not (mulligan-state (state to-shuffle))
    )
    =>
    (printout t "Players have shuffled their hands back into the library and are now putting cards back." crlf)
    (modify ?gs (phase mulligan-put-back))
)

;; then puts a number of those cards equal to the number of times that player has taken a mulligan on the bottom of their library in any order
(defrule mulligan-put-back-card
    ?gs <- (game-state
        (phase mulligan-put-back)
    )

    ?ms <- (mulligan-state
        (player ?pp)
        (state to-put-back)
        (yes-counter ?c)
        (cards-puts-back ?cb)
    )

    ?ar <- (action-result
        (valid ?v)
        (reason ?r)
    )

    ?mcd <- (mulligan-cards-back-on-library
        (player ?pp)
        (cards ?card)
    )

    ?hand-card <- (card
        (owner ?pp)
        (zone hand)
        (card-id ?card)
    )

    ?bottom-card <- (card
        (owner ?pp)
        (zone library)
        (library-position ?max-pos)
    )

    (not (card
        (owner ?pp)
        (zone library)
        (library-position ?pos&:(> ?pos ?max-pos))
    ))

    =>

    (retract ?ar)

    (printout t
        "Player " ?pp
        " is putting card " ?card
        " back on the library."
        crlf
    )

    ;; Sposta la carta indicata in fondo al mazzo
    (modify ?hand-card
        (zone library)
        (library-position (+ ?max-pos 1))
    )

    (retract ?mcd)
    (modify ?ms (cards-puts-back (+ ?cb 1)))
    (assert
        (action-result
            (valid yes)
            (reason "Card put back on library.")
        )
    )
    
)

;; Regola per fare catch nel caso in cui la carta dichiarata non si trovi in mano

(defrule mulligan-put-back-card-invalid

    ?ms <- (mulligan-state
        (player ?pp)
        (state to-put-back)
    )

    ?ar <- (action-result
        (valid yes)
    )

    ?mcd <- (mulligan-cards-back-on-library
        (player ?pp)
        (cards ?card)
    )

    ;; Non esiste la carta dichiarata nella mano del giocatore
    (not (card
        (owner ?pp)
        (zone hand)
        (card-id ?card)
    ))

    =>

    (printout t
        "ERROR: Player " ?pp
        " cannot put back card " ?card
        ": card not found in hand."
        crlf
    )

    (retract ?ar)
    (retract ?mcd)

    (assert
        (action-result
            (valid no)
            (reason "Card not found in player's hand.")
        )
    )
)

(defrule mulligan-player-finished-put-back-cards
    ?gs <- (game-state
        (phase mulligan-put-back)
        (starting-hand-size ?shs)
    )

    ?ms <- (mulligan-state
        (player ?pp)
        (state to-put-back)
        (yes-counter ?c)
        (cards-puts-back ?cb)
    )

    (test (= ?cb ?c))
  
    =>
    (modify ?ms (state pending))
    (modify ?ms (cards-puts-back 0))
            (printout t
                "Player " ?pp
                " has finished putting back cards on the library."
                crlf
            )
)

(defrule mulligan-put-back-completed
    ?gs <- (game-state
        (phase mulligan-put-back)
        (active-player ?ap)
        (priority-player ?pp))

    (not
        (mulligan-state
            (state to-put-back)
        )
    )
    =>
    (printout t "Mulligan cards back terminate. Returning to mulligan phase." crlf)

    (modify ?gs (phase mulligan))
    (modify ?gs (priority-player ?ap))
)

;;  This process is then repeated until no player takes a mulligan
(defrule mulligan-completed
    ?gs <- (game-state (phase mulligan)(active-player ?ap)(priority-player ?pp))

    (mulligan-state (player ?p) (state end))
    (not (mulligan-state (state ~end)))
    =>
    (modify ?gs (phase draw)(priority-player ?ap))
    (printout t "Mulligan phase completed. Proceeding to the next phase." crlf)
)

;;; alcuni catch nel caso di asserzioni non valide

(defrule wrong-priority-player
    ?gs <- (game-state (priority-player ?p))
    ?md <- (mulligan-decision (player ?pp) )
    ?ar <- (action-result (valid ?v) (reason ?r))
    (test (neq ?p ?pp))
    =>
    (retract ?md)
    (retract ?ar)
    (printout t "Wrong priority player: " ?p crlf)
    (assert (action-result (valid no) (reason "Wrong priority player.")))
)
