# Corsa Blog
Un service per gestire il mio blog della corsa con i commenti.

## Commenti
Quando viene postato un commento sul blog, esso viene dapprima esaminato con il
service [askimet](https://akismet.com/plan/personal/) per vedere se è uno spam.
Se non lo è viene mandata una notifica via mail e/o telegram per approvare il commento.

## Dashboard Admin
La parte che riguarda l'amministrazione del blog è gestita con vue in modalità single page.
Per contro, la parte testuale dei vari post è generata staticamente in html.
Al momento la Dashboard mi serve solo per gestire i commenti. POtrebbe, però, essere usata per
creare anche dei nuovi contenuti. Questo vorrebbe dire gestire la generazione statica di html con Zine,
anche se la fase di preview mi sembra molto complessa da gestire all'interno del browser.
La URL di riferimento è: http://localhost:5572/blog-admin/

## Blog Corsa Static-Html
La parte testuale che fornisce il blog sulla corsa è situata in static/blog-corsa.
Questa è una subdirectory totalmente autonoma generata automaticamente. Al momento uso Zine.
L'output è per esempio in D:\scratch\zig\zine\blog-corsa\zig-out   
La URL di riferimento è  http://localhost:5572. Nota che è una URL root. Non può essere in 
un path che non sia / in quanto il codice html generato da Zine usa tutti i path relativi a root.

Uso htmx per avere il fetch delle parti dinamiche, come per esempio i commenti di un post, nella
parte statica html. Nota che in questi casi i servizi come _commento_, che si trova su [gitlab.com/commento](https://gitlab.com/commento/commento), usa un approccio differente. Vale a dire i commenti vengono
visualizzati usanto una request via javascript al server di Commento, il quale risponde con un json.
Il file javascript sul client crea poi on-demand il codice html che viene aggiunto all'elemento div indicato.
Commento.js è un file di 60k non min.  

Quando uso htmx, invece, il server fornisce la parte html già creata senza passare da json. 

### TODO
- config_custom.toml va cryptato nel valore di diversi campi
- La pagina admin deve essere protetta da un token di sign-in 
- Nella pagina di risposta newcomment dal form, va messo un testo e un link per tornare indietro
- Per testare la sezione _Discussione_ uso del codice html direttamente in index.html. Probabilmente
dovrebbe essere generato completamente sul server con un get.
- Gestione dei commenti nel file data.json in quanto non voglio usare un database
- Nella pagina admin, manca la gestione Edit/delete/approve/decline dei commenti

### Stop del service
Per stoppare il sevice si usa:

    sudo systemctl stop corsa-blog

## Deployment su ubuntu direttamente

    cd ~/build/corsa-blog
    git pull --all
    ./publish-service.sh
Oppure uso Visual Code in remoto dove uso il synch di git. Qui nel terminal mi basta usare:

    ./publish-service.sh

## Service Config
Per prima cosa va creato il file corsa-blog.service.
Il contenuto l'ho messo sotto in una sezione apposita.

    sudo nano /lib/systemd/system/corsa-blog.service
Poi si fa l'enable:

    sudo systemctl enable corsa-blog.service
E infine lo start:

    sudo systemctl start corsa-blog
Logs sono disponibili con:

    sudo journalctl -f -u corsa-blog

## corsa-blog.service
Qui segue il contenuto del file corsa-blog.service
Nota il Type=idle che è meglio di simple. Così 
viene fatto partire dopo che la wlan ha ottenuto l'IP intranet
e così si ha l'accesso.

```
[Install]
WantedBy=multi-user.target

[Unit]
Description=corsa-blog service
ConditionPathExists=/home/igor/app/go/corsa-blog/current/corsa-blog.bin
After=network.target

[Service]
Type=idle
User=igor
Group=igor
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/igor/app/go/corsa-blog/current/
ExecStart=/home/igor/app/go/corsa-blog/current/corsa-blog.bin

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/corsa-blog
ExecStartPre=/bin/chown igor:igor /var/log/corsa-blog
ExecStartPre=/bin/chmod 755 /var/log/corsa-blog
StandardOutput=syslog
StandardError=syslog
```

## Data.json
Nel file data.json ho messo la lista dei commenti. Non uso un database in quanto i
commenti da gestire sono molto pochi e possono essere benissimo scritti in un unico file.


### config_custom.toml
È il file che mi esegue un ovveride del file config.toml. 
Mi serve in quanto config.toml si trova su gitHub, mentre config_custom.toml è
solo locale fuori da git. Si trova in:

    /home/igor/app/go/corsa-blog/current/

## Visual Code
Per lo sviluppo iniziale ho usato windows, poi, per l'update del service,
ho usato Visual Code Remote nella directory ~/build/corsa-blog.

## nginx proxy
todo

## Links utili

- https://hypermedia.systems/more-htmx-patterns/
- https://wiki.selfhtml.org/wiki/Formulare/Benutzereingaben_zug%C3%A4nglich_gestalten
- Image server: https://github.com/johan-st/go-image-server/blob/main/api.go#L111


