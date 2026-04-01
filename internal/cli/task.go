package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/feature"
	"github.com/synchestra-io/specscore/pkg/task"
	"gopkg.in/yaml.v3"
)

func taskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task management — list, info, and create tasks",
	}
	cmd.AddCommand(
		taskListCommand(),
		taskInfoCommand(),
		taskNewCommand(),
	)
	return cmd
}

// resolveTasksDir resolves the tasks directory from a --project flag or CWD.
func resolveTasksDir(projectFlag string) (string, error) {
	var startDir string
	if projectFlag != "" {
		abs, err := filepath.Abs(projectFlag)
		if err != nil {
			return "", exitcode.InvalidArgsErrorf("resolving --project path: %v", err)
		}
		startDir = abs
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", exitcode.UnexpectedErrorf("cannot determine working directory: %v", err)
		}
		startDir = cwd
	}

	root, err := feature.FindSpecRepoRoot(startDir)
	if err != nil {
		return "", err
	}

	tasksDir := filepath.Join(root, "tasks")
	info, err := os.Stat(tasksDir)
	if err != nil || !info.IsDir() {
		return "", exitcode.NotFoundErrorf("tasks directory not found: %s", tasksDir)
	}
	return tasksDir, nil
}

// --- task list ---

func taskListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks from the board",
		Long: `Lists tasks from the project's tasks/README.md board. Optionally filter
by status. Default output format is YAML.`,
		Args: cobra.NoArgs,
		RunE: runTaskList,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("status", "", "filter by status (e.g., planning, queued, in_progress, completed)")
	cmd.Flags().String("format", "yaml", "output format: yaml, json, md")
	return cmd
}

