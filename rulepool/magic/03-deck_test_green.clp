;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; TEST DECKS 60 CARTE - FORMATO CASUAL/STANDARD-LIKE
;; Regola usata:
;; - max 4 copie per carta NON basic land
;; - basic land illimitate
;; - solo terre + creature
;; - monocolore
;; - creature semplici ora, estendibili dopo
;;
;; Template usato:
;; (card (card-id ...)
;;       (name ...)
;;       (type ...)
;;       (owner ...)
;;       (zone library)
;;       (mana-cost ...)
;;       (power ...)
;;       (toughness ...)
;;       (damage 0))
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;


;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;; MAZZO P1 - MONO GREEN STOMPY
;; 24 Forest
;; 36 Creature
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(deffacts p1-deck

;; 24 Forest

(card (card-id p1-land1)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land2)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land3)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land4)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land5)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land6)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land7)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land8)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land9)  (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land10) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land11) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land12) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land13) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land14) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land15) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land16) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land17) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land18) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land19) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land20) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land21) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land22) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land23) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))
(card (card-id p1-land24) (name "Forest") (type land) (owner p1) (zone library) (library-position 0) (mana-cost 0) (power 0) (toughness 0) (damage 0))

;; 4 Llanowar Elves
(card (card-id p1-c1) (name "Llanowar Elves") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c2) (name "Llanowar Elves") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c3) (name "Llanowar Elves") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c4) (name "Llanowar Elves") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))

;; 4 Elvish Mystic
(card (card-id p1-c5) (name "Elvish Mystic") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c6) (name "Elvish Mystic") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c7) (name "Elvish Mystic") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))
(card (card-id p1-c8) (name "Elvish Mystic") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 1) (power 1) (toughness 1) (damage 0))

;; 4 Grizzly Bears
(card (card-id p1-c9)  (name "Grizzly Bears") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p1-c10) (name "Grizzly Bears") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p1-c11) (name "Grizzly Bears") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))
(card (card-id p1-c12) (name "Grizzly Bears") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 2) (toughness 2) (damage 0))

;; 4 Kalonian Tusker
(card (card-id p1-c13) (name "Kalonian Tusker") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c14) (name "Kalonian Tusker") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c15) (name "Kalonian Tusker") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c16) (name "Kalonian Tusker") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 2) (power 3) (toughness 3) (damage 0))

;; 4 Centaur Courser
(card (card-id p1-c17) (name "Centaur Courser") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c18) (name "Centaur Courser") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c19) (name "Centaur Courser") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 3) (damage 0))
(card (card-id p1-c20) (name "Centaur Courser") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 3) (power 3) (toughness 3) (damage 0))

;; 4 Rumbling Baloth
(card (card-id p1-c21) (name "Rumbling Baloth") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 4) (power 4) (toughness 4) (damage 0))
(card (card-id p1-c22) (name "Rumbling Baloth") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 4) (power 4) (toughness 4) (damage 0))
(card (card-id p1-c23) (name "Rumbling Baloth") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 4) (power 4) (toughness 4) (damage 0))
(card (card-id p1-c24) (name "Rumbling Baloth") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 4) (power 4) (toughness 4) (damage 0))

;; 4 Craw Wurm
(card (card-id p1-c25) (name "Craw Wurm") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 4) (damage 0))
(card (card-id p1-c26) (name "Craw Wurm") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 4) (damage 0))
(card (card-id p1-c27) (name "Craw Wurm") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 4) (damage 0))
(card (card-id p1-c28) (name "Craw Wurm") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 4) (damage 0))

;; 4 Colossal Dreadmaw
(card (card-id p1-c29) (name "Colossal Dreadmaw") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 6) (damage 0))
(card (card-id p1-c30) (name "Colossal Dreadmaw") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 6) (damage 0))
(card (card-id p1-c31) (name "Colossal Dreadmaw") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 6) (damage 0))
(card (card-id p1-c32) (name "Colossal Dreadmaw") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 6) (power 6) (toughness 6) (damage 0))

;; 4 Gigantosaurus
(card (card-id p1-c33) (name "Gigantosaurus") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 5) (power 10) (toughness 10) (damage 0))
(card (card-id p1-c34) (name "Gigantosaurus") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 5) (power 10) (toughness 10) (damage 0))
(card (card-id p1-c35) (name "Gigantosaurus") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 5) (power 10) (toughness 10) (damage 0))
(card (card-id p1-c36) (name "Gigantosaurus") (type creature) (owner p1) (zone library) (library-position 0) (mana-cost 5) (power 10) (toughness 10) (damage 0))
)