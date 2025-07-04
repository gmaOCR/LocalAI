// Patch pour filtrer la langue avant TTS (Go)
package patch

import (
	"github.com/abadojack/whatlanggo"
)

// ApplyTTSLanguageFilter détecte la langue du texte et applique le filtrage
func ApplyTTSLanguageFilter(input string, cfgLang string) string {
	info := whatlanggo.Detect(input)
	lang := ""
	if info.Lang == whatlanggo.Fra {
		lang = "fr"
	} else if info.Lang == whatlanggo.Ukr {
		lang = "uk"
	} else {
		lang = "autre"
	}
	// Filtrage et éventuel message d'erreur
	input = FilterLanguageForTTS(lang, input)
	input = FilterLanguageForTTS(cfgLang, input)
	return input
}
