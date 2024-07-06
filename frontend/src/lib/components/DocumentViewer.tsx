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
import { pdfjs } from "react-pdf";
import { Document, Page } from "react-pdf";

pdfjs.GlobalWorkerOptions.workerSrc = new URL(
  "pdfjs-dist/build/pdf.worker.min.mjs",
  import.meta.url,
).toString();

const ViewDocumentButton: React.FC<{ document_uuid: string }> = ({
  document_uuid,
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [numPages, setNumPages] = useState<number>();
  const [pageNumber, setPageNumber] = useState<number>(1);

  const toast = useToast();
  const [markdownContent, setMarkdownContent] = useState<string>(
    "# Loading Document Contents",
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

  function onDocumentLoadSuccess({ numPages }: { numPages: number }): void {
    setNumPages(numPages);
  }

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
      <Button onClick={onOpen}>View Document</Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>View Document</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <div>
              <Document
                file="https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf"
                onLoadSuccess={onDocumentLoadSuccess}
              >
                <Page pageNumber={pageNumber} />
              </Document>
              <p>
                Page {pageNumber} of {numPages}
              </p>
            </div>
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
