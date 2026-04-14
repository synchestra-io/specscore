package lint

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/synchestra-io/specscore/pkg/gitremote"
	"github.com/synchestra-io/specscore/pkg/projectdef"
)

// hubViewLinkMarker is the unique prefix we write / look for. Any line
// starting with this under the H1 is treated as the Hub view-link
// blockquote (stale or fresh).
const hubViewLinkMarker = "> [View in Synchestra Hub]("

// hubViewLinkSuffix is the trailing copy describing what Hub offers.
const hubViewLinkSuffix = ") — graph, discussions, approvals"

// hubViewLinkChecker verifies that every feature README carries an opt-in
// "View in Synchestra Hub" blockquote directly under its H1. The rule is a
// no-op unless `hub.host` is set in specscore-spec-repo.yaml.
type hubViewLinkChecker struct{}

func newHubViewLinkChecker() checker { return &hubViewLinkChecker{} }

func (c *hubViewLinkChecker) name() string     { return "hub-view-link" }
func (c *hubViewLinkChecker) severity() string { return "warning" }

// BuildHubViewURL returns the canonical Hub view URL for a feature README
// at relPath (relative to the project root, e.g. "spec/features/bots").
// host is the Hub base URL without trailing slash.
func BuildHubViewURL(host string, r gitremote.Remote, relPath string) string {
	id := fmt.Sprintf("%s@%s@%s", r.Repo, r.Owner, r.Host)
	relPath = filepath.ToSlash(relPath)
	return fmt.Sprintf("%s/project/features?id=%s&path=%s",
		host, id, url.QueryEscape(relPath))
}

// hubViewLinkContext resolves the project-level inputs (Hub host + git
// remote) once per run. If Hub is not opted in, enabled == false and no
// walking happens. Any opt-in-but-broken configuration (no git remote,
// non-GitHub remote) is surfaced as a single warning at the config file.
type hubViewLinkContext struct {
	enabled  bool
	hubHost  string
	remote   gitremote.Remote
	skipWith Violation // populated if enabled but remote unusable
	skipSet  bool
}

func resolveHubViewLinkContext(specRoot string) hubViewLinkContext {
	projectRoot := filepath.Dir(specRoot)
	cfg, err := projectdef.ReadSpecConfig(projectRoot)
	if err != nil {
		return hubViewLinkContext{}
	}
	host := cfg.HubHost()
	if host == "" {
		return hubViewLinkContext{}
	}

	ctx := hubViewLinkContext{enabled: true, hubHost: host}
	originURL, err := gitremote.OriginURL(projectRoot)
	if err != nil {
		ctx.skipSet = true
		ctx.skipWith = Violation{
			File:     projectdef.SpecConfigFile,
			Line:     0,
			Severity: "warning",
			Rule:     "hub-view-link",
			Message:  "hub.host is set but origin remote could not be read; rule skipped",
		}
		return ctx
	}
	remote, ok := gitremote.Parse(originURL)
	if !ok {
		ctx.skipSet = true
		ctx.skipWith = Violation{
			File:     projectdef.SpecConfigFile,
			Line:     0,
			Severity: "warning",
			Rule:     "hub-view-link",
			Message:  "hub.host is set but origin remote is not a supported GitHub URL; rule skipped",
		}
		return ctx
	}
	ctx.remote = remote
	return ctx
}

func (c *hubViewLinkChecker) check(specRoot string) ([]Violation, error) {
	ctx := resolveHubViewLinkContext(specRoot)
	if !ctx.enabled {
		return nil, nil
	}
	if ctx.skipSet {
		return []Violation{ctx.skipWith}, nil
	}

	projectRoot := filepath.Dir(specRoot)
	var violations []Violation
	err := walkFeatureReadmes(specRoot, func(readmePath string, content []byte) {
		relFromRoot, _ := filepath.Rel(projectRoot, readmePath)
		relFromRoot = filepath.ToSlash(relFromRoot)
		expectedURL := BuildHubViewURL(ctx.hubHost, ctx.remote, filepath.Dir(relFromRoot))
		expectedLine := hubViewLinkMarker + expectedURL + hubViewLinkSuffix

		status := classifyHubViewLink(string(content), expectedLine)
		if status == hubViewLinkOK {
			return
		}
		relFromSpec, _ := filepath.Rel(specRoot, readmePath)
		msg := "missing 'View in Synchestra Hub' blockquote under the H1"
		if status == hubViewLinkStale {
			msg = "'View in Synchestra Hub' blockquote is out of date (wrong URL or copy)"
		}
		violations = append(violations, Violation{
			File:     relFromSpec,
			Line:     0,
			Severity: "warning",
			Rule:     "hub-view-link",
			Message:  msg,
		})
	})
	if err != nil {
		return nil, err
	}
	return violations, nil
}

