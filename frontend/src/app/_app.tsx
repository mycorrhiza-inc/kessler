// src/pages/_app.tsx

import { Chakra } from "../lib/Chakra";
import { AppProps } from "next/app";
import { ChakraProvider, ColorModeScript } from "@chakra-ui/react";
import theme from "../app/theme";
import "../app/globals.css";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <Chakra cookies={pageProps.cookies}>
      <Component {...pageProps} />
    </Chakra>
  );
}

export default MyApp;
export { getServerSideProps } from "../lib/Chakra";
