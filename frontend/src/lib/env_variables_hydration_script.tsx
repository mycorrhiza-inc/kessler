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
  runtimeConfig,
} from "./env_variables";
import { env } from "process";
import MarkdownRenderer from "@/components/MarkdownRenderer";

const defaultEnvVariables = {};

const EnvVariablesClientContext =
  createContext<RuntimeEnvConfig>(defaultEnvVariables);

type EnvClientProviderProps = {
  children: ReactNode;
};

export const EnvVariablesClientProvider: React.FC<EnvClientProviderProps> = ({
  children,
}) => {
  const [envs, setEnvs] = useState<RuntimeEnvConfig>(defaultEnvVariables);

  useEffect(() => {
    const runtimeEnvs = getRuntimeEnv();
    setEnvs(runtimeEnvs);
  }, []);

  return (
    <EnvVariablesClientContext.Provider value={envs}>
      {children}
    </EnvVariablesClientContext.Provider>
  );
};

export const useEnvVariablesClientConfig = (): RuntimeEnvConfig => {
  if (EnvVariablesClientContext === undefined) {
    throw new Error(
      "useEnvVariablesClientConfig must be used within an EnvVariablesClientProvider",
    );
  }

  return useContext(EnvVariablesClientContext);
};

export const envScriptId = "env-config";

const isSSR = typeof window === "undefined";

export const getRuntimeEnv = (): RuntimeEnvConfig => {
  if (isSSR) {
    throw new Error(
      "You must call this function in a client component, for a server component just import runtimeConfig from ./env_variables.ts",
    );
  }
  const script = window.document.getElementById(
    envScriptId,
  ) as HTMLScriptElement;

  return script ? JSON.parse(script.innerText) : emptyRuntimeConfig;
};

export const EnvironmentVariableTestMarkdown = () => {
  const config = useEnvVariablesClientConfig();
  const markdown_string = `# Environment Variables
INTERNAL_API_URL: ${config.internal_api_url}

PUBLIC_API_URL: ${config.public_api_url}

PUBLIC_POSTHOG_KEY: ${config.public_posthog_key}

PUBLIC_POSTHOG_HOST: ${config.public_posthog_host}
`;
  return <MarkdownRenderer>{markdown_string}</MarkdownRenderer>;
};
