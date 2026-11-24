; rule 103.5 Initial card-draw and mulligan


;;; 103.5: "Each player draws a number of cards equal to their starting hand size"
(defrule initialize-player-draw
   ?gs <- (game-state
             (phase initial-draw)
             (priority-player ?pp)
             (starting-hand-size ?shs))
   ?sd <- (starting-draw (player-id ?pp) (counter ?c) (state initialize))
   =>
   (modify ?sd (counter ?shs))
   (modify ?sd (state pending))
   (printout t "Initialize card draw for player: " ?pp crlf)

)

;;; 103.5: "Each player DRAW..."
(defrule initial-draw
    ?gs <- (game-state
             (phase initial-draw)
             (active-player ?ap)
             (priority-player ?pp))   
    ?sd <- (starting-draw (player-id ?pp) (counter ?c&:(> ?c 0)) (state pending))
     =>

    (draw-card ?pp)
    
    (printout t "Card drew for player: " ?pp crlf)
    ;(normalize-library ?pp)

    (modify ?sd (counter (- ?c 1)))
    (printout t "Card counter for player: " ?c crlf)
)

;;; 103.5: "Each player"
(defrule change-player-to-initial-draw
    ?gc <- (game-config (num-players ?np))
    ?gs <- (game-state
             (phase initial-draw)
             (active-player ?ap)
             (priority-player ?pp)
            )   
    ?sd <- (starting-draw (player-id ?pp) (counter 0) (state pending))

     =>

    (printout t "Player " ?pp " finished initial draw." crlf)
      (modify ?sd (state end))
   ; passo al prossimo giocatore
   (bind ?next (next-player ?pp ?np))
   (modify ?gs (priority-player ?next))
   ;; AGGIUNGERE QUI LA MODIFICA DEL PARAMETRO "HAS_PRIORITY" DEL FATTO "PLAYER" PER IL GIOCATORE CON ID = ?NEXT
)

;;; 103.5: "...equal to their starting hand size"
(defrule all-players-finished-initial-draw
    ?gs <- (game-state
             (phase initial-draw)
             (priority-player ?pp)
             (active-player ?ap))
    ?sd <- (starting-draw (player-id ?ap) (counter 0) (state end))
     =>

    (printout t "All players finished initial draw, going to mulligan phase." crlf)
    (modify ?gs (phase start-mulligan))
)

;;; 103.5: "A player who is dissatisfied with their initial hand may take a mulligan."
;;; questa regola inizializza i fatti da utilizzare per tenere traccia dell'andamento del mulligan
; è una regola di supporto, non sono regole di gioco
(defrule start-mulligan

 ?gs <- (game-state
          (phase start-mulligan)
          (starting-hand-size ?shs))

 ?gc <- (game-config (num-players ?np))

 =>

 (bind ?i 1)
 (printout t "Start mulligan" crlf)
 (while (<= ?i ?np)

   (assert
      (mulligan-state
         (player (sym-cat p ?i))
         (state pending)
         (still-to-draw ?shs)
         (has-shuffled no)
         (has-decided no)
         (processed no)))
   (assert 
      (mulligan-yes-counter
         (player (sym-cat p ?i))
         (counter 0)
         (temp-counter 0)))

   (bind ?i (+ ?i 1))
 )

 (modify ?gs (phase mulligan))
)

;;; FASI DEL MULLIGAN
; 1) Il giocatore con priorità dedice se effettuare mulligan o meno
; 2) Poi il giocatore successivo fa lo stesso fino a quando tutti i giocatori non hanno deciso
; 3) appena hanno deciso tutti, i giocatori che hanno deciso di sì effettuano il mulligan, rimescolando le carte e ripescando
; 4) se un giocatore ha deciso di fare mulligan, deve mettere una carta in fondo al mazzo per ogni volta che ha fatto mulligan fino a quel momento
; 5) ogni giocatore in contemporanea ripesca e decide quali carte mettere in fondo al mazzo
; 6) si ripete il ciclo fino a quando tutti hanno deciso di NON effettuarlo


; "A player declared to take a mulligan"
(defrule mulligan-decision-yes
   ?ar <- (action-result (valid ?v) (reason ?r))
   ?gs <- (game-state
             (phase mulligan)
             (priority-player ?pp))
   ?md <- (mulligan-decision (player ?pp) (decision yes))
   ?ms <- (mulligan-state (player ?pp) (state pending)(has-decided no))
   ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c))
   =>
   (retract ?ar)
   (retract ?md)
   (printout t "Player " ?pp " decided to mulligan." crlf)
   (modify ?mc (counter (+ ?c 1)))
   (modify ?mc (temp-counter (+ ?c 1)))
   (modify ?ms (has-decided yes))
   (modify ?gs (phase mulligan-advance-priority))

   (assert (action-result (valid yes) (reason "Mulligan done, move to the next player.")))        

)

