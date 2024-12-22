"use client";
import MarkdownRenderer, {
  testMarkdownContent,
} from "@/components/MarkdownRenderer";
import PageContainer from "@/components/Page/PageContainer";
import { ForceUseClient } from "@/components/UtilityComponents";
export const dynamic = "force-dynamic";
export default function Page() {
  return (
    <>
      <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
        <ForceUseClient>
          <MarkdownRenderer>{testMarkdownContent}</MarkdownRenderer>
        </ForceUseClient>
      </PageContainer>
    </>
  );
}
