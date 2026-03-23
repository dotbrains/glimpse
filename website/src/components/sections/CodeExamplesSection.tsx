'use client';

import React, { useState } from 'react';
import { CodeBlock } from '@/components/CodeBlock';

export function CodeExamplesSection() {
  const [activeTab, setActiveTab] = useState<'working' | 'branches' | 'pr' | 'review' | 'resolve'>('working');

  const examples = {
    working: `# View all uncommitted changes (staged + unstaged)
$ glimpse
→ 3 files changed, 12 insertions(+), 4 deletions(-)
→ Serving at http://localhost:5391`,
    branches: `# Compare current branch against main
$ glimpse main

# Compare two branches
$ glimpse main..feature
$ glimpse main feature
$ glimpse --base main --compare feature

# Compare tags
$ glimpse v1.0.0 v2.0.0`,
    pr: `# View a GitHub PR
$ glimpse https://github.com/owner/repo/pull/123
→ Fetching PR #123 from owner/repo...
→ PR #123: Fix authentication race condition
→ 12 files changed, 45 insertions(+), 18 deletions(-)
→ Serving at http://localhost:5391`,
    review: `# AI code review of working tree changes
$ glimpse review
→ Reviewing 3 files changed, 12 insertions(+), 4 deletions(-)...
✓ 8 comments posted to http://localhost:5391
  2 must-fix
  4 suggestion
  2 nit

# Focus on security issues
$ glimpse review --focus security

# Review specific refs
$ glimpse review main..feature`,
    resolve: `# Output all open comments for your agent
$ glimpse resolve
main.go:42 [must-fix] Nil pointer dereference (id: a1b2c3)
auth.go:18 [suggestion] Consider using sync.Once (id: d4e5f6)

# Resolve a specific comment
$ glimpse resolve a1b2c3`,
  };

  const tabs = [
    { key: 'working' as const, label: 'Diffs', language: 'bash' },
    { key: 'branches' as const, label: 'Branches', language: 'bash' },
    { key: 'pr' as const, label: 'GitHub PRs', language: 'bash' },
    { key: 'review' as const, label: 'AI Review', language: 'bash' },
    { key: 'resolve' as const, label: 'Resolve', language: 'bash' },
  ];

  return (
    <section id="code-examples" className="py-12 sm:py-16 lg:py-20 bg-dark-gray/50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">Code Examples</h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">
            See glimpse in action
          </p>
        </div>
        <div className="bg-dark-slate border border-accent-primary/30 rounded-xl overflow-hidden">
          <div className="flex border-b border-accent-primary/30 overflow-x-auto">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`flex-1 px-3 sm:px-6 py-3 sm:py-4 text-xs sm:text-sm font-semibold transition-colors whitespace-nowrap ${
                  activeTab === tab.key
                    ? 'bg-dark-gray/50 text-accent-primary border-b-2 border-accent-primary'
                    : 'text-cream/70 hover:text-cream hover:bg-dark-gray/30'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
          <div className="p-4 sm:p-6 overflow-x-auto">
            <CodeBlock
              code={examples[activeTab]}
              language={tabs.find((t) => t.key === activeTab)?.language}
            />
          </div>
        </div>
      </div>
    </section>
  );
}
