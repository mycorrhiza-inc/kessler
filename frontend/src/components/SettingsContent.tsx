import ThemeSelector from "./ThemeSelector";

import { useKesslerStore } from "@/lib/_store";
import Link from "next/link";
// The password reset is horribly insecure, but it was horribly insecure before and did allow a password reset with a stolen cookie, but now there is a button that does the same thing. Welp...
const SettingsContent = () => {
  const globalStore = useKesslerStore();
  const experimentalFeaturesEnabled = globalStore.experimentalFeaturesEnabled;
  const setExperimentalFeatures = globalStore.setExperimentalFeaturesEnabled;
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
      <div>
        <p>
          <br />
        </p>
      </div>
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
            // Also right now this changes on every page refresh, (but stays intact with react router navigations) so setting the global store value from user info is good!
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
