package helper

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func AskLLM(
	ctx context.Context,
	question string,
	contextText string,
) (string, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_API_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return "", err
	}

	systemInstruction := `
		Kamu adalah AI assistant untuk sistem complaint management.

		Kamu HARUS menjawab pertanyaan HANYA berdasarkan context
		dokumen yang diberikan.

		Aturan:
		- Gunakan informasi yang terdapat di context.
		- Jangan mengarang informasi.
		- Jangan menggunakan pengetahuan di luar context.
		- Jawab dalam bahasa Indonesia.
		- Jawab secara natural dan ringkas.
		- Jika informasi yang ditanyakan tidak ditemukan dalam context,
		katakan bahwa informasi tersebut tidak ditemukan.
	`

	prompt := fmt.Sprintf(`Context dokumen: %s Pertanyaan user: %s`, contextText, question)

	contents := []*genai.Content{
		genai.NewContentFromText(
			prompt,
			genai.RoleUser,
		),
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(
			systemInstruction,
			genai.RoleUser,
		),
	}

	response, err := client.Models.GenerateContent(ctx, "gemini-3.7-flash", contents, config)

	if err != nil {
		return "", err
	}

	return response.Text(), nil
}
