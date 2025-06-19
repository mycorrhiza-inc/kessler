import DefaultContainer from "@/components/stateful/PageContainer/DefaultContainer";

export default async function Page({
  params,
  // searchParams
}: {
  params: Promise<{ filling_id: string }>;
  // searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}) {

  // const urlParams = generateTypeUrlParams(await searchParams);
  const filling_id = (await params).filling_id;
  return <DefaultContainer><p>TODO: IMPLEMENT FILING PAGES</p></DefaultContainer>
}
