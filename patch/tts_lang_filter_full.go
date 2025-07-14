// Patch pour filtrer la langue avant TTS (Go)
package patch

import (
	"github.com/abadojack/whatlanggo"
)

// ApplyTTSLanguageFilter détecte la langue du texte et applique le filtrage
func ApplyTTSLanguageFilter(input string, cfgLang string) string {
	info := whatlanggo.Detect(input)
	if info.Lang == whatlanggo.Fra {
		// français détecté
	} else if info.Lang == whatlanggo.Ukr {
		// ukrainien détecté
	} else {
		// Si la langue détectée n'est pas supportée, message d'erreur
		return "La langue détectée n’est pas prise en charge. Veuillez parler en français ou en ukrainien."
	}
	// Si la langue est supportée, on laisse passer le texte original
	return input
}
