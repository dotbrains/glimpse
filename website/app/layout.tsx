import type { Metadata } from 'next';
import '@/styles/globals.css';

export const metadata: Metadata = {
  title: 'glimpse — GitHub-style git diff viewer CLI',
  description: 'Browser-based, GitHub-style diff viewer for git changes. View uncommitted changes, branch comparisons, commit ranges, and more with syntax-highlighted split diffs.',
  openGraph: {
    title: 'glimpse — GitHub-style git diff viewer CLI',
    description: 'Browser-based, GitHub-style diff viewer for git changes. View uncommitted changes, branch comparisons, commit ranges, and more with syntax-highlighted split diffs.',
    url: 'https://glimpse.dotbrains.io',
    siteName: 'glimpse',
    images: [
      {
        url: '/og-image.svg',
        width: 1200,
        height: 630,
        alt: 'glimpse — GitHub-style git diff viewer CLI',
      },
    ],
    locale: 'en_US',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: 'glimpse — GitHub-style git diff viewer CLI',
    description: 'Browser-based, GitHub-style diff viewer for git changes. View uncommitted changes, branch comparisons, commit ranges, and more with syntax-highlighted split diffs.',
    images: ['/og-image.svg'],
  },
  icons: {
    icon: [
      {
        url: '/favicon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: [
      {
        url: '/favicon.svg',
        type: 'image/svg+xml',
      },
    ],
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </head>
      <body>{children}</body>
    </html>
  );
}
