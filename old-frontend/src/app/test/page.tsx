"use client";
import React from "react";
import { Box, Text, VStack, HStack, Avatar } from "@chakra-ui/react";
import MarkdownRenderer from "../../lib/components/MarkdownRenderer";

import PDFTesting from "../../lib/components/PDFViewer";
import NoSSR from "../../lib/components/NoSSR";
import dynamic from "next/dynamic";

const Page = () => (
  <>
    <Text>This is a test page</Text>
    <NoSSR>
      <PDFTesting file="https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf"></PDFTesting>
    </NoSSR>
  </>
);
// <PDFTesting></PDFTesting>

export default dynamic(() => Promise.resolve(Page), {
  ssr: false,
});
