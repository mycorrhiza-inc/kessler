import {
  ConversationPage,
  generateConversationInfo,
} from "@/components/ObjectPages/ConversationPage";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { Metadata } from "next";
import { headers } from "next/headers";
import { cache } from "react";

const cachedConvoInfo = cache(generateConversationInfo);

export async function generateMetadata({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}): Promise<Metadata> {
  const slug = (await params).conversation_id;
  const convoInfo = await cachedConvoInfo(slug);
  return {
    title: convoInfo.displayTitle,
  };
}

export default async function Page({
  params,
}: {
  params: Promise<{ conversation_id: string }>;
}) {
  const slug = (await params).conversation_id;
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const convoInfo = await cachedConvoInfo(slug);
  convoInfo.breadcrumbs.state = state || "";

  return (
    <ConversationPage
      conversation={convoInfo.conversation}
      breadcrumbs={convoInfo.breadcrumbs}
    // inheritedFilters={convoInfo.inheritedFilters}
    />
  );
}
