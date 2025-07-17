package gitlab

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hbk619/gh-peruse/internal/git"
	"github.com/hbk619/gh-peruse/internal/requests"
)

type Repo interface {
	Base(remoteName string) (*git.Repo, error)
}

type GLRepo struct {
	runner requests.CommandLine
}

func NewGLRepo(runner requests.CommandLine) *GLRepo {
	return &GLRepo{
		runner: runner,
	}
}
func (glRepo *GLRepo) Base(remoteName string) (*git.Repo, error) {
	if remoteName == "" {
		remoteName = "origin"
	}

	output, err := glRepo.runner.Run("git", []string{"remote", "get-url", remoteName})
	if err != nil {
		return nil, fmt.Errorf("failed to get remote URL for '%s': %w", remoteName, err)
	}

	remoteURL := strings.TrimSpace(string(output))
	if remoteURL == "" {
		return nil, fmt.Errorf("no remote URL found for '%s'", remoteName)
	}

	return ParseGitRepoURL(remoteURL)
}

func ParseGitRepoURL(repoURL string) (*git.Repo, error) {
	if repoURL == "" {
		return nil, fmt.Errorf("repository URL cannot be empty")
	}

	repoURL = strings.TrimSuffix(repoURL, ".git")

	if strings.HasPrefix(repoURL, "git@") {
		return parseSSHURL(repoURL)
	}

	if strings.HasPrefix(repoURL, "https://") || strings.HasPrefix(repoURL, "http://") {
		return parseHTTPSURL(repoURL)
	}

	return nil, fmt.Errorf("unsupported URL format: %s", repoURL)
}

func parseSSHURL(repoURL string) (*git.Repo, error) {
	repoURL = strings.TrimPrefix(repoURL, "git@")

	parts := strings.SplitN(repoURL, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SSH URL format")
	}

	host := parts[0]
	path := parts[1]

	pathParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid repository path: %s", path)
	}

	repoName := pathParts[len(pathParts)-1]
	owner := strings.Join(pathParts[:len(pathParts)-1], "/")

	return &git.Repo{
		Host:  host,
		Owner: owner,
		Name:  repoName,
	}, nil
}
func parseHTTPSURL(repoURL string) (*git.Repo, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	host := parsedURL.Host
	path := strings.Trim(parsedURL.Path, "/")

	if path == "" {
		return nil, fmt.Errorf("empty repository path")
	}

	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid repository path: %s", path)
	}

	repoName := pathParts[len(pathParts)-1]
	owner := strings.Join(pathParts[:len(pathParts)-1], "/")

	return &git.Repo{
		Host:  host,
		Owner: owner,
		Name:  repoName,
	}, nil
}
func GetRepoInfoFromGitRemote(remoteURL string) (*git.Repo, error) {
	remoteURL = strings.TrimSpace(remoteURL)
	remoteURL = strings.Trim(remoteURL, "\"'")

	return ParseGitRepoURL(remoteURL)
}
