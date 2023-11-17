package models

type Data struct {
	Path     string `path:"filepath"`
	Routines int
}

type Counter struct {
	Words, Vowels, Punctuations, Lines int
}
