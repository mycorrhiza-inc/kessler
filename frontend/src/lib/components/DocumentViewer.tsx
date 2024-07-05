import React, { useState } from "react";
import {
  Button,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useDisclosure,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Input,
  FormControl,
  FormLabel,
  useToast,
} from "@chakra-ui/react";
import { useForm } from "react-hook-form";
import axios from "axios";
import MarkdownRenderer from "./MarkdownRenderer";

const ViewDocumentButton: any = ({
  document_uuid,
}: {
  document_uuid: string;
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();

  const getMarkdownContent = async (document_uuid: string) => {
    try {
      const response = await fetch("/api/files/get_markdown/" + document_uuid, {
        method: "GET",
      });

      return response.text;
    } catch (error) {
      toast({
        title: "No text found.",
        description: `Could not get text from server`,
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    }
  };
  // Fix async stuff
  const markdown_content: string = await getMarkdownContent(document_uuid);

  return (
    <>
      <Button onClick={onOpen}>Upload File</Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>View Document</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <MarkdownRenderer>{markdown_content}</MarkdownRenderer>
          </ModalBody>
          <ModalFooter>
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
export default ViewDocumentButton;