; "A player declared NOT to take a mulligan"
(defrule mulligan-decision-no
   ?ar <- (action-result (valid ?v) (reason ?r))
   ?gs <- (game-state
             (phase mulligan)
             (priority-player ?pp))
   ?md <- (mulligan-decision (player ?pp) (decision no))
   ?ms <- (mulligan-state (player ?pp) (state pending)(has-decided no))
   ?ps <- (player-state (player-id ?pp))
   =>
   (retract ?ar)
   (retract ?md)
   (printout t "Player " ?pp " decided not to mulligan." crlf)
   (modify ?ms
      (state end)
      (has-decided yes))
   (modify ?gs (phase mulligan-advance-priority))
   (assert (action-result (valid yes) (reason "The player decided not to mulligan, move to the next player.")))        
   
)

; "Then each other player in turn order does the same"
(defrule mulligan-advance-priority
   ?gs <- (game-state
             (phase mulligan-advance-priority)
             (priority-player ?pp)
             (active-player ?ap))
   ?gc <- (game-config (num-players ?np))
   =>
   (modify ?gs (phase mulligan))
   (modify ?gs (priority-player (next-player ?pp ?np)))

)

; "Once each player has made a declaration, all players who decided to take mulligans do so at the same time"
(defrule mulligan-terminated-first-step
   ?gs <- (game-state
             (phase mulligan))
   ?gc <- (game-config (num-players ?np))

   (exists (mulligan-state (state pending)))
   (not (mulligan-state (has-decided no)(state pending)))
   =>
   (printout t "All players finished first mulligan step, going to mulligan finalize." crlf)
   (modify ?gs (phase mulligan-finalize))
)

; "To take a mulligan, a player shuffles the cards in their hand back into their library,"
(defrule execute-mulligan-put-back-cards
   ?gs <- (game-state
             (phase mulligan-finalize))
   ?ms <- (mulligan-state (player ?p) (state pending)(has-shuffled no)(has-decided yes))
   ?mc <- (mulligan-yes-counter (player ?p) (counter ?c))
   =>
   (printout t "Player " ?p " has to put all cards to the library and shuffle." crlf)
   (put-all-cards-from-hand-to-library ?p)
   (shuffle-library ?p)
   (modify ?ms (has-shuffled yes))
)


(defrule mulligan-go-to-finalize-draw
   ?gs <- (game-state
             (phase mulligan-finalize))
   (not (mulligan-state (state pending)(has-shuffled no)))
   =>
   (printout t "All players have put cards back, going to mulligan finalize draw." crlf)
   (modify ?gs (phase mulligan-finalize-draw))
)

; "draws a new hand of cards equal to their starting hand size"
(defrule mulligan-finalize-draw
   ?gs <- (game-state
             (phase mulligan-finalize-draw))
   ?ms <- (mulligan-state (player ?p) (state pending)(still-to-draw ?sd&:(> ?sd 0))(has-shuffled yes))
   ?mc <- (mulligan-yes-counter (player ?p) (counter ?c&:(> ?c 0)))
   =>
   (printout t "Player " ?p " is drawing a card." crlf)
   (draw-card ?p)
   (normalize-library ?p)
   (modify ?ms (still-to-draw (- ?sd 1)))
)

(defrule mulligan-go-to-put-back-cards
   ?gs <- (game-state
             (phase mulligan-finalize-draw))
   (not (mulligan-state (state pending)(still-to-draw ?sd&:(> ?sd 0))(has-shuffled yes)))
   =>
   (printout t "All players have drawn cards, going to put back cards." crlf)
   (modify ?gs (phase mulligan-finalize-put-back-cards))
)

; "then puts a number of those cards equal to the number of times that player has taken a mulligan on the bottom of their library in any order"
(defrule mulligan-put-back-card
   ?gs <- (game-state
             (phase mulligan-finalize-put-back-cards)
             (active-player ?ap)
             (priority-player ?pp))
   ?ar <- (action-result (valid ?v) (reason ?r))
   ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c)(temp-counter ?tc&:(> ?tc 0)))
   ?mcd <- (mulligan-cards-back-on-library (player ?pp) (cards ?cards))
   =>
   (retract ?ar)
   (printout t "Player " ?pp " is putting a card back on the library." crlf)

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
   
   (modify ?mc (temp-counter (- ?tc 1)))
   (assert (action-result (valid yes) (reason "Card put back on library.")))
   
)

