import SearchApp from "@/components/SearchApp";
import Landing from "@/components/landing/Landing";

import { createClient } from "@/utils/supabase/server";
import { redirect } from "next/navigation";
export default async function Page() {
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  const isLoggedIn = Boolean(user);
  return (
    <div className="w-full">
      {isLoggedIn ? (
        <SearchApp />
      ) : (
        <Landing />
      )}
    </div>
  );
}
