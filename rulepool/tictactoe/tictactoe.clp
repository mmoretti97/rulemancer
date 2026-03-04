(deftemplate move
  (slot x)
  (slot y)
  (slot player)) ; x | o

(deftemplate last-move
  (slot valid) ; yes | no | none
  (slot reason)) ; description of reason for invalid move

(deftemplate cell
  (slot x)
  (slot y)
  (slot value)) ; x | o

(deftemplate turn
  (slot player)) ; x o o

(deftemplate winner
  (slot player)) ; x | o

(deftemplate state
  (slot phase)) ; playing | ended

(deffacts start
  (turn (player x))
  (state (phase playing))
  (last-move (valid none)))

(deffunction switch-player (?current)
  (if (eq ?current x) then o else x))

(defrule invalid-move-not-playing
  ?s <- (state (phase ended))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  =>
  (retract ?m)
  (assert (last-move (valid no) (reason "Game has ended."))))

(defrule invalid-move-cell-occupied
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?c <- (cell (x ?x) (y ?y) (value ?v))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "Cell is already occupied."))))

(defrule invalid-move-wrong-turn
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?t <- (turn (player ?current&:(not (eq ?current ?p))))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "It's not your turn."))))

(defrule invalid-move-out-of-bounds-x
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x&:(or (< ?x 1) (> ?x 3))) (y ?y) (player ?p))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "X coordinate out of bounds."))))

(defrule invalid-move-out-of-bounds-y
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y&:(or (< ?y 1) (> ?y 3))) (player ?p))
  =>
  (retract ?m)
  (retract ?l)
  (assert (last-move (valid no) (reason "Y coordinate out of bounds."))))

(defrule valid-move
  ?l <- (last-move (valid ?a) (reason ?r))
  ?s <- (state (phase playing))
  ?m <- (move (x ?x) (y ?y) (player ?p))
  ?t <- (turn (player ?p))
  =>
  (retract ?m)
  (retract ?t)
  (retract ?l)
  (assert (turn (player (switch-player ?p))))
  (assert (last-move (valid yes) (reason "Move accepted.")))
  (assert (cell (x ?x) (y ?y) (value ?p))))

;(defrule last-move-show
;  ?l <- (last-move (valid ?v) (reason ?r))
;  =>
;  (printout t "Last move valid: " ?v ", Reason: " ?r crlf))

(defrule win-row
  ?c1 <- (cell (x 1) (y ?y) (value ?p))
  ?c2 <- (cell (x 2) (y ?y) (value ?p))
  ?c3 <- (cell (x 3) (y ?y) (value ?p))
  ?s <- (state (phase playing))
  =>
  (retract ?s)
  (assert (state (phase ended)))
  (assert (winner (player ?p))))

(defrule win-column
  ?c1 <- (cell (x ?x) (y 1) (value ?p))
  ?c2 <- (cell (x ?x) (y 2) (value ?p))
  ?c3 <- (cell (x ?x) (y 3) (value ?p))
  ?s <- (state (phase playing))
  =>
  (retract ?s)
  (assert (state (phase ended)))
  (assert (winner (player ?p))))

(defrule win-diagonal-1
  ?c1 <- (cell (x 1) (y 1) (value ?p))
  ?c2 <- (cell (x 2) (y 2) (value ?p))
  ?c3 <- (cell (x 3) (y 3) (value ?p))
  ?s <- (state (phase playing))
  =>
  (retract ?s)
  (assert (state (phase ended)))
  (assert (winner (player ?p))))

(defrule win-diagonal-2
  ?c1 <- (cell (x 1) (y 3) (value ?p))
  ?c2 <- (cell (x 2) (y 2) (value ?p))
  ?c3 <- (cell (x 3) (y 1) (value ?p))
  ?s <- (state (phase playing))
  =>
  (retract ?s)
  (assert (state (phase ended)))
  (assert (winner (player ?p))))

(defrule draw
  ?c1 <- (cell (x 1) (y 1) (value ?v1))
  ?c2 <- (cell (x 1) (y 2) (value ?v2))
  ?c3 <- (cell (x 1) (y 3) (value ?v3))
  ?c4 <- (cell (x 2) (y 1) (value ?v4))
  ?c5 <- (cell (x 2) (y 2) (value ?v5))
  ?c6 <- (cell (x 2) (y 3) (value ?v6))
  ?c7 <- (cell (x 3) (y 1) (value ?v7))
  ?c8 <- (cell (x 3) (y 2) (value ?v8))
  ?c9 <- (cell (x 3) (y 3) (value ?v9))
  ?s <- (state (phase playing))
  =>
  (retract ?s)
  (assert (state (phase ended)))
  (assert (winner (player draw))))