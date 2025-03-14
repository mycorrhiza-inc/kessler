import Link from "next/link";
import ConversationTableInfiniteScroll from "../LookupPages/ConvoLookup/ConversationTable";
import { ExperimentalChatModalClickDiv } from "../Chat/ChatModal";
import OrganizationTableInfiniteScroll from "../LookupPages/OrgLookup/OrganizationTable";
import FileSearchView from "../Search/FileSearch/FileSearchView";

export default function HomePage() {
  return (
    <>
      <div>
        <div className="grid grid-cols-2 w-full z-1">
          <div>
            <Link
              className="text-3xl font-bold hover:underline"
              href="/dockets"
            >
              Dockets
            </Link>
            <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
              <ConversationTableInfiniteScroll
                truncate
                lookup_data={{ query: "" }}
              />
            </div>
          </div>
          <div className="z-[1]">
            <Link
              className="text-3xl font-bold hover:underline mb-5 p-10"
              href="/orgs"
            >
              Organizations
            </Link>
            <div className="max-h-[600px] overflow-x-hidden pl-4">
              <OrganizationTableInfiniteScroll />
            </div>
          </div>
        </div>
        <ExperimentalChatModalClickDiv
          className="btn btn-accent w-full"
          inheritedFilters={[]}
        >
          Unsure of what to do? Try chatting with the entire New York PUC
        </ExperimentalChatModalClickDiv>

        <h1 className=" text-2xl font-bold">Newest Docs</h1>
        <FileSearchView inheritedFilters={[]} />
      </div>
    </>
  );
}
