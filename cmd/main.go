package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/arthureichelberger/autobackstage/pkg/env"
	"github.com/arthureichelberger/autobackstage/pkg/github"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	repo := env.Get("GITHUB_REPOSITORY", "")
	backstageBrach := env.Get("BACKSTAGE_BRANCH", "backstage")
	gitSha := env.Get("GITHUB_SHA", "")

	ghClient := github.NewDefaultClient("https://api.github.com", env.Get("GITHUB_TOKEN", ""), repo)

	if err := run(ctx, ghClient, backstageBrach, gitSha); err != nil {
		log.Error().Err(err).Msg("could not execute action")
		os.Exit(1)
		return
	}

	os.Exit(0)
}

func run(ctx context.Context, ghClient github.Client, backstageBranch, gitSha string) error {
	if err := ghClient.DeleteBranch(ctx, backstageBranch); err != nil {
		log.Error().Err(err).Msg("could not delete branch")
		return fmt.Errorf("could not delete branch")
	}

	if err := ghClient.CreateBranch(ctx, backstageBranch, gitSha); err != nil {
		log.Error().Err(err).Msg("could not create branch")
		return fmt.Errorf("could not create branch")
	}

	return nil
}
