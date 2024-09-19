import FetchDataSteps from "@/components/supabasetutorial/tutorial/fetch-data-steps";
import { createClient } from "@/utils/supabase/server";
import { InfoIcon } from "lucide-react";
import { redirect } from "next/navigation";

import PlanetStartPage from "@/components/experiments/PlanetHomepage";
import { exampleArticle } from "@/utils/interfaces";
import SearchApp from "@/components/SearchApp";
export default async function ProtectedPage() {
  const supabase = createClient();

  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (user) {
    return redirect("/");
  }

  return (
    <div className="w-full">
      <SearchApp></SearchApp>
    </div>
  );
}
