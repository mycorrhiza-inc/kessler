import {
  ConversationPage,
  generateConversationInfo,
} from "@/components/Conversations/ConversationPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { createClient } from "@/utils/supabase/server";
import { Metadata } from "next";
import { headers } from "next/headers";

export const metadata: Metadata = {
  title: "ERROR IN SITE NAME",
};
export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const slug = (await params).conversation_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const convoInfo = await generateConversationInfo(slug, state || "");
  metadata.title = convoInfo.displayTitle;

  return (
    <ConversationPage
      conversation={convoInfo.conversation}
      breadcrumbs={convoInfo.breadcrumbs}
      inheritedFilters={convoInfo.inheritedFilters}
    />
  );
}
