(deftemplate player-state
  (slot player-id)     ; p1 | p2
  (slot life)          ; life total (starts at 20)
  (slot max-mana)      ; maximum mana available
  (slot current-mana)  ; mana available this turn, this is simplified and doesn't track colors
  (slot lands-played)  ; number of lands played this turn (max 1)
 ; (slot has-priority)) ; yes | no : RIMOSSO in quanto ridondante
)
(deftemplate card
  (slot card-id)       ; unique card identifier
  (slot name)          ; card name
  (slot type)          ; land | creature | sorcery
  (slot owner)         ; p1 | p2
  (slot zone)          ; deck | hand | battlefield | graveyard
  (slot library-position) ; position of a card in the library
  (slot mana-cost)     ; mana cost (for non-lands)
  (slot color)         ; color of the card 
  (slot power)         ; creature power
  (slot toughness)     ; creature toughness
  (slot damage))       ; damage on creature

; Permanent on battlefield (creatures, lands)
(deftemplate permanent
  (slot card-id)       ; reference to card
  (slot controller)    ; p1 | p2
  (slot tapped)        ; yes | no
  (slot summoning-sick); yes | no (creatures can't attack first turn)
  (slot attacking)     ; yes | no
  (slot blocking)      ; yes | no
  (slot blocked-attacker)) ; card-id of attacker being blocked (or none)

; Game state
(deftemplate game-state
  (slot phase)         ; start-game-players | initial-draw | start-mulligan | mulligan | mulligan-finalize | untap | upkeep | draw | main1 | combat-declare-attackers | 
                       ; combat-declare-blockers | combat-damage | main2 | end | game-over
  (slot turn-number)   ; current turn number
  (slot starting-hand-size) ; number of cards in initial hand (7 usually)
  (slot active-player) ; p1 | p2 - player whose turn it is
  (slot priority-player)) ; p1 | p2 - player with priority

(deftemplate starting-draw
  (slot player-id)     ; p1 | p2
  (slot has-drawn (default no))         ; yes | no 
  )    

(deftemplate initial-shuffle
  (slot player-id)     ; p1 | p2
  (slot has-shuffled (default no)) ; yes | no
)

(deftemplate shuffle-request
  (slot player-id)     ; p1 | p2
)

(deftemplate shuffled
  (slot player-id)     ; p1 | p2
)

(deftemplate draw-request
  (slot player)     ; p1 | p2
)

(deftemplate draw-multiple-cards-request
  (slot player)     ; p1 | p2
  (slot counter(default 1)) ; number of cards to draw
)

(deftemplate card-drawn
  (slot player)     ; p1 | p2
  (slot card-id)    ; 
)

; Actions
(deftemplate mulligan-decision
  (slot player)       ; p1 | p2
  (slot decision))     ; yes | no    

(deftemplate mulligan-state
  (slot player)       ; p1 | p2
  (slot state)      ; pending | to-shuffle | to-put-back | end
  (slot yes-counter (default 0)) ; number of times the player has done mulligan (starts at 0)
  (slot cards-puts-back (default 0)) ; n cards
)      

;(deftemplate mulligan-yes-counter
;  (slot player)       ; p1 | p2
;  (slot counter)    ; number of times the player has mulliganed (starts at 0)
;)     ; temporary counter for drawing cards after mulligan (starts at 0)

(deftemplate mulligan-cards-back-on-library
  (slot player)       ; p1 | p2
  (slot cards)        ; list of card-ids to put back on library
  )

(deftemplate heads-or-tails
  (slot player)       ; p1 | p2
  (slot choice))      ; heads | tails

(deftemplate heads-or-tails-result
  (slot result))      ; heads | tails

(deftemplate starting-player-chooser
    (slot player)
)

(deftemplate starting-player-selection
    (slot chooser)
    (slot player)
)

(deftemplate play-land
  (slot player)
  (slot card-id))

(deftemplate cast-creature
  (slot player)
  (slot card-id))

(deftemplate declare-attacker
  (slot player)
  (slot card-id))

(deftemplate declare-blocker
  (slot player)
  (slot blocker-id)
  (slot attacker-id))

(deftemplate pass-priority
  (slot player))

; The only result with a sort of error code
(deftemplate action-result
  (slot valid)  ; yes | no
  (slot reason)) ; description

; Game final result
(deftemplate winner
  (slot player)) ; p1 | p2 | draw
