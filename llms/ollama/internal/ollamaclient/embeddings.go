package ollamaclient

import (
	"context"
	"net/http"
)

const (
	defaultEmbeddingModel = "text-embedding-ada-002"
)

func (c *Client) CreateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	resp := &EmbeddingResponse{}
	if err := c.do(ctx, http.MethodPost, "/api/embeddings", req, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
