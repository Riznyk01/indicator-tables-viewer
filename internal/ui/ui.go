package ui

import "github.com/go-vgo/robotgo"

func SetResolution() (w, h float32) {
	width, height := robotgo.GetScreenSize()
	if width > 1920 && height > 1080 {
		w = 0.5 * float32(width)
		h = 0.4 * float32(height)
		return w, h
	}
	w = 0.8 * float32(width)
	h = 0.7 * float32(height)
	return w, h
}
