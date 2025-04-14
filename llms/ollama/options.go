package ollama

import (
	"net/http"
	"net/url"

	"github.com/mateors/llmg/llms/ollama/internal/ollamaclient"
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
