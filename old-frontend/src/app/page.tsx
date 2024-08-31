import DefaultShell from "../lib/components/DefaultShell";
import DashboardPrompt from "../lib/components/DashboardPrompt";
import ChatUI from "../lib/components/ChatUI";

export default function Page() {
  return (
    <DefaultShell>
      {/* <DashboardPrompt /> */}
      <ChatUI
        chatUrl="/api/rag/rag_chat"
        modelOptions={[]}
        useCitations={true}
      />
    </DefaultShell>
  );
}
// export { getServerSideProps } from "../lib/Chakra";
