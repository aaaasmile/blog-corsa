# Corsa Blog
Un service per gestire il mio blog sulla corsa.
L'intenzione è quella di rimpiazzare completamente il mio sito https://stesosopra.blogspot.com/.

## Blog Corsa Static-Html
La parte testuale che fornisce il blog sulla corsa è situata in static/blog-corsa.
Questa è una subdirectory generata in gran parte automaticamente. Ho iniziato ad usare Zine,
per poi abbandonarlo quando la generazione dell'html è diventata complessa.

La URL di riferimento del blog è  http://localhost:5572. 
Nota che è una URL root. Non può essere in 
un path che non sia / in quanto il codice html generato richiede questo tipo di percorso.

### Creare il contenuto Html
Per creare i post html uso il generatore che si trova in content/src.
Nella sottodirectory Content metto tutti i vari post in una directory singola.
Qui ho il mio file _mdhtml_ e i vari files delle immagini relative al post. 
Il file mdhtml contiene solo la parte all'interno del tag article principale.

L'ho chiamato mdhtml in quanto, come nei files md, c'è una parte di dati in testa seguita da
una parte in html. Il testo dell'articolo lo edito in html per avere la massima flessibilità
di generazione del codice html. La parte di dati mi serve solo per quei campi che hanno bisogno
di un valore eplicito, altrimenti ci sono spazi per delle ambiguità nella generazione del codice. 

Per quanto mi riguarda, usando un html strettamente semantico, non vedo il bisogno di editare il post
in md con tutte le restrizioni del caso. Qui vengono generati i files statici della directory _posts_.

I contenuti della directory _pages_ li creo manualmente.

Il file _index.html_ in blog-corsa avrà una parte generata automaticamante quando aggiorno i post.

### Editare un post
In src/Content si lancia il watcher, che agisce quando il file mdhtml cambia. Oppure viene inserita un'immagine o rinominata. Per esempio se voglio modificare il file 24-11-08-ProssimaGara.mdhtml:

    cd ./content/src
    go run .\main.go -config ..\..\config.toml  -watch -target ..\2024\11\08\
Poi mentre cambio il file 24-11-08-ProssimaGara.mdhtml, mi piazzo col browser su:

    http://localhost:5572/posts/2024/11/24-11-08-ProssimaGara/
per vedere il cambiamenti nell'output statico dopo un browser reload

### Abbandono di un generatore statico standard
Ho abbandonato Zine per diversi motivi. Voglio editare dei files html manualmente, almeno alcuni.
Altri li voglio generare editando, però, solo alcune parti in html e non in md.
Il contenuto dei post voglio metterlo anche in un db per avere la full search out of the box.
Il db mi serve, oltre che per i commenti, anche per inserire dei post in runtime alla fine dell'articolo
che sono correlati all'articolo appena letto. Nel sito blogspot, invece, la navigazione è sempre lineare,
tranne quando l'utente epsplicitamente ne cerca un altro.

Altro tema è la gestione delle immagini. In blogspot, quando eseguo un upload di un'immagine,
ne viene creata una copia delle dimensioni di 320 pixels in larghezza. Questo è anche il compito
del mio generatore. Quando un'immagine va a finire nella directory del mio post, viene creata una
copia delle dimensioni di 320 pixel in larghezza.

Non ho idea se sia possibile usare un generatore come Hugo o Jeckill per la mia applicazione.
Ho trovato più divertente crearne uno mio.

## Formato mdhtml
È un file che ha una sezione per i dati come i files md e una con il contenuto.
Nella parte del contenuto uso il codice html. Per velocizzare la generazione dei tag, uso
un preprocessor che mi genera un codice html. Esso supporta queste macro:

- link
- figstack
- youtube
- latest_posts
- archive_posts

Tutti i comandi sono compresi tra parantesi quadre. La lista la trovo nel file _lexer-builtin-func.go_.

### Link
Il comando _link_ serve per avere un <a href> con il link url uguale al testo mostrato.
Esempio:

    [link 'https://wien-rundumadum-2024-130k.legendstracking.com/']
genera:

    <a href='https://wien-rundumadum-2024-130k.legendstracking.com/'>https://wien-rundumadum-2024-130k.legendstracking.com/ </a> 

### Link caption
Un link che però ha anche la caption.
Esempio:

    [linkcap 'Tracker', 'https://wien-rundumadum-2024-130k.legendstracking.com/']
genera:

    <a href='https://wien-rundumadum-2024-130k.legendstracking.com/'>Tracker</a> 

### figstack
Serve per creare velocemente una galleria di immagini.
Esempio:

    [figstack
        'AustriaBackyardUltra2024011.jpg', 'Partenza mondiale Backyard',
        'backyard_award.png', 'Certificato finale'
    ]
Ogni coppia è rappresentata dal nome del file dell'immagine integrale e dal titolo.
Il codice html generato lo trovo di seguito. Col file dell'immagine integrale 
cosidero per dato il file d'immagine in formato ridotto di larghezza 320 pixel.

