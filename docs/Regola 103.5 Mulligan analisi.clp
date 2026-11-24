; Analisi per lo sviluppo della regola 103.5 relativa al mulligan

;1) Ogni giocatore pesca un numero di carte = starting hand size
;2) Un giocatore non contento delle proprie carte può fare mulligan
;3) Il giocatore iniziale dichiara se vuole fare mulligan o meno
;4) Di seguito, ogni giocatore fa la stessa scelta
;5) Quando ogni giocatore ha effettuato la scelta, tutti i giocatori che hanno deciso di fare mulligan lo eseguono allo stesso momento
;6) per eseguire mulligan, un giocatore mischia le carte nel mazzo, ripesca le carte = hand-size e poi procede a mettere carte in fondo al mazzo in base al numero di mulligan
;7) Un giocatore che non effettua mulligan si tiene quelle carte come mano iniziale e non può effettuare più mulligan
;8) il processo continua fino a quando tutti decidono di non fare mulligan
;9) il limite di mulligan è uguale al numero di starting-hand-size (non può essere ovviamente meno di 0)

; IMPLEMENTAZIONI CONCETTUALI

;1) ogni giocatore pesca in contemporanea la mano iniziale, ci si deve trovare in una fase iniziale di pesca con la mano vuota
;2), 3) a partire dal giocatore con priorità, si può scegliere se effettuare mulligan o meno
;    - se si, si aggiorna un contatore/fatto che tiene conto del numero di mulligan effettuati, poi deve passare priorità al giocatore successivo
;    - se no, il giocatore non potrà effettuare più mulligan, e poi deve passare priorità al giocatore successivo
;    . aspetti da considerare: quante volte ha eseguito mulligan, se ha deciso di fare keep, se al giro attuale risponde si/no -> riassumo tutto in un fatto stato unico
;        - se non ha ancora scelto: state pending
;        - se ha scelto di NON fare mulligan: state end
;        - se ha scelto di fare mulligan: state to-shuffle
;    . tenere comunque un contatore di quante volte ha effettuato mulligan
;4) passaggio manuale di priorità
;5) quando si ritorna al giocatore iniziale, si leggono i fatti dei giocatori che hanno effettuato mulligan e si procede a mischiare le carte e a ripescarle
;6) poi, ogni giocatore deve scegliere tot carte da mettere nel mazzo in base al numero di mulligan effettuati
;7) si procede ripartendo dal giocatore iniziale a stabilire chi vuole effettuare mulligan
;8) se tutti hanno detto NO, il mulligan termina e si prosegue col gioco


; IMPLEMENTAZIONE

; ______     1) initial-draw ___________

(defrule initial-draw
    (game-state
        (phase initial-draw)
        (starting-hand-size ?shs)
    )
    (player-state (player-id ?p))
    =>
    (bind ?hand
        (find-all-facts ((?c card))
            (and (eq ?c:owner ?p)
                 (eq ?c:zone hand))
        )
    )
    (if (< (length$ ?hand) ?shs)
        then
            (draw-card ?p)
    )
)

(defrule initial-draw-finished
    ?gs <- (game-state
        (phase initial-draw)
        (starting-hand-size ?shs)
        (active-player ?ap)
    )
    (not
        (and
            (player-state (player-id ?p))
            (test
                (<
                    (length$
                        (find-all-facts ((?c card))
                            (and (eq ?c:owner ?p)
                                 (eq ?c:zone hand))
                        )
                    )
                    ?shs
                )
            )
        )
    )
    =>
    (modify ?gs
        (phase start-mulligan)
        (priority-player ?ap)
    )
)


; 2) ______     2) mulligan-decision ___________

(defrule start-mulligan-initialize-state
   (declare (salience 20))
   ?gs <- (game-state (phase start-mulligan)
                      (starting-hand-size ?shs)
                      (active-player ?ap))

   =>

   (do-for-all-facts ((?ps player-state))
      TRUE
      (if (not (any-factp ((?ms mulligan-state))
                 (eq ?ms:player (fact-slot-value ?ps player-id))))
         then
         (assert (mulligan-state
                    (player (fact-slot-value ?ps player-id))
                    (state pending)))
         (assert (mulligan-yes-counter
                    (player (fact-slot-value ?ps player-id))
                    (counter 0))))
   (modify ?gs (phase mulligan) (priority-player ?ap)))
)

