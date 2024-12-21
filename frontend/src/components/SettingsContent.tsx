import ThemeSelector from "./ThemeSelector";

import { useKesslerStore } from "@/lib/store";
import Link from "next/link";
// The password reset is horribly insecure, but it was horribly insecure before and did allow a password reset with a stolen cookie, but now there is a button that does the same thing. Welp...
const SettingsContent = () => {
  const globalStore = useKesslerStore();
  const experimentalFeaturesEnabled = globalStore.enableExperimentalFeatures;
  const setExperimentalFeatures = globalStore.setEnableExperimentalFeatures;
  return (
    <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box w-full">
      <h1 className="text-5xl font-extrabold">Settings</h1>
      <br />
      <Link href="/sign-out" className="btn btn-outline btn-secondary">
        Sign Out
      </Link>
      <Link
        href="/protected/reset-password"
        className="btn btn-outline btn-secondary"
      >
        Reset Password
      </Link>
      <ThemeSelector />
      <h2 className="text-3xl font-bold">Options</h2>
      <div className="form-control">
        <label className="cursor-pointer label">
          <span className="label-text">Enable Experimental Features</span>
          <input
            type="checkbox"
            className="checkbox checkbox-success"
            checked={experimentalFeaturesEnabled}
            // this seems hackish, but it kinda works?
            onChange={() =>
              setExperimentalFeatures(!experimentalFeaturesEnabled)
            }
          />
        </label>
      </div>
    </div>
  );
};
export default SettingsContent;
