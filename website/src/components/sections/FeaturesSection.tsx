'use client';

import { GitBranch, Eye, Layers, Globe, Terminal, Zap, GitPullRequest, Bot, MessageSquare, CheckCircle } from 'lucide-react';

export function FeaturesSection() {
  const features = [
    {
      icon: <Eye className="w-6 h-6" />,
      title: 'GitHub-Style Diffs',
      description: 'Syntax-highlighted split diff view with addition/deletion coloring, line numbers, hunk headers, and file status badges.',
    },
    {
      icon: <GitPullRequest className="w-6 h-6" />,
      title: 'GitHub PRs',
      description: 'Pass a PR URL to view any pull request locally. Fetches the diff via gh CLI — no cloning needed.',
    },
    {
      icon: <MessageSquare className="w-6 h-6" />,
      title: 'Inline Comments',
      description: 'Click any line to leave comments with severity tags: must-fix, suggestion, nit, question. Persisted to disk.',
    },
    {
      icon: <Bot className="w-6 h-6" />,
      title: 'AI Code Review',
      description: 'Run glimpse review to have an AI agent review the diff and post severity-tagged inline comments to the viewer.',
    },
    {
      icon: <CheckCircle className="w-6 h-6" />,
      title: 'Resolve Workflow',
      description: 'Run glimpse resolve to output open comments for your AI agent. Review → check → resolve, all from the terminal.',
    },
    {
      icon: <GitBranch className="w-6 h-6" />,
      title: 'Any Git Ref',
      description: 'Branches, tags, commits, ranges (main..feature), HEAD~N — glimpse resolves them all.',
    },
    {
      icon: <Terminal className="w-6 h-6" />,
      title: 'Working Tree Diffs',
      description: 'Run glimpse with no args to view all uncommitted changes — both staged and unstaged.',
    },
    {
      icon: <Layers className="w-6 h-6" />,
      title: 'Multi-Instance',
      description: 'Run multiple repos simultaneously. Each gets its own auto-assigned port. Re-running opens the existing instance.',
    },
    {
      icon: <Zap className="w-6 h-6" />,
      title: 'Single Binary',
      description: 'Written in Go. All assets embedded. Install via go install, Homebrew, or download from GitHub Releases.',
    },
  ];

  return (
    <section id="features" className="py-12 sm:py-16 lg:py-20 bg-dark-slate">
      <div className="max-w-7xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-10 sm:mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-bold text-cream mb-3 sm:mb-4">
            Built for Real Engineering Workflows
          </h2>
          <p className="text-cream/70 text-base sm:text-lg lg:text-xl max-w-3xl mx-auto">
            glimpse fits into how you already work
          </p>
        </div>
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
          {features.map((feature, index) => (
            <div
              key={index}
              className="group bg-dark-gray/50 border border-accent-primary/20 hover:border-accent-secondary/40 rounded-xl p-5 sm:p-6 transition-all hover:shadow-lg hover:shadow-accent-primary/10"
            >
              <div className="w-10 h-10 sm:w-12 sm:h-12 bg-gradient-to-br from-accent-primary to-accent-secondary rounded-lg flex items-center justify-center text-white mb-3 sm:mb-4 group-hover:scale-110 transition-transform">
                {feature.icon}
              </div>
              <h3 className="text-lg sm:text-xl font-semibold text-cream mb-2">{feature.title}</h3>
              <p className="text-cream/60 text-sm sm:text-base leading-relaxed">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
