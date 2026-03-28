# /glimpse-resolve

Reads all open comments from the glimpse diff viewer and makes the requested code changes.

## Usage

```
/glimpse-resolve                       # resolve all open comments
/glimpse-resolve abc123                # resolve a specific comment by ID
```

## Implementation

Run the resolve command to get all open comments:

```bash
glimpse resolve
```

This outputs comments in the format:
```
file:line [severity] comment text (id: abc123)
```

For each comment:
1. Read the file referenced in the comment.
2. Navigate to the specified line.
3. Make the code change described in the comment body.
4. After making changes, the comment will be resolved in the viewer.

If a specific ID is provided:
```bash
glimpse resolve abc123
```

Only resolve that single comment.
