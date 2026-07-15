// Package player lance la lecture d'une URL ou d'un fichier via mpv.
//
// mpv est un lecteur vidéo léger, pilotable au clavier, qui lit aussi bien les
// fichiers locaux que les flux réseau (HLS/m3u8) et, via yt-dlp, ~1800 sites.
package player

import (
	"fmt"
	"os"
	"os/exec"
)

// Disponible indique si mpv est présent dans le PATH.
func Disponible() bool {
	_, err := exec.LookPath("mpv")
	return err == nil
}

// Play ouvre `cible` (URL ou chemin de fichier) dans mpv et attend la fin.
// Le terminal reste rendu à l'utilisateur (mpv gère ses propres touches).
func Play(cible string) error {
	if !Disponible() {
		return fmt.Errorf("mpv introuvable dans le PATH — installez-le : scoop install mpv")
	}
	cmd := exec.Command("mpv",
		"--force-window=immediate", // fenêtre vidéo dès le lancement
		cible,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
