import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import Header from "@/components/Header";
import "./globals.css";
// import InitColorSchemeScript from "@mui/joy/InitColorSchemeScript";
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
  // const { theme } = useTheme();
  // console.log(theme);
  // const bg_color_tw = theme === "light" ? "bg-white" : "bg-black";
  return (
    <html
      lang="en"
      className={GeistSans.className}
      style={{ backgroundColor: "var(--background)" }}
      suppressHydrationWarning
    >
      <ThemeProvider
        attribute="class"
        defaultTheme="system"
        enableSystem
        disableTransitionOnChange
      >
        <body className="bg-background text-foreground">
          <Header></Header>
          <main className="flex flex-col items-center">
            <div className="flex-1 w-100vw flex flex-col items-center">
              <div className="flex flex-col">{children}</div>
            </div>
          </main>
        </body>
      </ThemeProvider>
    </html>
  );
}
