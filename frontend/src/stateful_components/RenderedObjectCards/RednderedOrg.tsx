import { sleep } from "@/lib/utils";

export default async function RenderedOrg({ org_id }: { org_id: string }) {
  await sleep(1000)
  return <p>IMPLEMENT THE RENDERED CARD</p>
}
