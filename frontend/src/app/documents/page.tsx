"use client";
import FileBrowser from "../../lib/components/FileBrowser";
import DefaultShell from "../../lib/components/DefaultShell";

import { GetAllFiles } from "../../lib/requests";

export default function Page() {
  return (
    <DefaultShell>
      <FileBrowser getFileFunc={GetAllFiles}/>
    </DefaultShell>
  );
}
