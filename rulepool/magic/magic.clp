(deftemplate move
  (slot x)
  (slot y)
  (slot player)) ; x | o

(deftemplate last-action
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
  (last-action (valid none)))

(deffunction switch-player (?current)
  (if (eq ?current x) then o else x))


