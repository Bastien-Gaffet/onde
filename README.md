# onde

**Lecteur de streaming en terminal, pour sources légales.** Recherche et lis des
vidéos/flux sans quitter le terminal — l'esprit d'`ani-cli`, appliqué à des sources
légitimes (YouTube, Twitch, IPTV/HLS, ta bibliothèque locale).

> `onde` ne scrape aucun site de piratage. Il s'appuie sur **yt-dlp** (recherche et
> résolution de ~1800 sites) et **mpv** (lecture), et lit ce que tu as le droit de lire.

---

## Prérequis

- [**mpv**](https://mpv.io) — le lecteur : `scoop install mpv` (Windows) · `brew install mpv` · `apt install mpv`
- [**yt-dlp**](https://github.com/yt-dlp/yt-dlp) — recherche & résolution : `pip install yt-dlp` ou `scoop install yt-dlp`

## Installation

```bash
git clone https://github.com/Bastien-Gaffet/onde
cd onde
go build -o onde.exe .        # Windows (ou : go install .)
```

## Utilisation

```bash
# Rechercher sur YouTube puis choisir
onde lofi hip hop radio

# Lire un flux / une URL directement (Twitch, HLS, site yt-dlp…)
onde "https://www.twitch.tv/xxxx"
onde "https://exemple.fr/flux.m3u8"

# Lire un fichier local
onde "C:\films\video.mkv"
```

Dans le lecteur (mpv) : `espace` pause · `f` plein écran · `←/→` avancer/reculer ·
`9/0` volume · `q` quitter.

---

## Feuille de route

- [x] **v0.1** — recherche YouTube (yt-dlp) → menu → lecture (mpv) ; URL/fichier direct
- [ ] **v0.2** — vraie TUI (Bubble Tea) : navigation fléchée, aperçu, recherche live
- [ ] **v0.3** — sources multiples : Twitch, playlists IPTV/`.m3u`, **bibliothèque locale**
- [ ] **v0.4** — historique + **reprise de lecture** (« continuer où j'en étais »)
- [ ] **v0.5** — favoris, qualité par défaut, sous-titres

---

## Architecture

```
main.go                 entrée CLI (recherche vs lecture directe)
internal/source/        fournisseurs de contenu (youtube.go ; puis twitch, iptv, local…)
internal/player/        lecture (mpv.go)
```

Chaque source renvoie une liste de `Result` (titre, URL, auteur, durée) ; le lecteur
lit une URL/fichier. Ajouter une source = un nouveau fichier dans `internal/source/`.

## Licence

MIT — voir [LICENSE](LICENSE).
