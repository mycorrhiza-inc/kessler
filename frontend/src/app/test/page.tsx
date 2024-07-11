"use client"
import React from "react";
import { Box, Text, VStack, HStack, Avatar } from "@chakra-ui/react";
import MarkdownRenderer from "../../lib/components/MarkdownRenderer";

import PDFTesting from "../../lib/components/PDFViewer"
import NoSSR from "../../lib/components/NoSSR";
import dynamic from 'next/dynamic'

const Page = () => 
 <>
    <Text>This is a test page</Text>
    <NoSSR>
      <PDFTesting></PDFTesting>
    </NoSSR>
 </>
  ;
// <PDFTesting></PDFTesting>


export default dynamic(() => Promise.resolve(Page), {
  ssr: false
})
