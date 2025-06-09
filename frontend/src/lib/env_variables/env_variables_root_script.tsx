import { unstable_noStore as noStore } from "next/cache";
import { headers } from "next/headers";
import { EnvScript } from "next-runtime-env";
import { getEnvConfig } from "./env_variables";

/**
 * Embeds runtime environment configuration into the HTML.
 * Uses next-runtime-env under the hood.
 */
export default async function EnvVariablesScript() {
  noStore();

  const nonce = (await headers()).get("x-nonce") || undefined;

  const env_obj = getEnvConfig();

  return <EnvScript nonce={nonce} env={env_obj} />;
}

