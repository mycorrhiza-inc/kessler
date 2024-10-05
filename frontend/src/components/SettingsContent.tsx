import { User } from "@supabase/supabase-js";
import ThemeSelector from "./ThemeSelector";

import { createClient } from "@/utils/supabase/server";
// The password reset is horribly insecure, but it was horribly insecure before and did allow a password reset with a stolen cookie, but now there is a button that does the same thing. Welp...
const SettingsContent = ({ user }: { user: User }) => {
  return (
    <div className=" p-5 m-5 justify-center border-2 border-['accent'] rounded-box w-full">
      <h1 className="text-5xl font-extrabold">Settings</h1>

      <button className="btn btn-outline btn-secondary">Sign Out</button>
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
