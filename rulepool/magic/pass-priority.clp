(defrule pass-priority
    ?gs <- (game-state (priority-player ?p))
    ?pass <- (pass-priority (player ?pp))
    ?ar <- (action-result (valid ?v) (reason ?r))
    =>
    (modify ?gs (priority-player ?pp))
    (retract ?pass)
    (retract ?ar)
    (printout t "Passing priority to player: " ?pp crlf)
    (assert (action-result (valid yes) (reason "Priority passed to the next player.")))
)