"use client";
import MarkdownRenderer, {
  testMarkdownContent,
} from "@/components/MarkdownRenderer";
import { ForceUseClient } from "@/components/UtilityComponents";
export const dynamic = "force-dynamic";
export default function Page() {
  return (
    <>
        <ForceUseClient>
          <MarkdownRenderer>{testMarkdownContent}</MarkdownRenderer>
        </ForceUseClient>
    </>
  );
}
