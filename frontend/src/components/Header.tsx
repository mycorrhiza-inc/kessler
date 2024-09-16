import { signOutAction } from "@/app/actions";
import { Button } from "@/components/supabasetutorial/ui/button";
import { createClient } from "@/utils/supabase/server";
import ThemeSelector from "./ThemeSelector";

async function HeaderAuth() {
  const {
    data: { user },
  } = await createClient().auth.getUser();

  return user ? (
    <div className="flex items-center gap-4">
      Hey, {user.email}!
      <form action={signOutAction} method="post">
        <button type="submit" className="btn btn-outline btn-secondary">
          Sign out
        </button>
      </form>
    </div>
  ) : (
    <div className="flex gap-2">
      <a href="/sign-in" className="btn btn-outline btn-secondary">
        Sign in
      </a>
      <a href="/sign-up" className="btn btn-outline btn-secondary">
        Sign up
      </a>
    </div>
  );
}
const Header = () => {
  return (
    <nav className="w-full flex justify-center border-b border-b-foreground/10 h-16 bg-base-100 text-base-content">
      <div
        className="w-full max-w-5xl flex justify-between items-center p-3 px-5 text-sm"
        style={{ zIndex: 3000 }}
      >
        <div className="flex gap-5 items-center font-semibold">
          <a href="/">Kessler</a>
        </div>
        <HeaderAuth />
        <ThemeSelector />
      </div>
    </nav>
  );
};
export default Header;
