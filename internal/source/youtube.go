// Package source fournit les fournisseurs de contenu (recherche → résultats).
//
// Chaque source légale (YouTube, plus tard Twitch, IPTV, bibliothèque locale…)
// implémente la même idée : à partir d'une requête, renvoyer une liste de
// Result que l'utilisateur choisit, puis qu'on lit avec le lecteur.
package source

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

// Result est un élément lisible (une vidéo, un flux…).
type Result struct {
	Title    string
	URL      string
	Uploader string
	Duration string // format lisible, ex "3:45"
}

// SearchYouTube interroge YouTube via yt-dlp et renvoie jusqu'à n résultats.
// yt-dlp doit être installé et dans le PATH.
func SearchYouTube(query string, n int) ([]Result, error) {
	if n <= 0 {
		n = 10
	}
	// --flat-playlist : ne résout pas chaque vidéo (rapide) ; --dump-json : 1 objet/ligne.
	cmd := exec.Command("yt-dlp",
		"--flat-playlist",
		"--dump-json",
		"--no-warnings",
		fmt.Sprintf("ytsearch%d:%s", n, query),
	)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("yt-dlp : %v (%s)", err, errBuf.String())
	}

	var results []Result
	scanner := bufio.NewScanner(&out)
	scanner.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for scanner.Scan() {
		ligne := scanner.Bytes()
		if len(ligne) == 0 {
			continue
		}
		var e ytEntry
		if err := json.Unmarshal(ligne, &e); err != nil {
			continue // on ignore les lignes non-JSON éventuelles
		}
		results = append(results, e.toResult())
	}
	return results, scanner.Err()
}

// ytEntry mappe les champs utiles d'une entrée yt-dlp.
type ytEntry struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	URL      string  `json:"url"`
	Uploader string  `json:"uploader"`
	Channel  string  `json:"channel"`
	Duration float64 `json:"duration"`
}

func (e ytEntry) toResult() Result {
	url := e.URL
	if e.ID != "" {
		url = "https://www.youtube.com/watch?v=" + e.ID
	}
	auteur := e.Uploader
	if auteur == "" {
		auteur = e.Channel
	}
	return Result{
		Title:    e.Title,
		URL:      url,
		Uploader: auteur,
		Duration: formatDuree(int(e.Duration)),
	}
}

func formatDuree(secondes int) string {
	if secondes <= 0 {
		return "?"
	}
	h := secondes / 3600
	m := (secondes % 3600) / 60
	s := secondes % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return strconv.Itoa(m) + ":" + fmt.Sprintf("%02d", s)
}
