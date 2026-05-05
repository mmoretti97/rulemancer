; rule 103.5 Initial card-draw and mulligan
; single-player game

(defrule initial-draw
    ?ar <- (action-result (valid ?v) (reason ?r))
    ?gs <- (game-state
             (phase initial-draw)
             (active-player ?ap)
             (priority-player ?pp))   
    ?gc <- (game-config (num-players ?np))
     =>

    (retract ?ar)
    (bind ?i 0)

    (while (< ?i 7)
        (bind ?i (+ ?i 1))
        (draw-card ?pp)
    )

    (printout t "Initial draw for player: " ?pp crlf)
    (normalize-library ?pp)

    ; Prossimo giocatore a pescare
    (bind ?next (next-player ?pp ?np))

    
        ;(printout t "Attivo=" ?ap crlf)
        ;(printout t "Priorità=" ?pp crlf)
        ;(printout t "Next=" ?next crlf)
    ; Se torniamo al primo, fine draw
    (if (eq ?next ?ap)
      then
         (modify ?gs
            (phase start-mulligan)
            (priority-player ?ap))

         (assert
            (action-result
               (valid yes)
               (reason "All players completed initial draw, going to mulligan phase.")))
       else
       ; Faccio pescare il player successivo
         (modify ?gs (priority-player ?next))
         (assert (action-result (valid yes) (reason "Next player to draw.")))        
    )
)

(defrule start-mulligan

 ?gs <- (game-state
          (phase start-mulligan))

 ?gc <- (game-config (num-players ?np))

 =>

 (bind ?i 1)
 (printout t "Start mulligan" crlf)
 (while (<= ?i ?np)

   (assert
      (mulligan-state
         (player (sym-cat p ?i))
         (state pending)))
   (assert 
      (mulligan-yes-counter
         (player (sym-cat p ?i))
         (counter 0)))

   (bind ?i (+ ?i 1))
 )

 (modify ?gs (phase mulligan))
)

(defrule mulligan-yes
    ?ar <- (action-result (valid ?v) (reason ?r))
    ?gs <- (game-state
             (phase mulligan)
             (active-player ?ap)
             (priority-player ?pp))   
    ?gc <- (game-config (num-players ?np))
    ?ps <- (player-state (player-id ?pp))
    ?md <- (mulligan-decision (player ?pp) (decision yes))
    ?ms <- (mulligan-state (player ?pp) (state pending))
    ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c))
     =>

    (retract ?ar)

   ; eseguo il mulligan, ovvero rimetto le carte nel mazzo, mischio e ripesco 7 carte
    (put-all-cards-from-hand-to-library ?pp)
    (shuffle-library ?pp)
    (bind ?i 0)

    (while (< ?i 7)
        (bind ?i (+ ?i 1))
        (draw-card ?pp)
    )

    (normalize-library ?pp)

        (printout t "Mulligan yes from:" ?pp crlf)

    (modify ?mc (counter (+ ?c 1)))
    (modify ?ms (state pending))

    (retract ?md)
    ; Prossimo giocatore a decidere se mulligan
    ;; TO-DO devo controllare se tutti hanno finito di fare mulligan (ovvero tutti gli altri giocatori hanno detto no)
    ; per evitare un ciclo inutile
    ; se sì passo alla fase successiva
    (modify ?ps (has-priority no))
    (bind ?next (next-player ?pp ?np))
    (modify ?gs (priority-player ?next))
    (do-for-all-facts
        ((?psn player-state))
        (eq ?psn:player-id ?next)
        (modify ?psn (has-priority yes)))

    (assert (action-result (valid yes) (reason "Mulligan done, move to the next player.")))        

  
)

(defrule mulligan-no

 ?ar <- (action-result (valid ?v) (reason ?r))

 ?gs <- (game-state
          (phase mulligan)
          (priority-player ?pp))

 ?md <- (mulligan-decision
          (player ?pp)
          (decision no))
 ?ms <- (mulligan-state (player ?pp) (state pending))
 ?ps <- (player-state (player-id ?pp))
 ?gc <- (game-config (num-players ?np))

 =>
 (retract ?ar)

 (printout t "Mulligan no from:" ?pp crlf)
 (modify ?ps (has-priority no))
 (retract ?md)
 (modify ?ms (state end))
 (bind ?next (next-player ?pp ?np))
 (modify ?gs (priority-player ?next))
 (do-for-all-facts
    ((?psn player-state))
    (eq ?psn:player-id ?next)
    (modify ?psn (has-priority yes)))
 (assert (action-result (valid yes) (reason "Player decided no to mulligan, move to the next player.")))
)

(defrule skip-mulligan-finalize
 
 ?gs <- (game-state (phase mulligan-finalize))

   ; controllo che nessun giocatore abbia effettuato mulligan
   (not (mulligan-yes-counter (counter ?c&:(> ?c 0))))

=>

   (printout t "No players need to put cards back. Skipping mulligan-finalize." crlf)

   (modify ?gs (turn-number 1))
   (modify ?gs (phase main1))
 
)