func (c *hubViewLinkChecker) fix(specRoot string) error {
	ctx := resolveHubViewLinkContext(specRoot)
	if !ctx.enabled || ctx.skipSet {
		return nil
	}
	projectRoot := filepath.Dir(specRoot)
	return walkFeatureReadmes(specRoot, func(readmePath string, content []byte) {
		relFromRoot, _ := filepath.Rel(projectRoot, readmePath)
		relFromRoot = filepath.ToSlash(relFromRoot)
		expectedURL := BuildHubViewURL(ctx.hubHost, ctx.remote, filepath.Dir(relFromRoot))
		expectedLine := hubViewLinkMarker + expectedURL + hubViewLinkSuffix

		updated, changed := applyHubViewLink(string(content), expectedLine)
		if !changed {
			return
		}
		_ = os.WriteFile(readmePath, []byte(updated), 0o644)
	})
}

type hubViewLinkStatus int

const (
	hubViewLinkOK hubViewLinkStatus = iota
	hubViewLinkMissing
	hubViewLinkStale
)

func classifyHubViewLink(content, expectedLine string) hubViewLinkStatus {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, hubViewLinkMarker) {
			if line == expectedLine {
				return hubViewLinkOK
			}
			return hubViewLinkStale
		}
	}
	return hubViewLinkMissing
}

// applyHubViewLink inserts or replaces the Hub view-link blockquote. The
// blockquote is placed on its own line immediately after the first H1,
// separated by a single blank line above and below.
func applyHubViewLink(content, expectedLine string) (string, bool) {
	lines := strings.Split(content, "\n")

	// Replace an existing blockquote anywhere in the file (idempotent on stale).
	for i, line := range lines {
		if strings.HasPrefix(line, hubViewLinkMarker) {
			if line == expectedLine {
				return content, false
			}
			lines[i] = expectedLine
			return strings.Join(lines, "\n"), true
		}
	}

	// Insert after the first H1.
	h1 := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") {
			h1 = i
			break
		}
	}
	if h1 < 0 {
		return content, false
	}

	// Build the insertion: [blank][blockquote][blank]. Collapse if the
	// next line is already blank to avoid double blanks.
	insert := []string{"", expectedLine, ""}
	tail := lines[h1+1:]
	if len(tail) > 0 && tail[0] == "" {
		tail = tail[1:]
	}
	newLines := append([]string{}, lines[:h1+1]...)
	newLines = append(newLines, insert...)
	newLines = append(newLines, tail...)
	return strings.Join(newLines, "\n"), true
}

// walkFeatureReadmes invokes fn for every feature README under specRoot,
// skipping reserved _-prefixed subtrees. It mirrors the walk used by
// adherenceFooterChecker so both rules agree on scope.
func walkFeatureReadmes(specRoot string, fn func(readmePath string, content []byte)) error {
	featuresDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featuresDir)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(featuresDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if path != featuresDir && strings.HasPrefix(info.Name(), "_") {
			return filepath.SkipDir
		}
		readmePath := filepath.Join(path, "README.md")
		readmeInfo, statErr := os.Stat(readmePath)
		if statErr != nil || readmeInfo.IsDir() {
			return nil
		}
		content, readErr := os.ReadFile(readmePath)
		if readErr != nil {
			return nil
		}
		fn(readmePath, content)
		return nil
	})
}
