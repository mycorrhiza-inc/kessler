import { signOutAction } from "@/app/actions";
import { createClient } from "@/utils/supabase/server";
import { UserIcon } from "@/components/Icons";

async function HeaderAuth() {
  const {
    data: { user },
  } = await createClient().auth.getUser();

  return user ? (
    <div className="flex items-center gap-4">
      <div class="dropdown dropdown-end">
        <div tabindex={0} role="button" class="btn btn-ghost rounded-btn">
          <UserIcon />
        </div>
        <form action={signOutAction} method="post">
          <ul
            tabindex={0}
            class="menu dropdown-content bg-base-200 rounded-box z-[1] w-52 p-2 ">
            <li>Hey, {user.email}!</li>
            <li><a href="/settings">Settings</a></li>
            <li>
              <button type="submit" className="btn btn-outline btn-secondary">
                Sign out
              </button>
            </li>
          </ul>
        </form>
      </div>
    </div >
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
    <nav
      className="w-full flex justify-center border-b border-b-foreground/10 h-16 bg-base-200 text-base-content"
      style={{ zIndex: 3000 }}
    >
      <div
        className="w-full max-w-5xl flex justify-between items-center bg-base-200 p-3 px-5 text-sm"
        style={{ zIndex: 3000 }}
      >
        <div className="flex gap-5 items-center font-semibold">
          <a href="/">Kessler</a>
        </div>
        <HeaderAuth />
      </div>
    </nav>
  );
};
export default Header;
