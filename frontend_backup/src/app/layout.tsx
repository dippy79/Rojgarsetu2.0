import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "react-hot-toast";
import { Header } from "@/components/layout/header";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "RojgarSetu 2.0 - AI-Powered Job & Course Aggregator",
  description: "Your Gateway to Government Jobs & Professional Courses with AI-powered matching and gamified career tracking",
  keywords: ["jobs", "courses", "government jobs", "career", "AI matching", "skills"],
  authors: [{ name: "RojgarSetu" }],
  icons: {
    icon: "/favicon.svg",
  },
  openGraph: {
    title: "RojgarSetu 2.0 - AI-Powered Job & Course Aggregator",
    description: "Your Gateway to Government Jobs & Professional Courses",
    type: "website",
    locale: "en_US",
  },
};

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60 * 1000,
      refetchOnWindowFocus: false,
    },
  },
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <QueryClientProvider client={queryClient}>
          <ThemeProvider
            attribute="class"
            defaultTheme="system"
            enableSystem
            disableTransitionOnChange
          >
            <Header />
            <main className="min-h-screen">
              {children}
            </main>
            <Toaster position="top-right" />
          </ThemeProvider>
        </QueryClientProvider>
      </body>
    </html>
  );
}