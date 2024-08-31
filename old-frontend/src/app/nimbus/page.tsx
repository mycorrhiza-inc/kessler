"use client";
import ModelPageBrowser from "../../lib/components/ModelTable";
import DefaultShell from "../../lib/components/DefaultShell";

export default function Page() {
  return (
    <DefaultShell>
      <ModelPageBrowser
        modelUrl="https://nimbus.kessler.xyz/api/v1/models/all"
        data={null}
      ></ModelPageBrowser>
    </DefaultShell>
  );
}
