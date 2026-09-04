package helper

import (
	"errors"
	"strings"

	"github.com/Nixon-Alexander/resolve_now.git/model"
)

func ParseSections(text string) ([]model.DocumentSection, error) {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSpace(text)

	if text == "" {
		return nil, errors.New("document is empty")
	}

	lines := strings.Split(text, "\n")

	var sections []model.DocumentSection

	var currentTitle string
	var currentContent []string
	var insideSection bool

	flushSection := func() {
		if !insideSection {
			return
		}

		content := strings.TrimSpace(
			strings.Join(currentContent, "\n"),
		)

		if currentTitle != "" && content != "" {
			sections = append(
				sections,
				model.DocumentSection{
					Title:   currentTitle,
					Content: content,
				},
			)
		}

		currentTitle = ""
		currentContent = nil
		insideSection = false
	}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		switch line {

		// ========================================
		// TITLE
		// ========================================
		case "[TITLE]":
			// TITLE tidak dimasukkan ke section.
			// Untuk sekarang kita skip.
			continue

		// ========================================
		// INTRODUCTION
		// ========================================
		case "[INTRODUCTION]":
			// Introduction kita anggap sebagai
			// section juga.
			flushSection()

			currentTitle = "Introduction"
			insideSection = true

		// ========================================
		// SECTION
		// ========================================
		case "[SECTION]":
			// Simpan section sebelumnya
			flushSection()

			// Baris berikutnya harus menjadi title
			if i+1 >= len(lines) {
				return nil, errors.New(
					"section title is missing",
				)
			}

			title := strings.TrimSpace(lines[i+1])

			if title == "" {
				return nil, errors.New(
					"section title cannot be empty",
				)
			}

			currentTitle = title
			currentContent = nil
			insideSection = true

			// Skip baris title
			i++

		// ========================================
		// CONTENT
		// ========================================
		default:
			if insideSection {
				currentContent = append(
					currentContent,
					line,
				)
			}
		}
	}

	// Simpan section terakhir
	flushSection()

	if len(sections) == 0 {
		return nil, errors.New(
			"document has no sections",
		)
	}

	return sections, nil
}
