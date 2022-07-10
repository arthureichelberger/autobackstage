package github

import "context"

type Client interface {
	CreateBranch(ctx context.Context, branch, sha string) error
	DeleteBranch(ctx context.Context, branch string) error
}