### youtube
Genera l'iframe che serve per contenere il video player di youtube.  
Esempio:

    [youtube 'IOP7RhDnLnw'] 
Dove IOP7RhDnLnw è il video ID su youtube.

Per centrare il video come le figure:

    <figure>
      [youtube 'vsC8SXH6Ffg']
      <figcaption>Il video ufficiale della gara</figcaption>
    </figure>
Oppure

    <p class="center">
        [youtube 'vsC8SXH6Ffg']
    </p>

## Immagini (html creato da figstack)
Quondo ho una serie di immagini da inserire nel post, uso il seguente html:

    <section class="vertstack">
      <figure>
        <a href="tabella.png"><img src="tabella_320.png" alt="Tabella finale" /></a>
        <figcaption>Tabella finale</figcaption>
      </figure>
      <figure>
        <a href="partenza.jpg"><img src="partenza_320.jpg" alt="Appena partiti" /></a>
        <figcaption>Appena partiti</figcaption>
      </figure>
    </section>
Per questo ho bisogno delle immagini in formato ridotto _xxx\_320_.
Qui si vede che le immagini sono nella stessa directory del post in quanto non riutilizzo mai
la stessa immagine in un altro post.

## latest_posts
Nella pagina proncipale ho bisogno di un sommario degli ultimi post. Per questo uso la macro:

    [latest_posts 'IgorRun Blog', '7']

Dove 'IgorRun Blog' rappresenta il titolo e '7' è il numero dei post da mettere.
Il risultato è un html con la lista degli ultimi 7 post. L'elenco viene creato leggendo il
database.

## Commenti
I commenti sono parte integrante dei post. Siccome i songoli post sono creati staticamente,
i commenti vengono mostrati tramite htmx in fase di rendering sul browser.

