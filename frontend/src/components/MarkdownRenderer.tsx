import React from "react";
import Markdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import remarkMath from "remark-math";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dracula } from "react-syntax-highlighter/dist/cjs/styles/prism";
import rehypeComponents from "rehype-components";
import { LinkDocket, LinkFile } from "./Chat/LLMComponents";
import { subdividedHueFromSeed } from "./Tables/TextPills";
interface MarkdownRendererProps {
  children: string;
  color?: string;
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
  name: "John",
  "age": 30,
  "city": "New York"
}
\`\`\`

\`\`\`python
# Python code block
def greet():
    print("Hello, world!")
\`\`\`

# Custom Components

In order to access the docket, <link-docket text="click here" docket_id="18-M-0084"/>. 

The organization <link-organization text="Public Service Comission" name="Public Service Comission"/> created the document.

Their report <link-file text="1" uuid="777b5c2d-d19e-4711-b2ed-2ba9bcfe449a" /> claims xcel energy failed to meet its renewable energy targets.
`;
const CodeBlock = ({ node, inline, className, children, ...props }: any) => {
  const match = /language-(\w+)/.exec(className || "");
  const codeContent = String(children).replace(/\n$/, "");

  const copyToClipboard = () => {
    navigator.clipboard.writeText(codeContent);
  };

  return inline! && match ? (
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
      <label
        className="swap"
        style={{ position: "absolute", top: "10px", right: "10px", zIndex: 1 }}
      >
        <input type="checkbox" onClick={copyToClipboard} />
        <div className="swap-on">Copied!</div>
        <div className="swap-off">Copy</div>
      </label>
    </div>
  ) : (
    <code className={className} {...props}>
      {children}
    </code>
  );
};
// Im sorry - nic
// Make a component that will take in tags like this from a markdown string
//
// \`\`\`python
// # Python code block
// def greet():
//     print("Hello, world!")
// \`\`\`
//
// # Custom Components
//
// In order to access the docket, <link-docket text="click here" docket_id="18-M-0084"/>.
//
// The organization <link-organization text="Public Service Comission" uuid="Public Service Comission"/> created the document.
//
// Their report <link-file text="1" uuid="777b5c2d-d19e-4711-b2ed-2ba9bcfe449a" /> claims xcel energy failed to meet its renewable energy targets.
//
// and change every instance of <link-docket text="click here" docket_id="18-M-0084"/> to a component
// <a href="/docket/18-M-0084" className
// style={{ backgroundColor: buttonColor }}
// className="btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty	" />
const horrificMarkdownComponentMangle = (inputMarkdown: string): string => {
  let mangledMarkdown = inputMarkdown;

  // Helper function to extract attributes from tag string
  const getAttributes = (tagContent: string) => {
    const attrs: Record<string, string> = {};
    const attrRegex = /(\w+)="([^"]+)"/g;
    let match;
    while ((match = attrRegex.exec(tagContent)) !== null) {
      attrs[match[1]] = match[2];
    }
    return attrs;
  };

  // Replace link-docket tags
  mangledMarkdown = mangledMarkdown.replace(
    /<link-docket([^/>]+)\/>/g,
    (match, attributes) => {
      const attrs = getAttributes(attributes);
      const color = subdividedHueFromSeed(attrs.docket_id);
      return `<a  target="_blank" href="/dockets/${attrs.docket_id}" class="btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty" style="background-color: ${color}"> ${attrs.text} </a>`;
    },
  );

  // Replace link-file tags
  mangledMarkdown = mangledMarkdown.replace(
    /<link-file([^/>]+)\/>/g,
    (match, attributes) => {
      const attrs = getAttributes(attributes);
      const color = subdividedHueFromSeed(attrs.uuid);
      return `<a target="_blank" href="/files/${attrs.uuid}" class="btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty" style="background-color: ${color}"> ${attrs.text} </a>`;
    },
  );

  // Replace link-organization tags
  mangledMarkdown = mangledMarkdown.replace(
    /<link-organization([^/>]+)\/>/g,
    (match, attributes) => {
      const attrs = getAttributes(attributes);
      const color = subdividedHueFromSeed(attrs.name);
      return `<a target="_blank" href="/orgs/${attrs.uuid}" class="btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty" style="background-color: ${color}"> ${attrs.text} </a>`;
    },
  );

  return mangledMarkdown;
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
          rehypePlugins={[
            rehypeRaw,
            [
              rehypeComponents,
              {
                components: {
                  // Uncomenting these with rehypeRaw causes the bug: TypeError: cyclic object value
                  // Uncomennting w/o rehypeRaw reverts to the default html escaping behavior.
                  // "link-file": LinkFile,
                  // "link-docket": LinkDocket,
                  // When trying with server side rendering I get a more detailed error:
                  //
                  // 67.36    Generating static pages (12/17)
                  // 67.63 Error: rehype-components: Component function is expected to return ElementContent or an array of ElementContent, but got [{"key":null,"ref":null,"props":{"text":"click here"},"_owner":null}].
                  // 67.63     at /app/.next/server/chunks/223.js:1:1069893
                  // 67.63     at /app/.next/server/chunks/223.js:1:1393842
                  // 67.63     at node (element<link-docket>) (/app/.next/server/chunks/223.js:1:1393285)
                  // 67.63     at node (element<p>) (/app/.next/server/chunks/223.js:1:1393505)
                  // 67.63     at node (root) (/app/.next/server/chunks/223.js:1:1393505)
                  // 67.63     at s (/app/.next/server/chunks/223.js:1:1393582)
                  // 67.63     at a (/app/.next/server/chunks/223.js:1:1393763)
                  // 67.63     at /app/.next/server/chunks/223.js:1:1069722
                  // 67.63     at /app/.next/server/chunks/223.js:1:1060981
                  // 67.63     at a (/app/.next/server/chunks/223.js:1:1061186) {
                  // 67.63   digest: '3998901393'
                  // 67.63 }
                },
              },
            ],
            rehypeKatex,
          ]}
          components={{
            code: CodeBlock,
          }}
        >
          {horrificMarkdownComponentMangle(children as string)}
        </Markdown>
      </article>
    </>
  );
};

export default MarkdownRenderer;
