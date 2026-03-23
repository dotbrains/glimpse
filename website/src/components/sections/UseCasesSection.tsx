'use client';

import { GitBranch, Eye, Clock, GitPullRequest, Tag, RefreshCw } from 'lucide-react';

export function UseCasesSection() {
  const useCases = [
    {
      icon: <Eye className="w-6 h-6" />,
      title: 'Pre-Commit Review',
      description: 'Run glimpse before committing to visually review all your uncommitted changes in a proper diff view.',
    },
    {
      icon: <GitBranch className="w-6 h-6" />,
      title: 'Branch Comparison',
      description: 'Compare your feature branch against main before opening a PR. Spot issues early.',
    },
    {
      icon: <Clock className="w-6 h-6" />,
      title: 'Commit History',
      description: 'Review your last N commits with glimpse HEAD~3. See exactly what shipped in a batch of changes.',
    },
    {
      icon: <Tag className="w-6 h-6" />,
      title: 'Release Diffs',
      description: 'Compare two release tags (glimpse v1.0.0 v2.0.0) to see everything that changed between versions.',
    },
    {
      icon: <GitPullRequest className="w-6 h-6" />,
      title: 'AI Agent Output',
      description: 'Review changes made by AI coding agents (Cursor, Claude Code, Codex) in a proper diff view before accepting.',
    },
    {
      icon: <RefreshCw className="w-6 h-6" />,
      title: 'Multi-Repo Workflow',
      description: 'Run glimpse in multiple repos simultaneously. Each gets its own port, existing instances are reused automatically.',
    },
  ];

  return (
    <section id="use-cases" className="py-12 sm:py-16 lg:py-20 bg-dark-slate">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">Use Cases</h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">
            glimpse adapts to how your team works
          </p>
        </div>
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
          {useCases.map((useCase, index) => (
            <div key={index} className="bg-dark-gray/50 border border-accent-primary/20 rounded-xl p-5 sm:p-6 hover:border-accent-secondary/40 transition-all">
              <div className="w-10 h-10 sm:w-12 sm:h-12 bg-gradient-to-br from-accent-primary to-accent-secondary rounded-lg flex items-center justify-center text-white mb-3 sm:mb-4">
                {useCase.icon}
              </div>
              <h3 className="text-lg sm:text-xl font-semibold text-cream mb-2">{useCase.title}</h3>
              <p className="text-cream/60 text-sm sm:text-base leading-relaxed">{useCase.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
