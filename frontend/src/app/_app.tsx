// src/pages/_app.tsx

import { Chakra } from "../lib/Chakra";
import { AppProps } from "next/app";
import "../app/globals.css";
import { DarkMode } from "@chakra-ui/react";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <Chakra cookies={pageProps.cookies}>
      <DarkMode>
        <Component {...pageProps} />
      </DarkMode>
    </Chakra>
  );
}

export default MyApp;
export { getServerSideProps } from "../lib/Chakra";
