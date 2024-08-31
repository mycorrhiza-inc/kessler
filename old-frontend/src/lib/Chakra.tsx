import React, { ReactNode } from "react";
import {
  ChakraProvider,
  cookieStorageManagerSSR,
  localStorageManager,
} from "@chakra-ui/react";
import { GetServerSideProps } from "next";

// Define the props type for the `Chakra` component
interface ChakraProps {
  cookies: string;
  children: ReactNode;
}

// Functional component with proper types for props
export function Chakra({ cookies, children }: ChakraProps) {
  const colorModeManager =
    typeof cookies === "string"
      ? cookieStorageManagerSSR(cookies)
      : localStorageManager;

  return (
    <ChakraProvider colorModeManager={colorModeManager}>
      {children}
    </ChakraProvider>
  );
}

// Type for getServerSideProps context
interface ServerSideContext {
  req: {
    headers: {
      cookie?: string;
    };
  };
}

// Define getServerSideProps with proper types
export const getServerSideProps: GetServerSideProps = async ({
  req,
}: ServerSideContext) => {
  return {
    props: {
      cookies: req.headers.cookie ?? "",
    },
  };
};
