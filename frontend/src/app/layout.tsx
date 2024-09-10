import HeaderAuth from "@/components/supabasetutorial/header-auth";
import { ThemeSwitcher } from "@/components/supabasetutorial/theme-switcher";
import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import Link from "next/link";
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
          <main className="h-100vh flex flex-col items-center">
            <div className="flex-1 w-100vw h-100vw flex flex-col items-center">
              <nav className="w-full flex justify-center border-b border-b-foreground/10 h-16">
                <div className="w-full max-w-5xl flex justify-between items-center p-3 px-5 text-sm" style={{zIndex: 3000}}>
                  <div className="flex gap-5 items-center font-semibold">
                    <Link href={"/"}>Kessler</Link>
                  </div>
                  <HeaderAuth />
                  <ThemeSwitcher />
                </div>
              </nav>
              <div className="flex flex-col">
                {children}
              </div>
            </div>
          </main>
        </ThemeProvider>
      </body>
    </html>
  );
}
