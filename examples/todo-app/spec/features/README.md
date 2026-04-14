# Todo App — Features

A demo CLI todo application specified using all four SpecScore layers: Features, Requirements, Acceptance Criteria, and Scenarios.

## CLI Interface

```
todo add <title> [--due YYYY-MM-DD]
todo edit <id> [--title <text>] [--due YYYY-MM-DD]
todo complete <id>
todo reopen <id>
todo delete <id>
todo list [--all | --completed | --overdue]
```

| Feature | Status | Description |
|---------|--------|-------------|
| [todo-item](todo-item/README.md) | Draft | Todo entity definition, management, and completion lifecycle |
| [todo-list](todo-list/README.md) | Draft | Listing, filtering, and count summaries |
| [due-dates](due-dates/README.md) | Draft | Optional due dates and overdue detection |

## Feature Hierarchy

```
spec/features/
├── todo-item/              # The todo entity
│   ├── manage/            # Create, edit, delete operations
│   └── completion/        # Complete/reopen lifecycle
├── todo-list/             # Listing and filtering
└── due-dates/             # Due dates and overdue detection
```

## Outstanding Questions

None at this time.

---
*This document follows the https://specscore.md/feature-specification*
