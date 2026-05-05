# Moretti Mirco, A.A. 2025-2026

# Regole, template e funzioni implementate


# 5/05/2026

- modifiche ai template
    - **mulligan-decision**: assume solo valori "yes" o "no", adesso è l'azione che compie il giocatore
    - **mulligan-state**: nuovo template che rappresenta invece lo stato nella fase di mulligan per ogni giocatore: assume "pending" o "end" in base alla decisione del giocatore di terminare o proseguire con il mulligan
- modifiche alle regole
    - **103.5**: aggiornata la regola per distinguere tra lo stato del mulligan e la "decision" (asserzione) del giocatore come indicato sopra

# 30/04/2026

- modifiche ai template
    - **mulligan-yes-counter**: template che tiene il conto delle volte che ciascun giocatore ha detto "yes" al mulligan, in modo da sapere quante carte rimettere in fondo al mazzo
    - **mulligan-decision**: rimosso contatore, sostituito dal template di cui sopra
    - **cards**: aggiunto il campo "color" per il colore/mana della carta

- modifiche a fatti
    - aggiunto ad ogni carta dei mazzi rosso e verde di test il rispettivo colore di mana

- nuove funzioni / modifiche
    - **count-all-player-cards**: conta il numero totale di carte di un giocatore (tutte le zone)
    - **shuffle-library**: corretto il comportamento della funzione, ora riassegna correttamente tutte le posizioni da 1 al numero di carte nel mazzo
    - **normalize-library**: funzione utile per ricompattare i valori interi di "library-position" da 1 a N, siccome ogni volta che pesco estraggo la carta/le carte con il valore minore e si creano dei "buchi" tra le posizioni
    - **put-a-card-on-bottom-deck**: funzione che inserisce una carta, in base all'id, in fondo al mazzo. Prima cerca la carta nel mazzo per quel giocatore, e se la trova, cambia la zona in "library" e le assegna "library-position" uguale al numero di carte nel mazzo + 1 (considerando che il mazzo è stato precedentemente normalizzato), poi restituisce TRUE. Se la carta non è trovata la funzione restituisce FALSE

- nuovi fatti
    - query
        - **mulligan-yes-counter**: stampa quante carte ciascun giocatore deve ancora rimettere in fondo al mazzo

- **Modifiche regole**
    - 103.5
        - **initial-draw**: aggiunta normalizzazione del mazzo per ogni giocatore dopo la pesca
        - **start-mulligan**: inizializzati i fatti per la "decision" del mulligan e il contatore per ogni giocatore a 0
        - **mulligan-yes**: aggiunta normalizzazione dopo aver ripescato le carte, aggiorno la "mulligan-decision" e il contatore di mulligan
        - **mulligan-check-next-valid-player**: questa funzione passa la priorità nella fase di mulligan al giocatore successivo, controllando che questo non abbia già terminato il mulligan (decision=end): se trova almeno un giocatore "pending", gli passa la priorità, altrimenti tutti hann terminato la prima fase di mulligan e si passa alla fase di mulligan-finalize
        - **skip-mulligan-finalize**: aggiunta regola che controlla se saltare la fase di "mulligan-finalize" nel caso in cui nessun giocatore avesse effettuato mulligan
        - **put-card-on-bottom-of-deck-and-check-next-player**: funzione che viene attivata nella fase di mulligan-finalize, quando un giocatore fa l'asserzione di scegliere una carta da mettere in fondo al mazzo. Viene chiamata la funzione "put-a-card-on-bottom-deck": con esito negativo, si interrompe la funzione e si ritratta l'asserzione; con esito positivo, si aggiorna il contatore delle carte da rimettere in fondo al mazzo. Poi si controlla se tutti i giocatori hanno terminato il mulligan-finalize e in quel caso si passa alla fase di main, Altrimenti, se il giocatore attuale ha ancora carte da rimpilare la priorità resta a lui, altrimenti viene passata al primo giocatore che ha ancora carte da rimpilare.
        - **invalid-mulligan-player-and-decision**: regola che controlla, nella fase di mulligan, quando un giocatore è sia senza priorità sia effettua una decisione non idonea (yes | no). Ritratta la decisione e interrompe l'operazione.
        - **invalid-mulligan-wrong-phase-and-decision**: intercetta il caso in cui un giocatore voglia asserire una mulligan fuori dalla fase di mulligan assieme ad una decision non valida ( yes | no )
        - **invalid-mulligan-wrong-phase**: intercetta il caso in cui un giocatore tenti di compiere un mulligan fuori dalla fase di mulligan
        - **invalid-player-mulligan-decision**: regola che intercetta quando un giocatore effettua una decisione valida nella fase di mulligan ma non è il giocatore di priorità, ritrattando l'asserzione.
        - **invalid-mulligan-decision**: intercetta quando il giocatore con priorità nella fase di mulligan effettua una decision non valida, ritrattando l'asserzione.
        - **invalid-player-mulligan-card-on-bottom**: intercetta quando nella fase di mulligan-finalize, un giocatore senza priorità tenta di mettere una carta in fondo al mazzo, ritrattando l'asserzione.


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
    - nuova fase iniziale per stabilire chi inizia a giocare "start-game-players"
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


**Nuove regole**

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