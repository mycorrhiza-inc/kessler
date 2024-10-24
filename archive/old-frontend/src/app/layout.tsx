// layout.tsx
// "use client";
import { Metadata } from "next";
import { Inter } from "next/font/google";
import { SaasProvider } from "@saas-ui/react";
import { ChakraProvider, ColorModeScript, DarkMode } from "@chakra-ui/react";
import theme from "../app/theme";
import "../app/globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Kessler Search",
  description: "the kessler search engine",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ColorModeScript initialColorMode={theme.config.initialColorMode} />
        <ChakraProvider>
          <SaasProvider>{children}</SaasProvider>
        </ChakraProvider>
      </body>
    </html>
  );
}
