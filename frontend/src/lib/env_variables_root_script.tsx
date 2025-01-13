import { unstable_noStore as noStore } from "next/cache";
import { runtimeEnvConfig } from "./env_variables";
import { headers } from "next/headers";

export default function EnvVariablesScript() {
  noStore();

  const nonce = headers().get("x-nonce");

  return (
    <script
      id="env-config"
      nonce={nonce || ""}
      dangerouslySetInnerHTML={{
        __html: JSON.stringify(runtimeEnvConfig),
      }}
    />
  );
}
