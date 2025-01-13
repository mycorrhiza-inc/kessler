"use client";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import PageContainer from "@/components/Page/PageContainer";
import { ForceUseClient } from "@/components/UtilityComponents";
export const dynamic = "force-dynamic";
export default function Page() {
  const value = process.env.INTERNAL_KESSLER_API_URL;
  const markdown_string = `# Environment Variables
INTERNAL_API_URL: ${value}

PUBLIC_API_URL: ${process.env.PUBLIC_API_URL}
`;

  return (
    <>
      <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
        <ForceUseClient>
          <MarkdownRenderer>{markdown_string}</MarkdownRenderer>
        </ForceUseClient>
      </PageContainer>
    </>
  );
}
