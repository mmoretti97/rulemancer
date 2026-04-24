; rule 103.5 Initial card-draw and mulligan

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
      (mulligan-decision
         (player (sym-cat p ?i))
         (decision pending)
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
    ?md <- (mulligan-decision (player ?pp) (decision yes) (counter ?c))
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
        (printout t "Mulligan yes from:" ?pp crlf)

    (modify ?md (decision pending) (counter (+ ?c 1)))

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
 ?ps <- (player-state (player-id ?pp))
 ?gc <- (game-config (num-players ?np))

 =>
 (retract ?ar)

 (printout t "Mulligan no from:" ?pp crlf)
 (modify ?ps (has-priority no))
 (modify ?md (decision end))
 (bind ?next (next-player ?pp ?np))
 (modify ?gs (priority-player ?next))
 (do-for-all-facts
    ((?psn player-state))
    (eq ?psn:player-id ?next)
    (modify ?psn (has-priority yes)))
 (assert (action-result (valid yes) (reason "Player decided no to mulligan, move to the next player.")))
)

(defrule mulligan-end

 ?gs <- (game-state (phase mulligan))
 ?gc <- (game-config (num-players ?np))

 ;(not (mulligan-decision (decision yes)))
 ;(not (mulligan-decision (decision pending)))

 =>
(bind ?count 0)
   (bind ?i 1)

   (while (<= ?i ?np)

      (bind ?player (sym-cat p ?i))

      (if (any-factp
            ((?md mulligan-decision))
            (and
               (eq ?md:player ?player)
               (eq ?md:decision end)))
         then
            (bind ?count (+ ?count 1))
      )

      (bind ?i (+ ?i 1))
   )

   (if (= ?count ?np)
      then
         (printout t "Mulligan finished for all players." crlf)

         (modify ?gs (phase mulligan-finalize))
   )
)

(defrule finalize-mulligan

 ?gs <- (game-state (phase mulligan-finalize))
 ?ar <- (action-result (valid ?v) (reason ?r))

 =>

 (retract ?ar)

 (do-for-all-facts
   ((?md mulligan-decision))
   TRUE
    (if (eq ?md:decision end)
       then
         ;; TO-DO il giocatore deve scegliere quale carta metter in fondo al mazzo, ora metto a caso
   (put-n-cards-from-hand-to-library
      ?md:player
      ?md:counter))
 )
        (printout t "Mulligan finalize" crlf)

 (modify ?gs (phase main1))
 (assert (action-result (valid yes) (reason "Mulligan terminated, moving to main 1.")))        

)