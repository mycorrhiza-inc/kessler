"use client";

import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";
import { RuntimeEnvConfig } from "./env_variables";
import { env } from "process";

type EnvVariablesClientConfig = {};

const defaultEnvVariables = {};

const EnvVariablesClientContext =
  createContext<RuntimeEnvConfig>(defaultEnvVariables);

type EnvClientProviderProps = {
  children: ReactNode;
};

export const EnvVariablesClientProvider: React.FC<EnvClientProviderProps> = ({
  children,
}) => {
  const [envs, setEnvs] =
    useState<EnvVariablesClientConfig>(defaultEnvVariables);

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

export const useEnvVariablesClientConfig = (): EnvVariablesClientConfig => {
  if (EnvVariablesClientContext === undefined) {
    throw new Error(
      "useEnvVariablesClientConfig must be used within an EnvVariablesClientProvider",
    );
  }

  return useContext(EnvVariablesClientContext);
};

export const envScriptId = "env-config";

const isSSR = typeof window === "undefined";

export const getRuntimeEnv = (): EnvVariablesClientConfig => {
  if (isSSR) return env;
  const script = window.document.getElementById(
    envScriptId,
  ) as HTMLScriptElement;

  return script ? JSON.parse(script.innerText) : undefined;
};
