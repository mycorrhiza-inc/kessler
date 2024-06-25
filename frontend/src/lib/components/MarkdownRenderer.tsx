// TODO : Maybe this is a bad idea since it could really bloat the application
// Math stuff could be thrown out, having tables is important, and code highlighting is kinda standard.
import React from "react";
import Markdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import remarkMath from "remark-math";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dracula } from "react-syntax-highlighter/dist/cjs/styles/prism";

import { Button } from "@chakra-ui/react";
interface MarkdownRendererProps {
  children: string;
}
const CodeBlock = ({ node, inline, className, children, ...props }: any) => {
  const match = /language-(\w+)/.exec(className || "");
  const codeContent = String(children).replace(/\n$/, "");

  const copyToClipboard = () => {
    navigator.clipboard.writeText(codeContent);
    alert("Code copied to clipboard!");
  };

  return !inline && match ? (
    <div style={{ position: "relative" }}>
      <div>
        <SyntaxHighlighter
          style={dracula}
          PreTag="div"
          language={match[1]}
          {...props}
        >
          {codeContent}
        </SyntaxHighlighter>
      </div>
      <Button
        onClick={copyToClipboard}
        style={{
          position: "absolute",
          top: "10px",
          right: "10px",
          zIndex: 1,
        }}
      >
        Copy
      </Button>
    </div>
  ) : (
    <code className={className} {...props}>
      {children}
    </code>
  );
};

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ children }) => {
  return (
    <Markdown
      remarkPlugins={[remarkMath, remarkGfm]}
      rehypePlugins={[rehypeKatex, rehypeRaw]}
      components={{
        code: CodeBlock,
      }}
    >
      {children}
    </Markdown>
  );
};

export default MarkdownRenderer;
