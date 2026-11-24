; rule 103.1 Starting the game
; rule 103.3 Deck shuffle
; 103.1. At the start of a game, the players determine which one of them will choose who takes the first turn. 
; In the first game of a match (including a single-game match), the players may use any mutually agreeable method 
; (flipping a coin, rolling dice, etc.) to do so. In a match of several games, the loser of the previous game chooses who takes the first turn. 
; If the previous game was a draw, the player who made the choice in that game makes the choice in this game.
; The player chosen to take the first turn is the starting player. The game's default turn order begins with the starting player and proceeds clockwise.


; "At the start of a game, the players determine which one of them will choose who takes the first turn"
;TO-DO: per adesso viene selezionato casualmente il giocatore iniziale, si potrebbe implementare un lancio della moneta o simile
; Mirko Mariotti deve implementare un metodo dal motore che generi un seed casuale
(defrule select-initial-player
    ?gc <- (game-config (num-players ?np))
    ?gs <- (game-state
             (phase start-game-players))   
     =>

    (init-random-seed)
    ;(seed (+ 1001 (random 1 99999))) ; da modificare per randomizzare il seed
    (bind ?rnd (random 1 ?np))
    (bind ?player (sym-cat p ?rnd))
    (printout t
      "Starting player chosen: "
      ?player crlf)

    (modify ?gs (active-player ?player)(priority-player ?player))
    (modify ?gs (phase start-game-shuffle))
    ;(assert (action-result (valid yes) (reason "Starting player choosed.")))
)

; "103.3. After the starting player has been determined and any additional steps performed, each player shuffles their deck so that the cards are in a random order."
(defrule initial-shuffle

  ?gc <- (game-config (num-players ?np))
  ?gs <- (game-state (phase start-game-shuffle)(active-player ?ap)(priority-player ?pp))
  
  => 
  (shuffle-library ?pp)
  (printout t "Player " ?pp " has shuffled his library." crlf)
  (bind ?next-player (next-player ?pp ?np))
  (if (eq ?next-player ?ap)
     then
        (printout t "All players have shuffled their libraries." crlf)
        (modify ?gs (phase initial-draw))
        (modify ?gs (priority-player ?ap))
     else
        (printout t "Next player to shuffle: " ?next-player crlf)
        (modify ?gs (priority-player ?next-player)))

)

