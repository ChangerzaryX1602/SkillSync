package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
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
func GenerateEmbeddingByOllama(ctx context.Context, FastHttpClient *fasthttp.Client, text string) ([]float32, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(viper.GetString("app.ollama_url") + "/api/embeddings")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	requestBody := map[string]interface{}{
		"model":  "bge-m3",
		"prompt": text,
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	req.SetBody(bodyBytes)

	if err := FastHttpClient.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("failed to generate embedding: %d", resp.StatusCode())
	}

	var response struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}

	return response.Embedding, nil
}
