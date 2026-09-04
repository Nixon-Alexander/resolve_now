package helper

import (
	"context"
	"errors"

	"google.golang.org/genai"
)

func EmbedText(
	ctx context.Context,
	text string,
) ([]float32, error) {

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}

	contents := []*genai.Content{
		genai.NewContentFromText(
			text,
			genai.RoleUser,
		),
	}

	outputDim := int32(768)

	result, err := client.Models.EmbedContent(
		ctx,
		"gemini-embedding-2",
		contents,
		&genai.EmbedContentConfig{
			OutputDimensionality: &outputDim,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, errors.New("embedding result is empty")
	}

	return result.Embeddings[0].Values, nil
}
