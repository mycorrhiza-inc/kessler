"use client";
import MarkdownRenderer, {
  testMarkdownContent,
} from "@/components/MarkdownRenderer";
import PageContainer from "@/components/Page/PageContainer";

export default function Page() {
  return (
    <>
      <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
        <MarkdownRenderer>{testMarkdownContent}</MarkdownRenderer>
      </PageContainer>
    </>
  );
}
