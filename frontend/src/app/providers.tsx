"use client";
import { getRuntimeEnv } from "@/lib/env_variables_hydration_script";
import posthog from "posthog-js";
import { PostHogProvider } from "posthog-js/react";

if (typeof window !== "undefined") {
  const runtimeConfig = getRuntimeEnv();
  if (!runtimeConfig.public_posthog_key && !runtimeConfig.public_posthog_host) {
    posthog.init(runtimeConfig.public_posthog_key!, {
      api_host: runtimeConfig.public_posthog_host,
      person_profiles: "identified_only",
      capture_pageview: true, // Disable automatic pageview capture, as we capture manually
    });
  }
}

export function PHProvider({ children }: { children: React.ReactNode }) {
  return <PostHogProvider client={posthog}>{children}</PostHogProvider>;
}
