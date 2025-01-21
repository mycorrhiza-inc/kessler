"use client";
import {
  EnvVariablesClientProvider,
  EnvironmentVariableTestMarkdown,
} from "@/lib/env_variables_hydration_script";
export const dynamic = "force-dynamic";
export default function Page() {
  return (
    <>
      <EnvVariablesClientProvider>
        <EnvironmentVariableTestMarkdown />
      </EnvVariablesClientProvider>
    </>
  );
}
