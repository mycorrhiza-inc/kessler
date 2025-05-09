import { unstable_noStore as noStore } from "next/cache";
import { headers } from "next/headers";
import { getUniversalEnvConfig } from "./env_variables";

export default function EnvVariablesScript() {
  noStore();

  const nonce = headers().get("x-nonce");

  return (
    <script
      id="env-config"
      nonce={nonce || ""}
      dangerouslySetInnerHTML={{
        __html: JSON.stringify(getUniversalEnvConfig()),
      }}
    />
  );
}