;; regola che individua quando il giocatore di priorità nella fase di mulligan ha già terminato la fase ( decision = no),
 ;; quindi passo direttamente al giocatore successivo
 ;; se tutti hanno terminato, passo alla fase di mulligan-finalize
(defrule mulligan-check-next-valid-player

 ?gs <- (game-state
          (phase mulligan)
          (priority-player ?pp)
          (active-player ?ap))

 ?ms <- (mulligan-state (player ?pp) (state end))
 ?ps <- (player-state (player-id ?pp))

 ?psa <- (player-state (player-id ?ap))
 
 ?gc <- (game-config (num-players ?np))

=>

 (modify ?ps (has-priority no))

 ; cerca se esiste qualcuno NON end
 (bind ?found FALSE)
 (bind ?next ?pp)

 (bind ?i 0)

 (while (< ?i ?np) do

    (bind ?next (next-player ?next ?np))

    (if (any-factp
          ((?m mulligan-state))
          (and
             (eq ?m:player ?next)
             (neq ?m:state end)))
       then
       (if (any-factp ((?m mulligan-state))
          (and
             (eq ?m:player ?next)
             (eq ?m:state end)))
       then
            (printout t "Player " ?next " has finished mulligan." crlf)
       else
          (bind ?found TRUE)
          (bind ?i ?np)
       )
    )

    (bind ?i (+ ?i 1))
 )

 (if ?found then

    (printout t "Player " ?pp " has finished. Next valid player: " ?next crlf)

    (modify ?gs (priority-player ?next))

    (do-for-all-facts
       ((?psn player-state))
       (eq ?psn:player-id ?next)
       (modify ?psn (has-priority yes)))

 else
    (modify ?gs (priority-player ?ap))
    (modify ?psa (has-priority yes))
    (printout t "All players finished first mulligan phase, going to mulligan-finalize." crlf)
    (modify ?gs (phase mulligan-finalize))
 )
)

;; Controllo se il giocatore con priorità ha terminato il mulligan mettendo le carte in fondo al mazzo
; se si passo priorità al giocatore successivo
; TO-DO: da aggiornare quando il motore potrà gestire multislot (opzionale, gestisce già correttamente una carta lla volta)
(defrule put-card-on-bottom-of-deck-and-check-next-player
  ?gs <- (game-state (phase mulligan-finalize)(priority-player ?pp)(active-player ?ap))
  ?ps <- (player-state (player-id ?pp))
  ?ar <- (action-result (valid ?v) (reason ?r))
  ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c))
  ?mcd <- (mulligan-cards-back-on-library (player ?pp) (cards ?cards))
  ?gc <- (game-config (num-players ?np))

  => 

   (retract ?ar)

 ;; esegue chiamata alla funzione per mettere carta in fondo al mazzo
 (bind ?function-result (put-a-card-on-bottom-deck ?pp ?cards))
 (retract ?mcd)

 (if (not ?function-result) then 
    (assert
      (action-result
         (valid no)
         (reason "Error: Card not found or not in hand.")))
   (return)
 )
 ;; cancello la richiesta per evitare loop

 (printout t "Counter before: " ?c "for player " ?pp crlf)

 ;; decrementa contatore di carte da reimpilare
 (bind ?newc (- ?c 1))

 (if (< ?newc 0) then (bind ?newc 0)) 

 (modify ?mc (counter ?newc))

 (printout t "Counter after: " ?newc "for player " ?pp crlf)

 ;; controllo globale se tutti a 0
 (bind ?allzero TRUE)

 (do-for-all-facts
   ((?x mulligan-yes-counter))
   TRUE
   (if (> ?x:counter 0)
      then
         (bind ?allzero FALSE)))

 (if ?allzero then

    (printout t "All players finished mulligan finalize." crlf)
    (modify ?ps (has-priority no))
    (modify ?gs (priority-player ?ap))
    (do-for-all-facts
      ((?p player-state))
      TRUE
      (if (eq ?p:player-id ?ap)
         then
            (modify ?p (has-priority yes))
         else
            (modify ?p (has-priority no)))
    )    
    (modify ?gs (turn-number 1))
    (modify ?gs (phase main1))

 else

    ;; se current ancora >0 resta lui
    (if (> ?newc 0) then

        (printout t "Player " ?pp " still needs " ?newc " cards." crlf)

    else

        ;; cerca prossimo player con counter >0
        (bind ?next ?pp)
        (bind ?found FALSE)
        (bind ?i 0)

        (while (< ?i ?np) do

           (bind ?next (next-player ?next ?np))

           (if (any-factp
                ((?z mulligan-yes-counter))
                (and
                   (eq ?z:player ?next)
                   (> ?z:counter 0)))
              then
                 (bind ?found TRUE)
                 (bind ?i ?np))

           (bind ?i (+ ?i 1))
        )

        (if ?found then
           (modify ?gs (priority-player ?next))

           (do-for-all-facts
             ((?ps player-state))
             TRUE
             (if (eq ?ps:player-id ?next)
                then
                   (modify ?ps (has-priority yes))
                else
                   (modify ?ps (has-priority no))))
        )
    )
 )

 (assert (action-result (valid yes) (reason "Processed mulligan card on bottom of deck")))

)


