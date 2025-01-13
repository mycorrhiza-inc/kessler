"use client";
import PageContainer from "@/components/Page/PageContainer";
import {
  EnvVariablesClientProvider,
  EnvironmentVariableTestMarkdown,
} from "@/lib/env_variables_hydration_script";
export const dynamic = "force-dynamic";
export default function Page() {
  return (
    <>
      <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
        <EnvVariablesClientProvider>
          <EnvironmentVariableTestMarkdown />
        </EnvVariablesClientProvider>
      </PageContainer>
    </>
  );
}
