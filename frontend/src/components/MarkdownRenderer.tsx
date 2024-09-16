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
  color: string | undefined;
}
export const testMarkdownContent = `
# Header 1
## Header 2
### Header 3
#### Header 4
##### Header 5
###### Header 6

**Bold Text**

*Italic Text*

***Bold and Italic Text***

~~Strikethrough Text~~

> Blockquote

* Unordered list item 1
* Unordered list item 2
  * Nested list item
    * Deeper nested list item

1. Ordered list item 1
2. Ordered list item 2
   1. Nested ordered list item
   2. Nested ordered list item

\`Inline code\`


#  \`Inline code\` in a header

\`\`\`javascript
// Code block
function helloWorld() {
  console.log("Hello, world!");
}
\`\`\`

[Link to Google](https://www.google.com)

![Image Alt Text](https://via.placeholder.com/150)

| Table Header 1 | Table Header 2 |
| -------------- | -------------- |
| Table Cell 1   | Table Cell 2   |
| Table Cell 3   | Table Cell 4   |

\`\`\`markdown
# Markdown code block
\`\`\`

- [ ] Task list item 1
- [x] Task list item 2

\`\`\`math
E = mc^2
\`\`\`

\`\`\`json
{
  "name": "John",
  "age": 30,
  "city": "New York"
}
\`\`\`

\`\`\`python
# Python code block
def greet():
    print("Hello, world!")
\`\`\`
`;
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

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
  children,
  color,
}) => {
  const textColor = color
    ? `text-${color} prose-headings:text-${color}`
    : "text-base-content prose-headings:text-base-content";
  return (
    <>
      <article
        // These colors should be style and component specific
        className={"prose prose-neutral " + textColor}
        style={{ maxWidth: "70vw" }}
      >
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
