import { signOutAction } from "@/app/actions";
import { ThemeSwitcher } from "@/components/supabasetutorial/theme-switcher";
import { Button } from "@/components/supabasetutorial/ui/button";
import { createClient } from "@/utils/supabase/server";

async function HeaderAuth() {
  const {
    data: { user },
  } = await createClient().auth.getUser();

  return user ? (
    <div className="flex items-center gap-4">
      Hey, {user.email}!
      <form action={signOutAction}>
        <Button type="submit" variant={"outline"}>
          Sign out
        </Button>
      </form>
    </div>
  ) : (
    <div className="flex gap-2">
      <Button asChild size="sm" variant={"outline"}>
        <a href="/sign-in">Sign in</a>
      </Button>
      <Button asChild size="sm" variant={"default"}>
        <a href="/sign-up">Sign up</a>
      </Button>
    </div>
  );
}
const Header = () => {
  return (
    <nav
      className="w-full flex justify-center border-b border-b-foreground/10 h-16 bg-white dark:bg-black"
      // style={{ position: "fixed", top: 0 }}
    >
      <div
        className="w-full max-w-5xl flex justify-between items-center p-3 px-5 text-sm"
        style={{ zIndex: 3000 }}
      >
        <div className="flex gap-5 items-center font-semibold">
          <a href="/">Kessler</a>
        </div>
        <HeaderAuth />
        <ThemeSwitcher />
      </div>
    </nav>
  );
};
export default Header;
