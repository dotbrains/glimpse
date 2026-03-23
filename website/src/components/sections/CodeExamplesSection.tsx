'use client';

import React, { useState } from 'react';
import { CodeBlock } from '@/components/CodeBlock';

export function CodeExamplesSection() {
  const [activeTab, setActiveTab] = useState<'working' | 'branches' | 'commits' | 'multi'>('working');

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
    commits: `# View last commit
$ glimpse HEAD~1

# View last 3 commits
$ glimpse HEAD~3

# Changes since a specific commit
$ glimpse abc1234`,
    multi: `# Terminal 1 — starts on :5391
$ cd ~/projects/app && glimpse

# Terminal 2 — starts on :5392
$ cd ~/projects/api && glimpse

# List all running instances
$ glimpse list
PORT   PID     REPO                                     REFS
────────────────────────────────────────────────────────────────────────────────
5391   12345   /Users/dev/projects/app                  working tree
5392   12346   /Users/dev/projects/api                  main..feature

# Stop existing and start fresh
$ glimpse --new`,
  };

  const tabs = [
    { key: 'working' as const, label: 'Working Tree', language: 'bash' },
    { key: 'branches' as const, label: 'Branches', language: 'bash' },
    { key: 'commits' as const, label: 'Commits', language: 'bash' },
    { key: 'multi' as const, label: 'Multi-Instance', language: 'bash' },
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
