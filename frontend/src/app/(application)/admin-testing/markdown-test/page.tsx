"use client";
import MarkdownRenderer, {
  testMarkdownContent,
} from "@/components/MarkdownRenderer";
import { CommandKSearch } from "@/components/Search/CommandK";
import { ForceUseClient } from "@/components/UtilityComponents";
export const dynamic = "force-dynamic";
export default function Page() {
  return (
    <>
      <MarkdownRenderer>{testMarkdownContent}</MarkdownRenderer>
      <CommandKSearch />
    </>
  );
}
