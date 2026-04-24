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

(deffunction shuffle-library-bkg (?player)

   ;; assegna numero casuale temporaneo
   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library))

      (modify ?c
         (library-position (random 1 100000))))
)

(deffunction shuffle-library (?player)

   ;; conta carte nel mazzo
   (bind ?size (count-library ?player))
   (bind ?pos 1)
   (bind ?remaining ?size)
   ;; reset posizioni
   (do-for-all-facts
      ((?c card))
      (and
         (eq ?c:owner ?player)
         (eq ?c:zone library))
      (modify ?c (library-position 0)))

   ;; assegna 1..N randomicamente
   (while (<= ?pos ?remaining)

      ;; scegli una carta ancora senza posizione
      (bind ?chosen FALSE)

      (while (eq ?chosen FALSE)

         (bind ?rnd (random 1 ?remaining))
         (bind ?i 1)

         (do-for-all-facts
            ((?c card))
            (and
               (eq ?c:owner ?player)
               (eq ?c:zone library)
               (= ?c:library-position 0))

            (if (= ?i ?rnd)
               then
                  (modify ?c (library-position ?pos))
                  (bind ?remaining (- ?remaining 1))
                  (bind ?chosen TRUE))

            (bind ?i (+ ?i 1))
         )
      )

      (bind ?pos (+ ?pos 1))
   )
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