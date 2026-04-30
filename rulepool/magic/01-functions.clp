(deffunction other-player (?player)
  (if (eq ?player p1) then p2 else p1))


;; TO-DO generalizzare i nomi dei giocatori, da considerare poi il gioco a squadre
(deffunction next-player (?current-player ?number-players)
   (switch ?number-players
      (case 2 then (if (eq ?current-player p1) then p2 else p1))
      (case 3 then (switch ?current-player
                    (case p1 then p2)
                    (case p2 then p3)
                    (case p3 then p1)
                    (default none)))
      (case 4 then (switch ?current-player
                     (case p1 then p2)
                     (case p2 then p3)
                     (case p3 then p4)
                     (case p4 then p1)))
   )
)


(deffunction next-phase (?current-phase)
  (switch ?current-phase
    (case start-game-players then initial-draw)
    (case initial-draw then mulligan)
    (case mulligan then mulligan-finalize)
    (case mulligan-finalize then main1)
    (case untap then upkeep)
    (case upkeep then draw)
    (case draw then main1)
    (case main1 then combat-declare-attackers)
    (case combat-declare-attackers then combat-declare-blockers)
    (case combat-declare-blockers then combat-damage)
    (case combat-damage then main2)
    (case main2 then end)
    (case end then untap)
    (default game-over)))

(deffunction draw-card (?player)

   (bind ?best-card FALSE)
   (bind ?min 999)

   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library))

      (if (< ?c:library-position ?min)
         then
            (bind ?min ?c:library-position)
            (bind ?best-card ?c)))

   (if (neq ?best-card FALSE)
      then
         (modify ?best-card
            (zone hand)
            (library-position 0)))
)

(deffunction put-all-cards-from-hand-to-library (?player)
   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone hand))

      (modify ?c
         (zone library)))
)

(deffunction put-n-cards-from-hand-to-library (?player ?n)
   (bind ?count 0)

   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone hand)
         (< ?count ?n))

      (modify ?c
         (zone library))

      (bind ?count (+ ?count 1)))
)

(deffunction count-library (?player)

   (bind ?n 0)

   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library))
      (bind ?n (+ ?n 1)))

   (return ?n)
)

(deffunction count-all-player-cards (?player)

   (bind ?n 0)

   (do-for-all-facts
      ((?c card))
      (eq ?c:owner ?player)
      (bind ?n (+ ?n 1)))

   (return ?n)
)


(deffunction shuffle-library (?player)

   ;; numero carte nel grimorio
   (bind ?size (count-library ?player))

   (printout t
      "Shuffling library for "
      ?player
      ": "
      ?size
      " cards."
      crlf)

   ;; reset posizioni
   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library))

      (modify ?c (library-position 0))
   )

   ;; assegna posizioni da 1 a N
   (bind ?pos 1)

   (while (<= ?pos ?size)

      ;; quante carte rimangono ancora senza posizione
      (bind ?remaining (- (+ ?size 1) ?pos))
      ;; scegli indice casuale tra quelle libere
      (bind ?rnd (random 1 ?remaining))

      (bind ?i 1)
      (bind ?assigned FALSE)

      (do-for-all-facts
         ((?c card))
         (and
            (eq ?c:owner ?player)
            (eq ?c:zone library)
            (= ?c:library-position 0))

         (if (and (= ?i ?rnd)
                  (eq ?assigned FALSE))
            then
               (modify ?c (library-position ?pos))

               (bind ?assigned TRUE)
         )

         (bind ?i (+ ?i 1))
      )

      (bind ?pos (+ ?pos 1))
   )

   (printout t "Shuffle completed for " ?player crlf)
)

(deffunction normalize-library (?player)

   (bind ?size (count-library ?player))

   (printout t "Normalizing library for " ?player ": " ?size " cards." crlf)

   (bind ?newPos 1)

   (while (<= ?newPos ?size) do

      ;; trova la carta con il minimo library-position
      (bind ?bestFact FALSE)
      (bind ?bestValue 999999)

      (do-for-all-facts
         ((?c card))
         (and
            (eq ?c:owner ?player)
            (eq ?c:zone library)
            (> ?c:library-position 0)
            (< ?c:library-position ?bestValue))

         (bind ?bestFact ?c)
         (bind ?bestValue ?c:library-position)
      )

      (if (neq ?bestFact FALSE) then

         ;; assegna nuova posizione compatta
         (modify ?bestFact (library-position ?newPos))

         ;; escludi temporaneamente questa carta dal prossimo ciclo
         (modify ?bestFact (library-position (+ 100000 ?newPos)))

         (bind ?newPos (+ ?newPos 1))
      )
   )

   ;; rimuove offset 10000 e ripristina i valori da 1 a N
   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library)
         (> ?c:library-position 100000))

      (modify ?c (library-position (- ?c:library-position 100000)))
   )

   (printout t "Library normalized for " ?player crlf)
)


(deffunction put-a-card-on-bottom-deck (?player ?card-id)

   (bind ?number-library-cards (count-library ?player))
   ;(bind ?total-cards (count-all-player-cards ?player))
   (bind ?found FALSE)

   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone hand)
         (eq ?c:card-id ?card-id))

      (modify ?c
         (library-position  (+ ?number-library-cards 1)))
      (modify ?c
         (zone library))  
   (bind ?found TRUE)
 
   )

   (if ?found then
      (return TRUE)
   else
      (printout t "Card " ?card-id " not found in hand for " ?player crlf)
      (return FALSE))

)

;; TO-DO: da rivedere, non è efficiente, ma implemetare nel motore
(deffunction init-random-seed ()

   (bind ?s
      (+ (* (random 1 99999) 17)
         (* (random 1 99999) 31)
         (random 1 99999)))

   (seed ?s)

   (printout t "Random seed initialized: " ?s crlf)
)