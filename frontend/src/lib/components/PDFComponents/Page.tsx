import React from "react";

interface PageProps {
  children: React.ReactNode;
  style?: React.CSSProperties;
}

const Page = React.memo((props: PageProps) => {
  const { children, style } = props;
  const internalStyle: React.CSSProperties = {
    ...style,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    outline: "1px solid #ccc"
  };
  return <div style={internalStyle}>{children}</div>;
});

export default Page;


