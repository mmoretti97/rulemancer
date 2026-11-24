; rule 103.1 Starting the game
; rule 103.3 Deck shuffle
; 103.1. At the start of a game, the players determine which one of them will choose who takes the first turn. 
; In the first game of a match (including a single-game match), the players may use any mutually agreeable method 
; (flipping a coin, rolling dice, etc.) to do so. In a match of several games, the loser of the previous game chooses who takes the first turn. 
; If the previous game was a draw, the player who made the choice in that game makes the choice in this game.
; The player chosen to take the first turn is the starting player. The game's default turn order begins with the starting player and proceeds clockwise.

;;; NB per adesso è gestito solo il caso di una partita singola

; "At the start of a game, the players determine which one of them will choose who takes the first turn"


; regola per indicare testa o croce nel lancio della moneta
(defrule declare-heads-or-tails
    ?gs <- (game-state
        (phase start-game-players)
    )

    (heads-or-tails
        (player ?pp)
        (choice ?c)
    )

    ?htr <- (heads-or-tails-result
        (result none)
    )

    ?ar <- (action-result
        (valid ?v)
        (reason ?r)
    )

    (test (or (eq ?c heads) (eq ?c tails)))    

    =>

    (retract ?ar)

    (printout t
        "Player " ?pp
        " has chosen: " ?c
        crlf
    )

    (modify ?gs
        (phase start-game-toss)
    )

    (assert
        (action-result
            (valid yes)
            (reason "Player has declared heads or tails.")
        )
    )
)

(defrule toss-coin
    ?gs <- (game-state
        (phase start-game-toss)
    )

    ?htr <- (heads-or-tails-result
        (result none)
    )
 
    =>
    (bind ?rnd (random 1 2))

    (if (= ?rnd 1)
        then
            (printout t "Coin toss result: HEADS" crlf)
            (modify ?htr (result heads))
        else
            (printout t "Coin toss result: TAILS" crlf)
            (modify ?htr (result tails))
    )

    (modify ?gs
        (phase start-game-select-chooser)
    )
)


(defrule select-starting-player-chooser
    ?gs <- (game-state
        (phase start-game-select-chooser)
    )

    (game-config
        (num-players ?np)
    )

    (heads-or-tails
        (player ?p)
        (choice ?c)
    )

    (heads-or-tails-result
        (result ?res&~none)
    )

    (not (starting-player-chooser))

    =>

    (if (eq ?c ?res)
        then
            (assert
                (starting-player-chooser
                    (player ?p)
                )
            )

            (printout t
                "Player " ?p
                " guessed correctly and chooses who starts."
                crlf
            )
        else
            ;; NB Per una partita a due giocatori,
            ;; la scelta passa all'altro giocatore.
            (bind ?other (next-player ?p ?np))

            (assert
                (starting-player-chooser
                    (player ?other)
                )
            )

            (printout t
                "Player " ?p
                " guessed incorrectly. Player "
                ?other
                " chooses who starts."
                crlf
            )
    )

    (modify ?gs
        (phase start-game-choose-player)
    )
)

(defrule select-initial-player
    ?gs <- (game-state
        (phase start-game-choose-player)
    )
    ?ar <- (action-result (valid ?v) (reason ?r))

    (starting-player-chooser
        (player ?chooser)
    )

    ?sel <- (starting-player-selection
        (chooser ?chooser)
        (player ?initial)
    )
    ?pl <- (player-state
        (player-id ?initial)
    )

    =>
    (retract ?ar)
    (printout t
        "Player " ?chooser
        " selected player " ?initial
        " as the starting player."
        crlf
    )

    (modify ?gs
        (active-player ?initial)
        (priority-player ?initial)
        (phase start-game-shuffle)
    )

    (assert
        (action-result
            (valid yes)
            (reason "Starting player selected.")
        )
    )
    (retract ?sel)
)

;; Regola per intercettare un giocatore diverso da quello del lancio vincente che tenta di decidere chi inizia
(defrule select-initial-player-invalid
    ?gs <- (game-state
        (phase start-game-choose-player)
    )
    ?ar <- (action-result (valid ?v) (reason ?r))

    (starting-player-chooser
        (player ?chooser)
    )

    ?sel <- (starting-player-selection
        (chooser ?other&~?chooser)
        (player ?initial)
    )

    =>

    (retract ?ar)

    (printout t
        "Player " ?other
        " attempted to select the starting player, but "
        ?chooser
        " is the chooser."
        crlf
    )
    (retract ?sel)
    (assert
        (action-result
            (valid no)
            (reason "Invalid starting player selection. Not the chooser's turn.")
        )
    )
    
)

(defrule select-initial-player-player-not-found
    ?gs <- (game-state
        (phase start-game-choose-player)
    )
    ?ar <- (action-result (valid ?v) (reason ?r))

    (starting-player-chooser
        (player ?chooser)
    )
    ?sel <- (starting-player-selection
        (chooser ?chooser)
        (player ?pl)
    )
    (not (player-state (player-id ?pl)))
    

    =>

    (retract ?ar)

    (printout t
        "Player " ?chooser
        " attempted to select an invalid player "
        ?pl
        " as the starting player."
        crlf
    )
    (retract ?sel)
    (assert
        (action-result
            (valid no)
            (reason "Invalid starting player selection. Player does not exist.")
        )
    )
    
)

(defrule invalid-coin-toss-choice
    ?gs <- (game-state
        (phase start-game-players)
    )

    ?ar <- (action-result
        (valid ?v)
        (reason ?r)
    )

    ?htr <- (heads-or-tails
        (player ?p)
        (choice ?c)
    )

    (test (not (or (eq ?c heads) (eq ?c tails))))

    =>

    (retract ?ar)
    (retract ?htr)
    (printout t
        "Player " ?p
        " made an invalid coin toss choice: "
        ?c
        crlf
    )

    (assert
        (action-result
            (valid no)
            (reason "Invalid coin toss choice. Must be 'heads' or 'tails'.")
        )
    )
)