package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/arthureichelberger/autobackstage/pkg/github/model"
)

type defaultClient struct {
	BaseURL string
	Token   string
	Repo    string
	client  http.Client
}

func NewDefaultClient(baseURL, token, repo string) defaultClient {
	return defaultClient{
		BaseURL: baseURL,
		Token:   token,
		Repo:    repo,
		client: http.Client{
			Timeout: time.Second,
		},
	}
}

func (dc defaultClient) CreateBranch(ctx context.Context, branch, sha string) error {
	url := fmt.Sprintf("%s/repos/%s/git/refs", dc.BaseURL, dc.Repo)
	payload := map[string]string{"ref": fmt.Sprintf("refs/heads/%s", branch), "sha": sha}
	payloadJSON, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("could not create request to create branch")
		return fmt.Errorf("could not create branch")
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", dc.Token))

	res, err := dc.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("could not execute request to create branch")
		return fmt.Errorf("could not execute request")
	}

	if res.StatusCode == http.StatusUnprocessableEntity {
		resError, err := dc.getResponseError(res.Body)
		if err != nil {
			log.Error().Err(err).Str("url", url).Msg("could not get response error")
			return err
		}

		log.Error().Str("message", resError.Message).Str("url", url).Msg("could not process entity while deleting branch")
	}

	if res.StatusCode != http.StatusCreated {
		log.Error().Str("url", url).Str("status", res.Status).Msg("could not create branch")
		return fmt.Errorf("could not create branch")
	}

	return nil
}

func (dc defaultClient) DeleteBranch(ctx context.Context, branch string) error {
	url := fmt.Sprintf("%s/repos/%s/git/refs/heads/%s", dc.BaseURL, dc.Repo, branch)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("could not create request to delete branch")
		return fmt.Errorf("could not create request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", dc.Token))

	res, err := dc.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("could not execute request to delete branch")
		return fmt.Errorf("could not execute request")
	}

	if res.StatusCode == http.StatusUnprocessableEntity {
		resError, err := dc.getResponseError(res.Body)
		if err != nil {
			log.Error().Err(err).Str("url", url).Msg("could not get response error")
			return err
		}

		if resError.Message == model.RefDoesNotExistError {
			return nil
		}

		log.Error().Str("message", resError.Message).Str("url", url).Msg("could not process entity while deleting branch")
	}

	if res.StatusCode != http.StatusNoContent {
		log.Error().Str("status", res.Status).Str("url", url).Msg("could not delete branch")
		return fmt.Errorf("could not delete branch")
	}

	return nil
}

func (dc defaultClient) getResponseError(body io.ReadCloser) (model.ResponseError, error) {
	var resError model.ResponseError
	if err := json.NewDecoder(body).Decode(&resError); err != nil {
		log.Error().Err(err).Msg("could not decode body")
		return model.ResponseError{}, fmt.Errorf("could not decode response error")
	}

	return resError, nil
}
