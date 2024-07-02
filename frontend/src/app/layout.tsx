// layout.tsx
import { Metadata } from "next";
import { Inter } from "next/font/google";
import { ClerkProvider } from "@clerk/nextjs";
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
    <ClerkProvider>
      <html lang="en">
        <body>
          <ColorModeScript initialColorMode={theme.config.initialColorMode} />
          <ChakraProvider theme={theme}>
            <DarkMode>
              <SaasProvider>{children}</SaasProvider>
            </DarkMode>
          </ChakraProvider>
        </body>
      </html>
    </ClerkProvider>
  );
}

// import type { Metadata } from "next";
// import { Inter } from "next/font/google";
// import "./globals.css";
// import { ClerkProvider } from "@clerk/nextjs";
// import { SaasProvider } from "@saas-ui/react";
//
// const inter = Inter({ subsets: ["latin"] });
//
// export const metadata: Metadata = {
//   title: "Kessler Search",
//   description: "the kessler search engine",
// };
//
// export default function RootLayout({
//   children,
// }: {
//   children: React.ReactNode;
// }) {
//   return (
//     <ClerkProvider>
//       <html lang="en">
//         <body>
//           <SaasProvider>{children}</SaasProvider>
//         </body>
//       </html>
//     </ClerkProvider>
//   );
// }
