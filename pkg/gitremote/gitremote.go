// Package gitremote parses git remote URLs into their owner / repo / host
// components. MVP supports GitHub hosts only; other hosts return ok=false.
package gitremote

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Remote is a parsed git remote URL.
type Remote struct {
	Owner string
	Repo  string
	Host  string // e.g. "github.com"
}

var (
	// https://github.com/owner/repo(.git)?  or  http://...
	httpsRE = regexp.MustCompile(`^https?://([^/]+)/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	// ssh://git@github.com/owner/repo(.git)?
	sshURLRE = regexp.MustCompile(`^ssh://[^@]+@([^/]+)/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	// git@github.com:owner/repo(.git)?
	sshSCPRE = regexp.MustCompile(`^[^@]+@([^:]+):([^/]+)/([^/]+?)(?:\.git)?$`)
)

// Parse extracts the owner/repo/host from a git remote URL. It returns
// (Remote, true) only for GitHub hosts in MVP; any other host (GitLab,
// Bitbucket, self-hosted) yields (_, false) so callers can gracefully
// skip rather than emit broken links.
func Parse(url string) (Remote, bool) {
	url = strings.TrimSpace(url)
	for _, re := range []*regexp.Regexp{httpsRE, sshURLRE, sshSCPRE} {
		if m := re.FindStringSubmatch(url); m != nil {
			r := Remote{Host: strings.ToLower(m[1]), Owner: m[2], Repo: m[3]}
			if r.Host != "github.com" {
				return Remote{}, false
			}
			return r, true
		}
	}
	return Remote{}, false
}

// OriginURL returns the URL of the "origin" remote for the git repository
// at dir. Returns an error if dir is not a git repo or has no origin remote.
func OriginURL(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git remote origin: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
