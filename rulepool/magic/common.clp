(deftemplate game-config
  (slot game-name)
  (slot description)
  (slot num-players))

(deftemplate assertable
  (slot name)
  (multislot relations))

(deftemplate results
  (slot name)
  (multislot relations))

(deftemplate queryable
  (slot name)
  (multislot relations))