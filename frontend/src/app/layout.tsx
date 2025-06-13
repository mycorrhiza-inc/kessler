// import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import "./globals.css";
import { PHProvider } from "./providers";
import dynamic from "next/dynamic";
import {
  ClerkProvider,
} from "@clerk/nextjs";
import PostHogPageView from "@/stateful_components/tracking/Posthog/PostHogPageView";
import { Suspense } from "react";

// const PostHogPageView = dynamic(() => import("../components/Posthog/PostHogPageView"), {
//   ssr: false,k
// });

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
      // className={GeistSans.className}
      // style={{ background-color: "oklch(var(--b1))" }
      suppressHydrationWarning
    >
      <PHProvider>
        {/* <ClerkProvider> */}
        <head>
          <meta
            name="viewport"
            content="width=device-width, initial-scale=1.0"
          />
        </head>
        <body className="bg-base-100">
          <ThemeProvider defaultTheme="kessler">
            <Suspense>
              <PostHogPageView />
            </Suspense>
            <main className="">
              <div className="flex-1 w-100vw flex flex-col items-center">
                {/* <Navbar user={user} /> */}
                {children}
              </div>
            </main>
          </ThemeProvider>
        </body>
        {/* </ClerkProvider> */}
      </PHProvider>
    </html>
  );
}
