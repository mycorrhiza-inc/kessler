import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import Header from "@/components/Header";
import "./globals.css";

const defaultUrl = "https://app.kessler.xyz";

export const metadata = {
  metadataBase: new URL(defaultUrl),
  title: "Kessler",
  description: "Inteligence and Research tools for Lobbying",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={GeistSans.className} suppressHydrationWarning>
      <body className="bg-background text-foreground">
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <Header></Header>
          <main className="h-100vh flex flex-col items-center">
            <div className="flex-1 w-100vw h-100vw flex flex-col items-center">
              <div className="flex flex-col">{children}</div>
            </div>
          </main>
        </ThemeProvider>
      </body>
    </html>
  );
}
