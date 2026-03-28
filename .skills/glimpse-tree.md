# /glimpse-tree

Opens a full file tree browser for the current repo — no diff required. Browse files with syntax highlighting and leave comments.

## Usage

```
/glimpse-tree
```

## Implementation

Run:

```bash
glimpse tree
```

This opens a file tree browser in the user's browser. They can click any file to view it with syntax highlighting, and leave comments on specific lines. When done, they can run `/glimpse-resolve-tree` to have you apply the changes.
