(deffacts magic-config
  (game-config
    (game-name magic)
    (description "Magic.")
    (num-players 2)))

(deffacts magic-interface
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
