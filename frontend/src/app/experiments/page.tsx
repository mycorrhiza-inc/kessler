import TestFlowVisuals from "@/components/SeriousGaming/SeriousGameAdmin";
import PlanetStartPage from "@/components/experiments/PlanetHomepage";
import { createClient } from "@/utils/supabase/server";
export default async function Page() {
  const supabase = createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();
  return (
    <div className="flex flex-col justify-center space-y-10">
      <TestFlowVisuals user={user}></TestFlowVisuals>
      <PlanetStartPage articles={{}} />
    </div>
  );
}
