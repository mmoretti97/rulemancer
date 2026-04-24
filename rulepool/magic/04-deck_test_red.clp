;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; MAZZO P2 - MONO RED AGGRO / CREATURES ONLY
;; 24 Mountain
;; 36 Creature
;; Regola:
;; - max 4 copie non basic land
;; - solo terre + creature
;; - nomi reali Magic
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(deffacts p2-deck

;; 24 Mountain

(card (card-id p2-land1)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land2)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land3)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land4)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land5)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land6)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land7)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land8)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land9)  (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land10) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land11) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land12) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land13) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land14) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land15) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land16) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land17) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land18) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land19) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land20) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land21) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land22) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land23) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p2-land24) (name "Mountain") (type land) (owner p2) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))

;; 4 Goblin Arsonist
(card (card-id p2-c1) (name "Goblin Arsonist") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p2-c2) (name "Goblin Arsonist") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p2-c3) (name "Goblin Arsonist") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p2-c4) (name "Goblin Arsonist") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))

;; 4 Goblin Piker
(card (card-id p2-c5) (name "Goblin Piker") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 1) (damage 0))
(card (card-id p2-c6) (name "Goblin Piker") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 1) (damage 0))
(card (card-id p2-c7) (name "Goblin Piker") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 1) (damage 0))
(card (card-id p2-c8) (name "Goblin Piker") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 1) (damage 0))

;; 4 Borderland Marauder
(card (card-id p2-c9)  (name "Borderland Marauder") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 1) (toughness 2) (damage 0))
(card (card-id p2-c10) (name "Borderland Marauder") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 1) (toughness 2) (damage 0))
(card (card-id p2-c11) (name "Borderland Marauder") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 1) (toughness 2) (damage 0))
(card (card-id p2-c12) (name "Borderland Marauder") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 1) (toughness 2) (damage 0))

;; 4 Ember Hauler
(card (card-id p2-c13) (name "Ember Hauler") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p2-c14) (name "Ember Hauler") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p2-c15) (name "Ember Hauler") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p2-c16) (name "Ember Hauler") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))

;; 4 Ahn-Crop Crasher
(card (card-id p2-c17) (name "Ahn-Crop Crasher") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 2) (damage 0))
(card (card-id p2-c18) (name "Ahn-Crop Crasher") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 2) (damage 0))
(card (card-id p2-c19) (name "Ahn-Crop Crasher") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 2) (damage 0))
(card (card-id p2-c20) (name "Ahn-Crop Crasher") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 2) (damage 0))

;; 4 Fire Elemental
(card (card-id p2-c21) (name "Fire Elemental") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 4) (damage 0))
(card (card-id p2-c22) (name "Fire Elemental") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 4) (damage 0))
(card (card-id p2-c23) (name "Fire Elemental") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 4) (damage 0))
(card (card-id p2-c24) (name "Fire Elemental") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 4) (damage 0))

;; 4 Charging Monstrosaur
(card (card-id p2-c25) (name "Charging Monstrosaur") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c26) (name "Charging Monstrosaur") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c27) (name "Charging Monstrosaur") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c28) (name "Charging Monstrosaur") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 5) (power 5) (toughness 5) (damage 0))

;; 4 Shivan Dragon
(card (card-id p2-c29) (name "Shivan Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c30) (name "Shivan Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c31) (name "Shivan Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 5) (toughness 5) (damage 0))
(card (card-id p2-c32) (name "Shivan Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 5) (toughness 5) (damage 0))

;; 4 Volcanic Dragon
(card (card-id p2-c33) (name "Volcanic Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 4) (toughness 4) (damage 0))
(card (card-id p2-c34) (name "Volcanic Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 4) (toughness 4) (damage 0))
(card (card-id p2-c35) (name "Volcanic Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 4) (toughness 4) (damage 0))
(card (card-id p2-c36) (name "Volcanic Dragon") (type creature) (owner p2) (zone library) (library-position 0) (mana-cost 6) (power 4) (toughness 4) (damage 0))
)