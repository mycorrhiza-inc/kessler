"use client";

import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import {
  EnvConfig,
  getEnvConfig,
} from "./env_variables";
import MarkdownRenderer from "@/components/MarkdownRenderer";

// Client-specific alias for getting runtime env config
// This is now provided by env_variables.ts
// export const getClientRuntimeEnv = (): EnvConfig => {
//   return getEnvConfig();
// };

// Initialize context with an empty validated config
const EnvVariablesClientContext =
  createContext<EnvConfig>(getEnvConfig());

type EnvClientProviderProps = {
  children: ReactNode;
};

export const EnvVariablesClientProvider: React.FC<EnvClientProviderProps> = ({
  children,
}) => {
  const [envs, setEnvs] = useState<EnvConfig>(getEnvConfig);

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

export const useEnvVariablesClientConfig = (): EnvConfig => {
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
`;
  return <MarkdownRenderer>{markdown_string}</MarkdownRenderer>;
};

export function getClientRuntimeEnv(): EnvConfig {
  return getEnvConfig();
}