### Flow del nuovo commento
Quando viene postato un commento sul blog, esso potrebbe venire dapprima esaminato con il
service [askimet](https://akismet.com/plan/personal/) per vedere se è uno spam.
Al momento questo non avviene. Il commento, che non contiene html, viene messo in moderazione.
Viene mandata una notifica via mail e/o telegram per approvare il commento attraverso la pagina di admin.

### Rendering dei commenti di un Post
Uso htmx per avere il fetch delle parti dinamiche, come per esempio i commenti di un post, nella
parte statica html. Nota che in questi casi i servizi come _commento_, che si trova su [gitlab.com/commento](https://gitlab.com/commento/commento), usa un approccio differente. Vale a dire i commenti vengono
visualizzati usanto una request via javascript al server di Commento, il quale risponde con un json.
Il file javascript sul client crea poi on-demand il codice html che viene aggiunto all'elemento div indicato.
Commento.js è un file di 60k non min.  

Quando uso htmx, invece, il server fornisce la parte html già creata senza passare da json. 

### Protezione Spam
Quando avevo un Guestbook gestito attraverso un Form, esso era bersaglio di
scraping che automaticamente mandavano dei messaggi. 
La mia protezione è il Form posizionato all'interno del tag html _details_.
Il form compare quando l'utente apre il tag details attraverso htmx.
L'altro step è quello della moderazione e dell'impossibilità di inviare html.

## Dashboard Admin
La parte che riguarda l'amministrazione del blog è gestita con vue in modalità single page.
Per contro, la parte testuale dei vari post è generata staticamente in html.
Al momento la Dashboard non gestisce i post, solo i commenti. 
Successivamente potrebbe essere usata per
creare anche dei nuovi contenuti. Questo vorrebbe dire gestire la generazione statica di html.
La URL di riferimento è: http://localhost:5572/blog-admin/

### TODO
- Nei commenti va implementata la risposta, per avere commenti nei livelli inferiori [DONE]
- La pagina admin deve essere protetta da un token di sign-in [DONE] 
- Nella pagina admin, manca la gestione Edit/delete/approve/decline dei commenti [DONE]
- buildmain dovrebbe creare: main, archivio e feed.

### Stop del service
Per stoppare il sevice si usa:

    sudo systemctl stop igorrun

## Deployment su ubuntu direttamente

    cd ~/build/blog-corsa
    git pull --all
    ./publish-service.sh
Oppure uso Visual Code in remoto dove uso il synch di git. Qui nel terminal mi basta usare:

    ./publish-service.sh

## Service Config
Per prima cosa va creato il file igorrun.service.
Il contenuto l'ho messo sotto in una sezione apposita.

    sudo nano /lib/systemd/system/igorrun.service
Poi si fa l'enable:

    sudo systemctl enable igorrun.service
E infine lo start:

    sudo systemctl start igorrun
Logs sono disponibili con:

    sudo journalctl -f -u igorrun

## igorrun.service
Qui segue il contenuto del file igorrun.service
Nota il Type=idle che è meglio di simple. Così 
viene fatto partire dopo che la wlan ha ottenuto l'IP intranet
e così si ha l'accesso.

```
[Install]
WantedBy=multi-user.target

[Unit]
Description=igorrun service
ConditionPathExists=/home/igor/app/go/igorrun/current/blog-corsa.bin
After=network.target

[Service]
Type=idle
User=igor
Group=igor
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/igor/app/go/igorrun/current/
ExecStart=/home/igor/app/go/igorrun/current/blog-corsa.bin

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/igorrun
ExecStartPre=/bin/chown igor:igor /var/log/igorrun
ExecStartPre=/bin/chmod 755 /var/log/igorrun
StandardOutput=syslog
StandardError=syslog
```

### config_custom.toml
È il file che mi esegue un ovveride del file config.toml. 
Mi serve in quanto config.toml si trova su gitHub, mentre config_custom.toml è
solo locale fuori da git. Si trova in:

    /home/igor/app/go/igorrun/current/

## Visual Code
Per lo sviluppo iniziale ho usato windows, poi, per l'update del service,
ho usato Visual Code Remote nella directory ~/build/blog-corsa.

## nginx proxy
Vedi il documento locale readme_Hetzner.txt.

## Links utili

- https://hypermedia.systems/more-htmx-patterns/
- https://wiki.selfhtml.org/wiki/Formulare/Benutzereingaben_zug%C3%A4nglich_gestalten
- https://github.com/johan-st/go-image-server/blob/main/api.go#L111
- https://developer.mozilla.org/en-US/docs/Glossary/Semantics#semantics_in_html
- https://akismet.com/plan/personal/

## Validazione utente

Ho bisogno di un  token e un hash per validare l'utente che amministra il sito.
Per generare il JWT token ho bisogno di una chiave privata.
Per validare il JWT token ho bisogno della chiave pubblica.
Per calcolare l'hash della password salvata in cred.json ho bisogno di una chiave privata key.pem.

La chiave privata criptata messa nel file key.pem è generata da questo programma (funzione savePrivateKeyInFile)
La password che genera l'hash viene chiesta all'atto della creazione dell'account admin.
Nel file cred.json c'è l'hash della password di admin.
Il salt per l'hash della password è generato casualmente nel momento in cui viene creata la
chiave privata key.pem.

Comando:

    go run .\main.go -initaccount

## Key.pem
È la chiave privata che viene usata per generare il token JWT e l'hash dell'utente.

## Chiave pubblica
Per validare il token jwt occorre la chiave pubblica, che ricavo con WSL dalla chiave pem
generata in con questo programma attraverso:

    openssl rsa -in key.pem -pubout -out pubkey.pem

Il token Jwt vale un'ora e non uso il refresh. Viene memorizzato nel browser session store. È attaccabile via XSS (https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps), ma senza cors, nessun utente extra, l'unico punto dovrebbe essere il commento, che, utilizzando html, potrebbe eseguire del codice esterno che va a modificare la app (la static get). Quando inserisco un nuovo commento, ho un check del contenuto con Bluemonday in StrictPolicy, che non ammette html.
Per riuscire a cambiare la app admin, il commento che arriva dal db deve essere un html. Qui il rendering del post non deve assolutamente generare html, ma semplicemente mostrare una stringa html-escaped. 

## Dominio
Ho riservato il nome: igorrun.invido.it

## Database
Ho separato due database. Il primo blog-comments.db è solo per i commenti e serve solo sul server remoto.
In generale non lo devo aggiornare in quanto la gestione dei commenti avviene con l'admin. 
Il secondo database è blog-corsa.db e mi serve per creare i link. Questo db serve per il programma src.exe
per creare i post e le pages.
La ricerca non l'ho ancora implementata, ma dovrebbe rimanere nel db blog-corsa.db.

## Ricreare il sito da zero
Se per caso devo ricreare il sito (links, pages e posts)

    .\src.exe -config ..\..\config.toml -rebuildall

## Creare un nuovo Post
Al momento il processo funziona con Visual Code (profilo Edit Post).
Il database sarebbe meglio scaricarlo da current su invido.it.
Per il nuovo post:

    cd .\content\src\
    .\src.exe -config ..\..\config.toml  -newpost "Prossima Backyard: Frankenmarkt" -date "2025-06-22" -watch

Ora edito il nuovo file mdhtml e vedo subito il risultato (nell'esempio di sopra su http://localhost:5572/posts/2025/04/17/25-04-17-NuovoSito/).

Ora devo attualizzare i links:

    .\src.exe -config ..\..\config.toml -scancontent

Creare i posts col feed:

    .\src.exe -config ..\..\config.toml -buildposts
Creare la main page (force non mi piace perchè cambia autore e statistiche, vedi todo):

    .\src.exe -config ..\..\config.toml -buildpages -force


Siccome ho separato i due db con i commenti, il sync dei commenti non è necessario. 
Il db blog-corsa.db rimane dove viene creato il post.
In futuro, con la funzione "cerca", il sync del db con i dati della ricerca sarà necessario.