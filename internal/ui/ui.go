package ui

import (
	"github.com/go-vgo/robotgo"
	"indicator-tables-viewer/internal/config"
)

func SetResolution(cfg *config.Config) (w, h float32) {
	width, height := robotgo.GetScreenSize()
	w = float32(width) * cfg.WidthMultiplier
	h = float32(height) * cfg.HeightMultiplier
	return w, h
}
