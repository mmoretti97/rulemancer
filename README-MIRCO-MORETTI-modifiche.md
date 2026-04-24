# Moretti Mirco, A.A. 2025-2026

# Regole e funzioni implementate:

# 24/04/2026

- modifiche ai template
    - **card**: aggiunto "library-position" per indicare la posizione della carta nel mazzo (altrimenti se non è nel mazzo è 0)
    - **mulligan-decision**: template asserzione per indicare se il giocatore vuole effettuare o no il mulligan
    - **mulligan-cards-back-on-library**: il giocatore deve indicare le carte da rimettere in fondo al mazzo al termine del mulligan 
    * TO-DO: scelta multipla carte, da implementare con il motore di gioco

- nuove funzioni
    - **next-player**: utile a stabilire il giocatore successivo, con siwtch-case in caso di più giocatori
    - **next-phase**: indica la fase successiva in base a quella attuale
        - nuove fasi aggiunte: start-game-players, initial-draw, mulligan, mulligan-finalize
    - **draw-card**: funzione di pesca, estrae una carta dalla cima del mazzo del giocatore indicato in base a quella con l'ordinamento minore (library-position)
    - **put-all-cards-from-hand-to-library**: reinserisce tutte le carte dalla mano nel mazzo (ad esempio nella fase di mulligan)
    - **put-n-cards-from-hand-to-library**: mette n carte dalla mano al mazzo (solo per test)
    - **count-library**: conta il numero di carte attualmente nel mazzo e restituisce la cifra
    - **shuffle-library**: mischia il mazzo, riassegnando ad ogni carta una posizioneda 1 a N, dove N è il valore ottenuto dalla funzione 'count-library'
    - **init-random-seed**: inizializza un seed casuale per la funzione di Random 
    * TO-DO: il seed viene generato dal motore

- caso iniziale di test
    - nuova fase iniziale per stabilire chi inizia a giocare
    - nessun giocatore con priorità

- nuovi fatti (02-magicmeta.clp)
    - asserzioni:
        - **mulligan-decision**: decidere se fare mulligan o meno
        - **mulligan-cards-back-on-library**: quali carte vanno rimesse nel mazzo
    - query: 
        - **mulligan-counter**: stampa tutte le scelte di mulligan dei giocatori
        - **only-game-and-player-state**: stampa solo info sullo stato della partita e dei giocatori

- mazzi di test
    - **03-deck_test_green.clp**: mazzo di test di 60 carte bilanciato, con solo creature base e terre
    - **04-deck_test_red.clp**: come mazzo verde
    * TO-DO gestire i colori dei mana, attualmente tipo di mana non presente


### NUOVE REGOLE 

- **rule 103.1**
    - **rule start-game**: implementata regola per gestire le fasi iniziali di gioco:
        - scelta casuale del giocatore iniziale in base al numero di giocatori totali
        - shuffle dei mazzi di ciascun giocatore
        - assegnata priorità al giocatore estratto
        - spostarsi alla fase successiva di pesca iniziale

- **rule 103.5**
    - **rule initial-draw**: a partire dal giocatore con priorità, tutti pescano 7 carte fino a quando non si completa il turno e si passa alla fase preliminare di mulligan
    - **rule start-mulligan**: per ogni giocatore, inizializzo asserzioni con una variabile "decision = pending", che assumerà valori "yes | no" in base alla scelta dei giocatori, poi si passa alla fase di mulligan vera e propria
    - **rule mulligan-yes**: controlla che il giocatore con priorità abbia selezionato "yes" come decisione, allora procede ad inserire le carte nel mazzo e a mischiarle, per poi ripescarne 7. Poi passo la priorità al giocatore successivo 
    * TO-DO: gestire meglio il counter dei "yes", evitare di passare priorità ad un giocatore che ha già detto "no" al mulligan
    - **rule mulligan-no**: se un giocatore dice "no", allora imposto "decision=end" per indicare che ha terminato il mulligan, poi passo la priorità al giocatore successivo 
    * TO-DO come sopra, controllare che il giocatore successivo non abbia terminato
    - **rule mulligan-end**: condizione che controlla quando tutti i giocatori hanno "decision=end", quindi tutti hanno terminato i mulligan, e quindi passo alla fase "mulligan-finalize"
    - **rule finalize-mulligan**: fase finale dove i giocatori con "decision=end" devono rimettere le carte in fondo al mazzo
    * TO-DO, attendere che tutti abbiano rimesso n carte nel mazzo quanto il numero di mulligan, aggiungere asserzione con scelta multipla (modifiche al motore)