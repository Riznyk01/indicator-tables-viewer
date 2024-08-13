package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type viewerTheme struct {
	fyne.Theme
	fontSize float32
}

func newTermTheme(size float32) *viewerTheme {
	return &viewerTheme{
		Theme: fyne.CurrentApp().Settings().Theme(), fontSize: size,
	}
}

func (t *viewerTheme) Size(n fyne.ThemeSizeName) float32 {
	if n == theme.SizeNameText {
		return t.fontSize
	}

	return t.Theme.Size(n)
}
