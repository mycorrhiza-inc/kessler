"use client";
import Link from "next/link";
import ConversationTableInfiniteScroll from "../Organizations/ConversationTable";
import PageContainer from "../Page/PageContainer";
import { ChatModalClickDiv } from "../Chat/ChatModal";
import OrganizationTableInfiniteScroll from "../Organizations/OrganizationTable";
import RecentUpdatesView from "./RecentUpdatesView";
import { useKesslerStore } from "@/lib/store";

export default function RecentUpdatesPage() {
  const globalStore = useKesslerStore();
  const experimentalFeatures = globalStore.experimentalFeaturesEnabled;
  return (
    <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
      <div className="grid grid-cols-2 w-full">
        <div>
          <Link className="text-3xl font-bold hover:underline" href="/dockets">
            Dockets
          </Link>
          <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
            <ConversationTableInfiniteScroll truncate />
          </div>
        </div>
        <div>
          <Link className="text-3xl font-bold hover:underline" href="/orgs">
            Organizations
          </Link>
          <div className="max-h-[600px] overflow-x-hidden pl-4">
            <OrganizationTableInfiniteScroll />
          </div>
        </div>
      </div>
      {experimentalFeatures && (
        <ChatModalClickDiv
          className="btn btn-accent w-full"
          inheritedFilters={[]}
        >
          Unsure of what to do? Try chatting with the entire New York PUC
        </ChatModalClickDiv>
      )}
      <h1 className=" text-2xl font-bold">Newest Docs</h1>
      <RecentUpdatesView />
    </PageContainer>
  );
}
