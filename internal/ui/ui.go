package ui

import (
	"github.com/go-vgo/robotgo"
	"indicator-tables-viewer/internal/config"
)

func SetResolution(cfg *config.Config) (w, h float32) {
	width, height := robotgo.GetScreenSize()
	if width > 1920 && height > 1080 {
		w = cfg.W1Size * float32(width)
		h = cfg.H1Size * float32(height)
		return w, h
	}
	w = cfg.W2Size * float32(width)
	h = cfg.H2Size * float32(height)
	return w, h
}
