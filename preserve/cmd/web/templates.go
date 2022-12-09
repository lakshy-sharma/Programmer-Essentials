package main

import "preserve/pkg/models"

// This struct allows us to collect several dynamic data before passing to templates.
type templateData struct {
	Note *models.Note
}
