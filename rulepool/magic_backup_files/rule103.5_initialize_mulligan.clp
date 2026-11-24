(defrule start-mulligan-initialize-state
   ?gs <- (game-state (phase start-mulligan)
                      (starting-hand-size ?shs)
                      (active-player ?ap))

   =>

   (do-for-all-facts ((?ps player-state))
      TRUE
       (bind ?found FALSE)

        (do-for-all-facts ((?ms mulligan-state))
            (eq ?ms:player ?ps:player-id)
            (bind ?found TRUE)
        )

        (if (not ?found)
            then
                (assert
                    (mulligan-state
                        (player ?ps:player-id)
                        (state pending)))
                (assert
                    (mulligan-yes-counter
                        (player ?ps:player-id)
                        (counter 0)))
        )
   )

   (modify ?gs (phase mulligan) (priority-player ?ap))
)