func runTaskList(cmd *cobra.Command, _ []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")
	statusFlag, _ := cmd.Flags().GetString("status")
	formatFlag, _ := cmd.Flags().GetString("format")

	if formatFlag != "yaml" && formatFlag != "json" && formatFlag != "md" {
		return exitcode.InvalidArgsErrorf("invalid format: %s (supported: yaml, json, md)", formatFlag)
	}

	tasksDir, err := resolveTasksDir(projectFlag)
	if err != nil {
		return err
	}

	boardPath := filepath.Join(tasksDir, "README.md")
	data, err := os.ReadFile(boardPath)
	if err != nil {
		return exitcode.NotFoundErrorf("reading board: %v", err)
	}

	bv, err := task.ParseBoard(data)
	if err != nil {
		return exitcode.UnexpectedErrorf("parsing board: %v", err)
	}

	// Filter by status if requested.
	rows := bv.Rows
	if statusFlag != "" {
		filterStatus := task.TaskStatus(statusFlag)
		var filtered []task.BoardRow
		for _, r := range rows {
			if r.Status == filterStatus {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}

	w := cmd.OutOrStdout()

	switch formatFlag {
	case "yaml":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		if err := enc.Encode(rows); err != nil {
			return exitcode.UnexpectedErrorf("encoding yaml: %v", err)
		}
		return enc.Close()
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(rows)
	case "md":
		filtered := &task.BoardView{Rows: rows}
		_, err := w.Write(task.RenderBoard(filtered))
		return err
	}
	return nil
}

// --- task info ---

func taskInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show detailed task information",
		Long: `Shows detailed information for a task by reading its tasks/{slug}/README.md
file and status from the board. Output includes title, status, description,
dependencies, and summary.`,
		Args: cobra.NoArgs,
		RunE: runTaskInfo,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("task", "", "task slug (required)")
	cmd.Flags().String("format", "yaml", "output format: yaml, json")
	return cmd
}

// taskInfoOutput is the combined output for task info.
type taskInfoOutput struct {
	Slug        string   `yaml:"slug" json:"slug"`
	Title       string   `yaml:"title" json:"title"`
	Status      string   `yaml:"status" json:"status"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	DependsOn   []string `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	Summary     string   `yaml:"summary,omitempty" json:"summary,omitempty"`
}

func runTaskInfo(cmd *cobra.Command, _ []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")
	taskFlag, _ := cmd.Flags().GetString("task")
	formatFlag, _ := cmd.Flags().GetString("format")

	if taskFlag == "" {
		return exitcode.InvalidArgsError("missing required flag: --task")
	}
	if formatFlag != "yaml" && formatFlag != "json" {
		return exitcode.InvalidArgsErrorf("invalid format: %s (supported: yaml, json)", formatFlag)
	}

	tasksDir, err := resolveTasksDir(projectFlag)
	if err != nil {
		return err
	}

	// Read the task file.
	taskFilePath := filepath.Join(tasksDir, taskFlag, "README.md")
	data, err := os.ReadFile(taskFilePath)
	if err != nil {
		return exitcode.NotFoundErrorf("task not found: %s", taskFlag)
	}

	tfd, err := task.ParseTaskFile(data)
	if err != nil {
		return exitcode.UnexpectedErrorf("parsing task file: %v", err)
	}

	// Read status from board.
	status := string(task.StatusPlanning) // default
	boardPath := filepath.Join(tasksDir, "README.md")
	boardData, err := os.ReadFile(boardPath)
	if err == nil {
		bv, parseErr := task.ParseBoard(boardData)
		if parseErr == nil {
			for _, r := range bv.Rows {
				if r.Task == taskFlag {
					status = string(r.Status)
					break
				}
			}
		}
	}

	out := taskInfoOutput{
		Slug:        taskFlag,
		Title:       tfd.Title,
		Status:      status,
		Description: tfd.Description,
		DependsOn:   tfd.DependsOn,
		Summary:     tfd.Summary,
	}

	w := cmd.OutOrStdout()
	switch formatFlag {
	case "yaml":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		if err := enc.Encode(out); err != nil {
			return exitcode.UnexpectedErrorf("encoding yaml: %v", err)
		}
		return enc.Close()
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	return nil
}

// --- task new ---

func taskNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new task",
		Long: `Creates a new task by writing tasks/{slug}/README.md and adding a row
to the tasks/README.md board. New tasks are always created in "planning" status.`,
		Args: cobra.NoArgs,
		RunE: runTaskNew,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("task", "", "task slug (required)")
	cmd.Flags().String("title", "", "task title (required)")
	cmd.Flags().String("description", "", "task description")
	cmd.Flags().String("depends-on", "", "comma-separated list of dependency slugs")
	cmd.Flags().String("format", "yaml", "output format: yaml, json")
	return cmd
}

func runTaskNew(cmd *cobra.Command, _ []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")
	taskFlag, _ := cmd.Flags().GetString("task")
	titleFlag, _ := cmd.Flags().GetString("title")
	descFlag, _ := cmd.Flags().GetString("description")
	depsFlag, _ := cmd.Flags().GetString("depends-on")
	formatFlag, _ := cmd.Flags().GetString("format")

	if taskFlag == "" {
		return exitcode.InvalidArgsError("missing required flag: --task")
	}
	if titleFlag == "" {
		return exitcode.InvalidArgsError("missing required flag: --title")
	}
	if formatFlag != "yaml" && formatFlag != "json" {
		return exitcode.InvalidArgsErrorf("invalid format: %s (supported: yaml, json)", formatFlag)
	}

	// Parse --depends-on
	var deps []string
	if depsFlag != "" {
		for _, d := range strings.Split(depsFlag, ",") {
			d = strings.TrimSpace(d)
			if d != "" {
				deps = append(deps, d)
			}
		}
	}

	tasksDir, err := resolveTasksDir(projectFlag)
	if err != nil {
		return err
	}

	// Check if task already exists.
	taskDir := filepath.Join(tasksDir, taskFlag)
	if _, err := os.Stat(taskDir); err == nil {
		return exitcode.InvalidArgsErrorf("task already exists: %s", taskFlag)
	}

	// Create task directory and README.
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return exitcode.UnexpectedErrorf("creating task directory: %v", err)
	}

	tfd := task.TaskFileData{
		Title:       titleFlag,
		Description: descFlag,
		DependsOn:   deps,
	}
	taskFilePath := filepath.Join(taskDir, "README.md")
	if err := os.WriteFile(taskFilePath, task.RenderTaskFile(tfd), 0o644); err != nil {
		return exitcode.UnexpectedErrorf("writing task file: %v", err)
	}

	// Update board: read, append row, write back.
	boardPath := filepath.Join(tasksDir, "README.md")
	boardData, err := os.ReadFile(boardPath)
	if err != nil {
		return exitcode.NotFoundErrorf("reading board: %v", err)
	}

	bv, err := task.ParseBoard(boardData)
	if err != nil {
		return exitcode.UnexpectedErrorf("parsing board: %v", err)
	}

	newRow := task.BoardRow{
		Task:      taskFlag,
		Status:    task.StatusPlanning,
		DependsOn: deps,
	}
	bv.Rows = append(bv.Rows, newRow)

	if err := os.WriteFile(boardPath, task.RenderBoard(bv), 0o644); err != nil {
		return exitcode.UnexpectedErrorf("writing board: %v", err)
	}

	// Output the created task info.
	out := taskInfoOutput{
		Slug:        taskFlag,
		Title:       titleFlag,
		Status:      string(task.StatusPlanning),
		Description: descFlag,
		DependsOn:   deps,
	}

	w := cmd.OutOrStdout()
	switch formatFlag {
	case "yaml":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		if err := enc.Encode(out); err != nil {
			return exitcode.UnexpectedErrorf("encoding yaml: %v", err)
		}
		return enc.Close()
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}
	return nil
}
