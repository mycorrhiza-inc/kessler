"use client";

import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import {
  RuntimeEnvConfig,
  emptyRuntimeConfig,
  getUniversalEnvConfig,
} from "./env_variables";
import MarkdownRenderer from "@/components/MarkdownRenderer";

// Client-specific alias for getting runtime env config
// This is now provided by env_variables.ts
// export const getClientRuntimeEnv = (): RuntimeEnvConfig => {
//   return getUniversalEnvConfig();
// };

// Initialize context with an empty validated config
const EnvVariablesClientContext =
  createContext<RuntimeEnvConfig>(emptyRuntimeConfig);

type EnvClientProviderProps = {
  children: ReactNode;
};

export const EnvVariablesClientProvider: React.FC<EnvClientProviderProps> = ({
  children,
}) => {
  const [envs, setEnvs] = useState<RuntimeEnvConfig>(emptyRuntimeConfig);

  useEffect(() => {
    const runtimeEnvs = getClientRuntimeEnv();
    setEnvs(runtimeEnvs);
  }, []);

  return (
    <EnvVariablesClientContext.Provider value={envs}>
      {children}
    </EnvVariablesClientContext.Provider>
  );
};

export const useEnvVariablesClientConfig = (): RuntimeEnvConfig => {
  const context = useContext(EnvVariablesClientContext);
  if (context === undefined) {
    throw new Error(
      "useEnvVariablesClientConfig must be used within an EnvVariablesClientProvider",
    );
  }

  return context;
};

export const EnvironmentVariableTestMarkdown = () => {
  const config = useEnvVariablesClientConfig();
  const markdown_string = `# Environment Variables
INTERNAL_API_URL: ${config.internal_api_url}

PUBLIC_API_URL: ${config.public_api_url}

NEXT_PUBLIC_POSTHOG_KEY: ${config.public_posthog_key}

NEXT_PUBLIC_POSTHOG_HOST: ${config.public_posthog_host}

VERSION_HASH: ${config.version_hash}
`;
  return <MarkdownRenderer>{markdown_string}</MarkdownRenderer>;
};

export function getClientRuntimeEnv(): RuntimeEnvConfig {
  return getUniversalEnvConfig();
}