;;;; Condizioni non valide


(defrule invalid-mulligan-player-and-decision
   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase mulligan)
             (priority-player ?pp))

   ?md <- (mulligan-decision
              (player ?req)
              (decision ?d))

   (test (neq ?req ?pp))
   (test (not (or (eq ?d yes) (eq ?d no))))
   
   =>

   (retract ?ar)
   (retract ?md)
   (printout t
      "Player " ?req
      " attempted to make an invalid mulligan decision without priority."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: player has no priority during mulligan phase and invalid decision")))
)


(defrule invalid-mulligan-wrong-phase
   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase ?ph)
             (priority-player ?pp))

   ?md <- (mulligan-decision
              (player ?req)
              (decision ?d))

   (test (not (eq ?ph mulligan)))
   (test (or (eq ?d yes) (eq ?d no)))

   =>
   (retract ?ar)
   (retract ?md)
   (printout t
      "Player " ?req
      " attempted to make a mulligan decision outside of mulligan phase."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: can only make mulligan decision during mulligan phase")))
)

(defrule invalid-mulligan-wrong-phase-and-decision
   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase ?ph)
             (priority-player ?pp))

   ?md <- (mulligan-decision
              (player ?req)
              (decision ?d))

   (test (not (eq ?ph mulligan)))
   (test (not (or (eq ?d yes) (eq ?d no))))

   =>
   (retract ?ar)
   (retract ?md)
   (printout t
      "Player " ?req
      " attempted to make an invalid mulligan decision outside of mulligan phase."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: can only make valid mulligan decision during mulligan phase")))
)

(defrule invalid-player-mulligan-decision

   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase mulligan)
             (priority-player ?pp))

   ?md <- (mulligan-decision
              (player ?req)
              (decision ?d))

   (test (neq ?req ?pp))
   (test (or (eq ?d yes) (eq ?d no)))
   
   =>

   (retract ?ar)
   (retract ?md)
   (printout t
      "Player " ?req
      " attempted to make a mulligan decision without priority."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: player has no priority during mulligan phase")))
)

(defrule invalid-mulligan-decision
?ar <- (action-result (valid ?v) (reason ?r))
?gs <- (game-state
          (phase mulligan)
          (priority-player ?pp))
?md <- (mulligan-decision
          (player ?req)
          (decision ?d))

(test (eq ?req ?pp))
(test (not (or (eq ?d yes) (eq ?d no))))
=>
(retract ?ar)
(retract ?md)
(printout t
   "Player " ?req
   " made an invalid mulligan decision: " ?d
   crlf)
(assert
   (action-result
      (valid no)
      (reason "Invalid action: mulligan decision must be yes or no")))

)


(defrule invalid-player-mulligan-card-on-bottom

   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase mulligan-finalize)
             (priority-player ?pp)
             (active-player ?ap))

   ?mcd <- (mulligan-cards-back-on-library
              (player ?req)
              (cards ?cards))

   ?mc <- (mulligan-yes-counter (player ?req) (counter ?c))

   (test (> ?c 0))
   (test (neq ?req ?pp))
=>

   (retract ?ar)

   (retract ?mcd)

   (printout t
      "Player " ?req
      " attempted to put a card on bottom of library without priority."
      crlf)

   (assert
      (action-result
         (valid no)
         (reason "Invalid action: player has no priority during mulligan-finalize")))
)


(defrule invalid-mulligan-finalize-phase

   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase ?ph)
             (priority-player ?pp)
             (active-player ?ap))

   ?mcd <- (mulligan-cards-back-on-library
              (player ?req)
              (cards ?cards))

   (test (not (eq ?ph mulligan-finalize)))

   =>
   (retract ?ar)
   (retract ?mcd)
   (printout t
      "Player " ?req
      " attempted to put a card on bottom of library outside of mulligan-finalize phase."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: can only put cards on bottom of library during mulligan-finalize phase")))
)

(defrule invalid-mulligan-finalize-player

   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase mulligan-finalize)
             (priority-player ?pp)
             (active-player ?ap))

   ?mcd <- (mulligan-cards-back-on-library
              (player ?req)
              (cards ?cards))

   (test (neq ?req ?pp))

   =>
   (retract ?ar)
   (retract ?mcd)
   (printout t
      "Player " ?req
      " attempted to put a card on bottom of library without priority."
      crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: player has no priority during mulligan-finalize phase")))

)