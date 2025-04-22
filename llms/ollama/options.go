package ollama

import (
	"log"
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

type Option func(*options)

// WithModel Set the model to use.
func WithModel(model string) Option {
	return func(opts *options) {
		opts.model = model
	}
}

// WithFormat Sets the Ollama output format (currently Ollama only supports "json").
func WithFormat(format string) Option {
	return func(opts *options) {
		opts.format = format
	}
}

// WithSystem Set the system prompt. This is only valid if
// WithCustomTemplate is not set and the ollama model use
// .System in its model template OR if WithCustomTemplate
// is set using {{.System}}.
func WithSystemPrompt(p string) Option {
	return func(opts *options) {
		opts.system = p
	}
}

// WithCustomTemplate To override the templating done on Ollama model side.
func WithCustomTemplate(template string) Option {
	return func(opts *options) {
		opts.customModelTemplate = template
	}
}

// WithServerURL Set the URL of the ollama instance to use.
func WithServerURL(rawURL string) Option {
	return func(opts *options) {
		var err error
		opts.ollamaServerURL, err = url.Parse(rawURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// WithRunnerEmbeddingOnly Only return the embbeding.
func WithRunnerEmbeddingOnly(val bool) Option {
	return func(opts *options) {
		opts.ollamaOptions.EmbeddingOnly = val
	}
}
