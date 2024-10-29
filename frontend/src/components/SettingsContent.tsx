import { User } from "@supabase/supabase-js";
import ThemeSelector from "./ThemeSelector";

import { createClient } from "@/utils/supabase/server";
import { signOutAction } from "@/app/actions";
import { useBearStore } from "@/lib/store";
// The password reset is horribly insecure, but it was horribly insecure before and did allow a password reset with a stolen cookie, but now there is a button that does the same thing. Welp...
const SettingsContent = () => {
  const globalStore = useBearStore();
  return (
    <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box w-full">
      <h1 className="text-5xl font-extrabold">Settings</h1>
      {/* Sign out broken, due to me not knowing how to do server components. */}
      {/* <button */}
      {/*   className="btn btn-outline btn-secondary" */}
      {/*   onClick={signOutAction} */}
      {/* ></button> */}
      <a href="/sign-out" className="btn btn-outline btn-secondary">
        Sign Out
      </a>
      <a
        href="/protected/reset-password"
        className="btn btn-outline btn-secondary"
      >
        Reset Password
      </a>
      <ThemeSelector />
    </div>
  );
};
export default SettingsContent;
