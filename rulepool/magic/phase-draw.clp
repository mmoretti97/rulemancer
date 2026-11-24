; Automatically move to main1
(defrule draw-to-main1
  ?gs <- (game-state (phase draw)(active-player ?ap))
  =>
  (assert (draw-request (player ?ap)))
  (printout t "Draw phase: Player " ?ap " draws a card. Moving to Main Phase 1" crlf)
  (modify ?gs (phase main1))
)
