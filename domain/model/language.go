package model

type Language string

const (
	LanguageEnglish Language = "EN"
	LanguageFarsi   Language = "FA"
)

var AvailableLanguages = []Language{
	LanguageEnglish,
	LanguageFarsi,
}
