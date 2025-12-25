package utils

import (
	"context"
	"log"

	"google.golang.org/genai"
)

func GenerateEmbedding(text string, embeddingKey string) ([]float32, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: embeddingKey,
	})
	if err != nil {
		log.Fatal(err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText(text, genai.RoleUser),
	}
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return result.Embeddings[0].Values, nil
}
