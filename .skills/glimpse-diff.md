# /glimpse-diff

Opens the glimpse diff viewer in the user's browser. Accepts the same refs as the CLI, plus natural language.

## Usage

```
/glimpse-diff                          # working tree changes
/glimpse-diff main                     # current branch against main
/glimpse-diff main..feature            # branch diff
/glimpse-diff HEAD~1                   # last commit
/glimpse-diff last 3 commits           # natural language
/glimpse-diff https://github.com/owner/repo/pull/123  # PR URL
```

## Implementation

Run the `glimpse` command with the user's arguments:

```bash
glimpse {{args}}
```

If no arguments are provided, run `glimpse` with no args to show working tree changes.

After running, tell the user the viewer is open and they can leave comments on any line. When they're done commenting, they can run `/glimpse-resolve` to have you fix the issues.
