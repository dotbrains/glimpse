# /glimpse-review

Runs an AI code review on the current diff and posts severity-tagged inline comments to the glimpse viewer.

## Usage

```
/glimpse-review                             # review working tree changes
/glimpse-review main                        # review what you're merging into main
/glimpse-review main..feature               # review branch diff
/glimpse-review identify security issues    # focus on security
/glimpse-review performance in src/lib      # focus on performance in a path
/glimpse-review last 3 commits              # natural language
```

## Implementation

First ensure a glimpse instance is running (run `glimpse` if not).

Then run the review command:

```bash
glimpse review {{args}}
```

If the user specifies a focus area (security, performance, testing, etc.), add `--focus`:

```bash
glimpse review --focus security
```

After the review completes, tell the user how many comments were posted and that they can view them in the browser. Suggest running `/glimpse-resolve` to apply fixes.
