import Image from "next/image";

import PlanetStartPage from "@/components/PlanetHomepage";
import { exampleArticle } from "@/interfaces";
export default function Home() {
  return <PlanetStartPage articles={[exampleArticle]}></PlanetStartPage>;
}
