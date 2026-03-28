# /glimpse-resolve-tree

Reads all open comments from the glimpse tree browser and makes the requested code changes.

## Usage

```
/glimpse-resolve-tree                  # resolve all open comments
/glimpse-resolve-tree abc123           # resolve a specific comment by ID
```

## Implementation

Run:

```bash
glimpse resolve-tree
```

This outputs comments in the same format as `/glimpse-resolve`:
```
file:line [severity] comment text (id: abc123)
```

For each comment, read the referenced file, make the described change, and move on. For a specific comment:

```bash
glimpse resolve-tree abc123
```
