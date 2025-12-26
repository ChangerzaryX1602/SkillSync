package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
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
	var OutputDimensionality int32 = 768
	result, err := client.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		&genai.EmbedContentConfig{
			OutputDimensionality: &OutputDimensionality,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return result.Embeddings[0].Values, nil
}
func GenerateEmbeddingByOllama(text string) ([]float32, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  "bge-m3",
		"prompt": text,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(viper.GetString("app.ollama_url")+"/api/embeddings", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to generate embedding: %s", resp.Status)
	}

	var response struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Embedding, nil
}
