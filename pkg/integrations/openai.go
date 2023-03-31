package integrations

import (
	"context"
	"fmt"

	tugit "github.com/b4nst/turbogit/pkg/git"
	git "github.com/libgit2/git2go/v33"
	"github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(r *git.Repository) (*OpenAIProvider, error) {
	c, err := r.Config()
	if err != nil {
		return nil, err
	}
	enabled, err := c.LookupBool("openai.enabled")
	if !enabled {
		return nil, err
	}
	token, err := c.LookupString("openai.token")
	if err != nil {
		return nil, err
	}
	oc := openai.NewClient(token)

	return &OpenAIProvider{client: oc}, nil
}

func (oai *OpenAIProvider) CommitMessages(diff *git.Diff) ([]string, error) {
	spatch, err := tugit.PatchFromDiff(diff)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Propose a conventional commit message for this git diff \n```\n%s\n```", spatch)

	resp, err := oai.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	cmos := make([]string, 0, len(resp.Choices))
	for _, c := range resp.Choices {
		cmos = append(cmos, c.Message.Content)
	}

	return cmos, nil
}