(defrule mulligan-next-player-to-put-back-cards
   ?gs <- (game-state
             (phase mulligan-finalize-put-back-cards)
             (active-player ?ap)
             (priority-player ?pp))
   ?gc <- (game-config (num-players ?np))
   ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c&:(> ?c 0))(temp-counter ?tc&:(= ?tc 0)))
   ?ms <- (mulligan-state (player ?pp) (state pending)(has-shuffled yes)(has-decided yes)(processed no))
   =>
   (printout t "Player " ?pp " has finished putting back cards." crlf)
   (bind ?next (next-player ?pp ?np))
   (modify ?ms (processed yes))
   (modify ?gs (priority-player ?next))
)

; "Salto i giocatori che non devono fare mulligan"
(defrule skip-end-player-mulligan-put-back-cards
   ?gs <- (game-state
             (phase mulligan-finalize-put-back-cards)
             (active-player ?ap)
             (priority-player ?pp))
   ?gc <- (game-config (num-players ?np))
   ?ms <- (mulligan-state (player ?pp) (state end))
   =>
   (printout t "Player " ?pp " has no cards to put back, skipping." crlf)
   (bind ?next (next-player ?pp ?np))
   (modify ?gs (priority-player ?next))
)

; "This process is then repeated until no player takes a mulligan"
(defrule mulligan-all-players-finished-put-back-cards
   ?gs <- (game-state
             (phase mulligan-finalize-put-back-cards)
             (starting-hand-size ?shs)
             (active-player ?ap))
   ?gc <- (game-config (num-players ?np))
   ?myc <- (mulligan-yes-counter)
   ?ms <- (mulligan-state (has-shuffled ?hs)(has-decided ?hd))

   (not (mulligan-state (processed no)(state pending)))
   =>
   (printout t "All players have finished putting back cards, going back to mulligan phase." crlf)
   (modify ?gs (priority-player ?ap))
   (do-for-all-facts
      ((?ms mulligan-state))
      TRUE
      (modify ?ms
         (still-to-draw ?shs)
         (has-shuffled no)
         (has-decided no)
         (processed no)))

   (do-for-all-facts
      ((?mc mulligan-yes-counter))
      TRUE
      (modify ?mc
         (temp-counter 0)))


   (modify ?gs (phase mulligan-check-end))
)

; "no player takes a mulligan"
(defrule mulligan-all-players-end
   ?gs <- (game-state
            (phase mulligan-check-end)
            (active-player ?ap))
   (not (mulligan-state (state pending)))
   =>
   (printout t "All players have finished mulligan, going to main phase 1." crlf)
   (modify ?gs (priority-player ?ap))
   (modify ?gs (phase main1))   
)

(defrule back-to-mulligan
   ?gs <- (game-state
             (phase mulligan-check-end)
             (priority-player ?pp))
   (exists (mulligan-state (state pending)(player ?pp)))
   =>
   (printout t "At least a player decided to mulligan, going back to mulligan phase." crlf)
   (modify ?gs (phase mulligan))
)

; "Salto i giocatori che hanno già risposto di NO al mulligan"
(defrule mulligan-pass-priority-if-state-end
  ?gs <- (game-state
            (phase mulligan-check-end)
            (priority-player ?pp))
  ?gc <- (game-config (num-players ?np))
  ?ms <- (mulligan-state (player ?pp) (state end))
  =>
  (printout t "Player " ?pp " has terminated the mulligan, passing priority." crlf)
  (bind ?next (next-player ?pp ?np))
  (modify ?gs (priority-player ?next))
)

; "Intercetto quando tutti i giocatori hanno scelto NO nella fase di decisione e termino la fase"
(defrule mulligan-phase-end
   ?gs <- (game-state
             (phase mulligan)
             (active-player ?ap))
   (not (mulligan-state (state pending)))
    =>
   (printout t "All players have finished mulligan, going to main phase 1." crlf)
   (modify ?gs (priority-player ?ap))
   (modify ?gs (phase main1))
)

; _________________ "VECCHIE" REGOLE ____________________________


;;;; Condizioni non valide

; "A player can take mulligans until their opening hand would be zero cards, after which they may not take further mulligans"
(defrule mulligan-limit-reached

   ?ar <- (action-result (valid ?v) (reason ?r))

   ?gs <- (game-state
             (phase mulligan)
             (priority-player ?pp)
             (starting-hand-size ?shs))
   ?mc <- (mulligan-yes-counter (player ?pp) (counter ?c&:(>= ?c ?shs)))
   ?md <- (mulligan-decision (player ?pp) (decision yes))
   =>
   (retract ?ar)
   (printout t "Player " ?pp " attempted to take a mulligan but has already reached the limit of " ?shs " cards." crlf)
   (assert
      (action-result
         (valid no)
         (reason "Invalid action: player has already taken the maximum number of mulligans")))
)

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
             (phase mulligan-finalize-put-back-cards)
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

   (test (not (eq ?ph mulligan-finalize-put-back-cards)))

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
             (phase mulligan-finalize-put-back-cards)
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