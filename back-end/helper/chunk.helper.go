package helper

import (
	"strings"

	"github.com/Nixon-Alexander/resolve_now.git/model"
)

func ChunkDocument(
	sections []model.DocumentSection,
	maxChars int,
) []model.DocumentChunk {

	if maxChars <= 0 {
		maxChars = 1000
	}

	var chunks []model.DocumentChunk

	for _, section := range sections {
		// Kalau content masih muat dalam satu chunk
		if len(section.Content) <= maxChars {
			chunks = append(chunks, model.DocumentChunk{
				Index:        len(chunks),
				SectionTitle: section.Title,
				Content:      section.Content,
			})

			continue
		}
		// Kalau terlalu panjang
		splitted := splitText(section.Content, maxChars)

		for _, text := range splitted {
			chunks = append(chunks, model.DocumentChunk{
				Index:        len(chunks),
				SectionTitle: section.Title,
				Content:      text,
			})
		}
	}

	return chunks
}

func splitText(text string, maxChars int) []string {

	var result []string
	text = strings.TrimSpace(text)
	for len(text) > maxChars {
		// Cari newline terakhir sebelum maxChars
		cutIndex := strings.LastIndex(
			text[:maxChars],
			"\n",
		)

		// Kalau tidak ada newline,
		// cari spasi terakhir
		if cutIndex <= 0 {
			cutIndex = strings.LastIndex(
				text[:maxChars],
				" ",
			)
		}

		// Kalau tidak menemukan titik pemisah,
		// terpaksa potong berdasarkan karakter
		if cutIndex <= 0 {
			cutIndex = maxChars
		}

		chunk := strings.TrimSpace(text[:cutIndex])

		if chunk != "" {
			result = append(result, chunk)
		}

		text = strings.TrimSpace(text[cutIndex:])
	}

	if text != "" {
		result = append(result, text)
	}

	return result
}
