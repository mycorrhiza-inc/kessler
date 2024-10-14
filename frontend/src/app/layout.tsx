import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import "./globals.css";
const defaultUrl = "https://kessler.xyz";

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
  // const { theme } = useTheme();
  // console.log(theme);
  // const bg_color_tw = theme === "light" ? "bg-white" : "bg-black";
  return (
    <html
      lang="en"
      className={GeistSans.className}
      // style={{ background-color: "oklch(var(--b1))" }
      suppressHydrationWarning
    >
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      </head>
      <body className="bg-base-100">
        <ThemeProvider>
          <main className="">
            <div className="flex-1 w-100vw flex flex-col items-center">
              {children}
            </div>
          </main>
        </ThemeProvider>
      </body>
    </html>
  );
}