(defrule mulligan-decision-yes
    ?gs <- (game-state (phase mulligan) (priority-player ?p))
    ?ms <- (mulligan-state (player ?p) (state pending))
    ?md <- (mulligan-decision (player ?p) (decision yes))

    =>

    (modify ?ms (state to-shuffle))
    (bind ?counter (fact-slot-value ?ms counter))
    (modify ?ms (counter (+ ?counter 1)))
    (assert (action-result (valid yes)
        (reason "Mulligan choice accepted for player " ?p)))
)

(defrule mulligan-decision-no
    ?gs <- (game-state (phase mulligan) (priority-player ?p))
    ?ms <- (mulligan-state (player ?p) (state pending))
    ?md <- (mulligan-decision (player ?p) (decision no))

    =>

    (modify ?ms (state end))
    (assert (action-result (valid yes)
        (reason "Mulligan choice accepted for player " ?p)))
)

; Regola che impedisce di fare mulligan se il numero di carte da pescare sarebbe 0
(defrule mulligan-limit-reached
   (declare (salience 30))
   (game-state (phase mulligan) (priority-player ?pp)
               (starting-hand-size ?shs))
   ?md <- (mulligan-decision (player ?pp) (decision yes))
   (mulligan-state (player ?pp) (state pending))
   (mulligan-yes-counter (player ?pp) (counter ?c&:(>= ?c (- ?shs 1))))
   ?ar <- (action-result (valid ?v) (reason ?r))
   =>
   (retract ?ar ?md)
   (assert (action-result (valid no)
      (reason "Mulligan illegale: la mano successiva sarebbe vuota.")))
)

(defrule mulligan-declarations-complete
   (declare (salience -10))
   ?gs <- (game-state (phase mulligan))
   (not (mulligan-state (state pending)))
   =>
   (modify ?gs (phase mulligan-shuffle))
)

(defrule mulligan-shuffle-back
    ?gs <- (game-state (phase mulligan-shuffle))
    ?ms <- (mulligan-state (player ?p) (state to-shuffle))
    ?mc <- (mulligan-yes-counter (player ?p) (counter ?c))
    =>
    (put-all-cards-from-hand-to-library ?p)
    (shuffle-library ?p)
    (loop-for-count 1 ?shs do
        (draw-card ?p)
        (normalize-library ?p)
    )
    (modify ?ms (state to-put-back))
)

(defrule mulligan-shuffle-to-put-back
    ?gs <- (game-state (phase mulligan-shuffle))
    (not (mulligan-state (state to-shuffle))
    )
    =>
    (modify ?gs (phase mulligan-put-back))
)

(defrule mulligan-put-back-card
    ?gs <- (game-state (phase mulligan-put-back) (priority-player ?pp))
    ?ms <- (mulligan-state (player ?pp) (state to-put-back))

    ?ar <- (action-result (valid ?v) (reason ?r))
    ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c))
    ?mcd <- (mulligan-cards-back-on-library (player ?pp)  (cards $?cards&:(> (length$ ?cards) 0)))

    => 

    (bind ?card (nth$ 1 ?cards))

    (retract ?ar)
    (printout t "Player " ?pp " is putting a card back on the library." crlf)

        ;; esegue chiamata alla funzione per mettere carta in fondo al mazzo
    (bind ?function-result (put-a-card-on-bottom-deck ?pp ?card))
    (retract ?mcd)

    (if (not ?function-result) then 
        (assert
            (action-result
                (valid no)
                (reason "Error: Card not found or not in hand.")))
        (return)
    )
    (bind ?remaining (delete$ ?cards 1 1))
    (modify ?mcd (cards $?remaining))

    (assert (action-result (valid yes) (reason "Card put back on library.")))
)

(defrule mulligan-put-back-completed
    ?gs <- (game-state
        (phase mulligan-put-back))

    (not
        (mulligan-state
            (state to-put-back)
        )
    )
    =>
    (printout t "Mulligan cards back terminate. Returrning to mulligan phase." crlf)

    (modify ?gs (phase mulligan))
)

(defrule mulligan-completed
    ?gs <- (game-state (phase mulligan))
    (mulligan-state (player ?p) (state end))
    (not (mulligan-state (state ~end)))
    =>
    (modify ?gs (phase main))
    (printout t "Mulligan phase completed. Proceeding to main phase." crlf)
)