import { sleep } from "@/lib/utils";

export default async function RenderedConvo({ convo_id }: { convo_id: string }) {
  await sleep(1000)
  return <p>IMPLEMENT THE RENDERED CARD</p>
}
