import { generateTypeUrlParams } from "@/lib/types/url_params";
import HomePageServer from "@/stateful_components/HomePage/HomePageServer";

export default async function Page(
  { searchParams }: { searchParams: Promise<{ [key: string]: string | string[] | undefined }> }
) {
  const untypedUrlParams = await searchParams;
  const urlParams = generateTypeUrlParams(untypedUrlParams)
  return <HomePageServer />;
}
