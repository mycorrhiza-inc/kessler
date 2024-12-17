import Link from "next/link";
import ConversationTableInfiniteScroll from "../Organizations/ConversationTable";
import PageContainer from "../Page/PageContainer";
import { ChatModalClickDiv } from "../Chat/ChatModal";
import OrganizationTableInfiniteScroll from "../Organizations/OrganizationTable";
import RecentUpdatesView from "./RecentUpdatesView";

export default function RecentUpdatesPage() {
  return (
    <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
      <div className="grid grid-cols-2 w-full">
        <div>
          <Link className="text-3xl font-bold hover:underline" href="/dockets">
            Dockets
          </Link>
          <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
            <ConversationTableInfiniteScroll />
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
      <ChatModalClickDiv
        className="btn btn-accent w-full"
        inheritedFilters={[]}
      >
        Unsure of what to do? Try chatting with the entire New York PUC
      </ChatModalClickDiv>
      <Link className="btn btn-primary w-full" href="/files">
        Search all Files
      </Link>
      <h1 className=" text-2xl font-bold">Newest Docs</h1>
      <RecentUpdatesView />
    </PageContainer>
  );
}
