import React from "react";
import Markdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import remarkMath from "remark-math";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dracula } from "react-syntax-highlighter/dist/cjs/styles/prism";
import { useState } from "react";
interface MarkdownRendererProps {
  children: string;
}
const CodeBlock = ({ node, inline, className, children, ...props }: any) => {
  const match = /language-(\w+)/.exec(className || "");
  const codeContent = String(children).replace(/\n$/, "");
  const [buttonText, setButtonText] = useState("Copy");

  const copyToClipboard = () => {
    navigator.clipboard.writeText(codeContent);
    setButtonText("Copied!");
    setTimeout(() => setButtonText("Copy"), 2000);
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
      <button
        onClick={copyToClipboard}
        style={{
          position: "absolute",
          top: "10px",
          right: "10px",
          zIndex: 1,
        }}
      >
        {buttonText}
      </button>
    </div>
  ) : (
    <code className={className} {...props}>
      {children}
    </code>
  );
};

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ children }) => {
  return (
    <>
      <article className="prose">
        <Markdown
          remarkPlugins={[remarkMath, remarkGfm]}
          rehypePlugins={[rehypeKatex, rehypeRaw]}
          components={{
            code: CodeBlock,
          }}
        >
          {children}
        </Markdown>
      </article>
    </>
  );
};

export default MarkdownRenderer;
