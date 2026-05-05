package applications

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alehbelskidev/job_appl_track/internal/repo"
	"github.com/alehbelskidev/job_appl_track/internal/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type service struct {
	openai *openai.Client
	repo   repo.Querier
}

func newService(c *shared.Config, repo repo.Querier) *service {
	client := openai.NewClient(option.WithAPIKey(string(c.OpenAISecret)))

	return &service{
		openai: &client,
		repo:   repo,
	}
}

func (s *service) parseHtmlPage(url string) (*string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("script, style, nav, footer, header, noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	cleanText := doc.Find("body").Text()

	result := strings.Join(strings.Fields(cleanText), " ")

	return &result, nil
}

func (s *service) createApplicationFromPrompt(
	ownerID uuid.UUID,
	ctx context.Context,
	prompt string,
	notes string,
	url string,
) (*repo.JobApplication, error) {
	resp, err := s.ask(ctx, prompt)
	if err != nil {
		return nil, err
	}

	createAppDto := repo.CreateJobApplicationParams{
		DateApplied: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Status:      int32(applied),
		OwnerID:     ownerID,
		Company:     resp.Company,
		Role:        resp.Role,
		Description: pgtype.Text{String: resp.Description, Valid: resp.Description != ""},
		Url:         pgtype.Text{String: url, Valid: url != ""},
		Notes:       pgtype.Text{String: notes, Valid: notes != ""},
	}

	app, err := s.repo.CreateJobApplication(ctx, createAppDto)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

type openAIApplicationResponse struct {
	Company     string `json:"company"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

func (s *service) ask(ctx context.Context, prompt string) (*openAIApplicationResponse, error) {
	resp, err := s.openai.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model: openai.ChatModelGPT4o,
	})

	if err != nil {
		return nil, err
	}

	respStr := resp.OutputText()

	var parsedResp openAIApplicationResponse
	json.Unmarshal([]byte(respStr), &parsedResp)

	return &parsedResp, nil
}
