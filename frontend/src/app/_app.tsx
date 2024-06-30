// src/pages/_app.tsx

import { Chakra } from "../lib/Chakra";
import { AppProps } from "next/app";
import { ChakraProvider, ColorModeScript } from "@chakra-ui/react";
import theme from "../app/theme";
import "../app/globals.css";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <Chakra cookies={pageProps.cookies}>
      <ChakraProvider theme={theme}>
        <ColorModeScript initialColorMode={theme.config.initialColorMode} />
        <Component {...pageProps} />
      </ChakraProvider>
    </Chakra>
  );
}

export default MyApp;
export { getServerSideProps } from "../lib/Chakra";
