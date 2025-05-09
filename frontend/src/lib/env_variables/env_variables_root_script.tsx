import { unstable_noStore as noStore } from "next/cache";
import { headers } from "next/headers";
import Script from "next/script"; // Use Next.js Script component for CSP compliance
import { getUniversalEnvConfig } from "./env_variables";

/**
 * Embeds runtime environment configuration into the HTML.
 * Uses next-runtime-env under the hood.
 */
export default function EnvVariablesScript() {
  noStore();

  const nonce = headers().get("x-nonce") || undefined;
  const envString = JSON.stringify(getUniversalEnvConfig());

  return (
    <Script
      id="env-config"
      nonce={nonce}
      strategy="beforeInteractive"
      dangerouslySetInnerHTML={{ __html: envString }}
    />
  );
}