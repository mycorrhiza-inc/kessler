import { GeistSans } from "geist/font/sans";
import { ThemeProvider } from "next-themes";
import "./globals.css";
import { PHProvider } from "./providers";
import dynamic from "next/dynamic";
import EnvVariablesScript from "@/lib/env_variables/env_variables_root_script";
import {
  ClerkProvider,
} from "@clerk/nextjs";
import { EnvVariablesClientProvider } from "@/lib/env_variables/env_variables_hydration_script";

import PostHogPageView from "@/components/Posthog/PostHogPageView";
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
      className={GeistSans.className}
      // style={{ background-color: "oklch(var(--b1))" }
      suppressHydrationWarning
    >
      <PHProvider>
        <ClerkProvider
          publishableKey="pk_test_YWNlLXdhbGxhYnktODQuY2xlcmsuYWNjb3VudHMuZGV2JA"
        >
          <head>
            <meta
              name="viewport"
              content="width=device-width, initial-scale=1.0"
            />
            <EnvVariablesScript />
          </head>
          <body className="bg-base-100">
            <ThemeProvider defaultTheme="kessler">
              <PostHogPageView />
              <main className="">
                <div className="flex-1 w-100vw flex flex-col items-center">
                  {/* <Navbar user={user} /> */}
                  <EnvVariablesClientProvider>
                    {children}
                  </EnvVariablesClientProvider>
                </div>
              </main>
            </ThemeProvider>
          </body>
        </ClerkProvider>
      </PHProvider>
    </html>
  );
}
