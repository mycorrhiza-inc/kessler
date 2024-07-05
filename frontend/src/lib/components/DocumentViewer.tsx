import React, { useState, useEffect } from "react";
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
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import MarkdownRenderer from "./MarkdownRenderer";

const ViewDocumentButton: React.FC<{ document_uuid: string }> = ({
  document_uuid,
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const toast = useToast();
  const [markdownContent, setMarkdownContent] = useState<string>(
    "Loading Document Contents",
  );

  const getMarkdownContent = async (document_uuid: string) => {
    try {
      const response = await fetch("/api/files/get_markdown/" + document_uuid, {
        method: "GET",
      });

      return response.text();
    } catch (error) {
      toast({
        title: "No text found.",
        description: "Could not get text from server",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      return "Error loading document";
    }
  };

  useEffect(() => {
    if (isOpen) {
      (async () => {
        const content = await getMarkdownContent(document_uuid);
        setMarkdownContent(content);
      })();
    }
  }, [isOpen, document_uuid]);

  return (
    <>
      <Button onClick={onOpen}>Upload File</Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>View Document</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <MarkdownRenderer>{markdownContent}</MarkdownRenderer>
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
