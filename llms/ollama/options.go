package ollama

import (
	"llmg/llms/ollama/internal/ollamaclient"
	"net/http"
	"net/url"
)

type options struct {
	ollamaServerURL     *url.URL
	httpClient          *http.Client
	model               string
	ollamaOptions       ollamaclient.Options
	customModelTemplate string
	system              string
	format              string
	keepAlive           string
}
