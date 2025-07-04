// Patch pour filtrer la langue avant TTS
package patch

// FilterLanguageForTTS vérifie si la langue détectée est supportée (fr ou uk)
// et retourne le texte à synthétiser (message d'erreur si non supportée)
func FilterLanguageForTTS(lang, text string) string {
	supported := map[string]bool{"fr": true, "uk": true}
	if !supported[lang] {
		return "La langue détectée n’est pas prise en charge. Veuillez parler en français ou en ukrainien."
	}
	return text
}
