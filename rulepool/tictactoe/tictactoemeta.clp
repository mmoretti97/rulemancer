(deffacts tictactoe-config
  (game-config
    (game-name tictactoe)
    (description "A simple Tic Tac Toe game between two players.")
    (num-players 2)))

(deffacts tictactoe-interface
  (assertable
    (name move)
    (relations move))
  (results 
    (name move)
    (relations last-move))
  (queryable
    (name winner)
    (relations winner cell))
  (queryable
    (name cell)
    (relations cell)))