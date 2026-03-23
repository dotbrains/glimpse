import type { Metadata } from 'next';
import '@/styles/globals.css';

export const metadata: Metadata = {
  title: '__PROJECT_NAME__ — __PROJECT_DESCRIPTION__',
  description: '__PROJECT_DESCRIPTION_LONG__',
  openGraph: {
    title: '__PROJECT_NAME__ — __PROJECT_DESCRIPTION__',
    description: '__PROJECT_DESCRIPTION_LONG__',
    url: 'https://__PROJECT_NAME__.dotbrains.io',
    siteName: '__PROJECT_NAME__',
    images: [
      {
        url: '/og-image.svg',
        width: 1200,
        height: 630,
        alt: '__PROJECT_NAME__ — __PROJECT_DESCRIPTION__',
      },
    ],
    locale: 'en_US',
    type: 'website',
  },
  twitter: {
    card: 'summary_large_image',
    title: '__PROJECT_NAME__ — __PROJECT_DESCRIPTION__',
    description: '__PROJECT_DESCRIPTION_LONG__',
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
