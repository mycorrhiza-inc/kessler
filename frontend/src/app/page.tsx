import SearchApp from "@/components/SearchApp";
import Landing from "@/components/landing/Landing";

import { createClient } from "@/utils/supabase/server";
import { redirect } from "next/navigation";

import { useKesslerStore } from "@/lib/store";
import AuthGuard from "@/components/AuthGuard";
import RecentUpdatesView from "@/components/RecentUpdates/RecentUpdatesView";

export default async function Page() {
  const checkLoggedIn = async () => {
    const supabase = createClient();
    const {
      data: { user },
    } = await supabase.auth.getUser();
    const userPresent = Boolean(user);
    return userPresent;
  };
  const isLoggedIn = await checkLoggedIn();
  return (
    <div className="w-full">
      {isLoggedIn ? (
        <AuthGuard isLoggedIn={isLoggedIn}>
          <RecentUpdatesView />
        </AuthGuard>
      ) : (
        <Landing />
      )}
    </div>
  );
}
