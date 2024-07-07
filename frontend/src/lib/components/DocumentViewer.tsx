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
  Grid,
  GridItem,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import ReactDOM from "react-dom";
import axios from "axios";
import MarkdownRenderer from "./MarkdownRenderer";
import { pdfjs } from "react-pdf";
import { Document, Page } from "react-pdf";

pdfjs.GlobalWorkerOptions.workerSrc = new URL(
  "pdfjs-dist/build/pdf.worker.min.mjs",
  import.meta.url,
).toString();

const DynamicModal: React.FC<{
  document_uuid: string;
  onClose: () => void;
}> = ({ document_uuid, onClose }) => {
  const {
    isOpen,
    onOpen,
    onClose: closeModal,
  } = useDisclosure({ defaultIsOpen: true });
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

  const handleClose = () => {
    closeModal();
    onClose();
  };

  return ReactDOM.createPortal(
    <Modal isOpen={isOpen} onClose={handleClose} size="6xl">
      <ModalOverlay />
      <ModalContent maxW="80%">
        <ModalHeader>View Document</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Grid templateColumns="repeat(2, 1fr)" gap={6}>
            <GridItem>
              <Document
                file="https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf"
                onLoadSuccess={onDocumentLoadSuccess}
              >
                <Page pageNumber={pageNumber} />
              </Document>
            </GridItem>
            <GridItem>
              <MarkdownRenderer>{markdownContent}</MarkdownRenderer>
            </GridItem>
          </Grid>
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" onClick={handleClose}>
            Close
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>,
    document.body,
  );
};

const ViewDocumentButton: React.FC<{ document_uuid: string }> = ({
  document_uuid,
}) => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <Button onClick={handleOpenModal}>View Document</Button>
      {isModalOpen && (
        <DynamicModal
          document_uuid={document_uuid}
          onClose={handleCloseModal}
        />
      )}
    </>
  );
};

export default ViewDocumentButton;
