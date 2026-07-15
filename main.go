// onde — lecteur de streaming en terminal, pour sources légales.
//
// v0.1 : recherche YouTube (via yt-dlp) → menu → lecture (via mpv).
// Lit aussi directement une URL (flux HLS, site géré par yt-dlp) ou un fichier
// local.
//
// Usage :
//
//	onde <recherche...>              recherche YouTube puis choix
//	onde "https://.../flux.m3u8"     lit un flux/URL directement
//	onde "C:\films\video.mkv"        lit un fichier local
//
// Prérequis : mpv et yt-dlp dans le PATH.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Bastien-Gaffet/onde/internal/player"
	"github.com/Bastien-Gaffet/onde/internal/source"
)

const version = "0.1.0"

// Couleurs ANSI (les terminaux Windows 10/11 les gèrent).
const (
	reset = "\033[0m"
	gras  = "\033[1m"
	cyan  = "\033[36m"
	vert  = "\033[32m"
	jaune = "\033[33m"
	dim   = "\033[2m"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		usage()
		return
	}
	if args[0] == "-v" || args[0] == "--version" {
		fmt.Println("onde", version)
		return
	}

	// URL ou fichier → lecture directe.
	if estCible(args[0]) {
		lire(args[0])
		return
	}

	// Sinon : recherche.
	requete := strings.Join(args, " ")
	rechercherEtLire(requete)
}

// estCible reconnaît une URL (http, https, flux) ou un fichier local existant.
func estCible(arg string) bool {
	if strings.Contains(arg, "://") {
		return true
	}
	if info, err := os.Stat(arg); err == nil && !info.IsDir() {
		return true
	}
	return false
}

func rechercherEtLire(requete string) {
	fmt.Printf("%s🔎 Recherche :%s %s\n", cyan, reset, requete)
	resultats, err := source.SearchYouTube(requete, 10)
	if err != nil {
		fmt.Printf("%sErreur de recherche :%s %v\n", jaune, reset, err)
		os.Exit(1)
	}
	if len(resultats) == 0 {
		fmt.Println("Aucun résultat.")
		return
	}

	for i, r := range resultats {
		fmt.Printf("  %s%2d%s  %s%s%s  %s%s · %s%s\n",
			gras, i+1, reset,
			r.Title, "", "",
			dim, r.Uploader, r.Duration, reset)
	}

	choix := demanderChoix(len(resultats))
	if choix < 0 {
		fmt.Println("Annulé.")
		return
	}
	lire(resultats[choix].URL)
}

// demanderChoix lit un numéro valide (1..n) ; renvoie l'index 0-based, ou -1 si annulé.
func demanderChoix(n int) int {
	lecteur := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%sNuméro à lire (Entrée = annuler) :%s ", cyan, reset)
		ligne, err := lecteur.ReadString('\n')
		if err != nil {
			return -1
		}
		ligne = strings.TrimSpace(ligne)
		if ligne == "" {
			return -1
		}
		num, err := strconv.Atoi(ligne)
		if err != nil || num < 1 || num > n {
			fmt.Printf("%s  Entrez un numéro entre 1 et %d.%s\n", jaune, n, reset)
			continue
		}
		return num - 1
	}
}

func lire(cible string) {
	fmt.Printf("%s▶ Lecture :%s %s\n", vert, reset, cible)
	if err := player.Play(cible); err != nil {
		fmt.Printf("%s%v%s\n", jaune, err, reset)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`%sonde%s %s — lecteur de streaming en terminal (sources légales)

%sUsage :%s
  onde <recherche...>            recherche YouTube puis choix au clavier
  onde "https://.../flux.m3u8"   lit un flux / une URL directement
  onde "C:\films\video.mkv"      lit un fichier local

%sExemples :%s
  onde lofi hip hop radio
  onde "https://www.twitch.tv/xxxx"
  onde -v

%sPrérequis :%s mpv et yt-dlp dans le PATH.
`, gras, reset, version, gras, reset, gras, reset, gras, reset)
}
