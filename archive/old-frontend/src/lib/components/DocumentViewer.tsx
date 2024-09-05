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
import MarkdownRenderer from "./MarkdownRenderer";
import PDFViewer from "./PDFViewer";
import { FileType } from "../interfaces";

import { Table, Thead, Tbody, Tr, Th, Td } from "@chakra-ui/react";
const DynamicModal: React.FC<{
  document_object: FileType;
  onClose: () => void;
}> = ({ document_object, onClose }) => {
  const document_uuid = document_object.id;
  const {
    isOpen,
    onOpen,
    onClose: closeModal,
  } = useDisclosure({ defaultIsOpen: true });
  const isPDF = document_object.doctype == "pdf";
  const [numPages, setNumPages] = useState<number>();
  const [pageNumber, setPageNumber] = useState<number>(1);

  const toast = useToast();
  const [markdownContent, setMarkdownContent] = useState<string>(
    "# Loading Document Contents",
  );
  const [pdfURL, setPdfURL] = useState<string>("");

  const getPDFURL = (document_uuid: string) =>
    "/api/files/raw/" + document_uuid;

  const getMarkdownContent = async (document_uuid: string) => {
    try {
      const response = await fetch("/api/files/markdown/" + document_uuid, {
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

  const MetadataTable = ({ mdata }: { mdata: Object }) => {
    return (
      <Table variant="simple" border="1px" borderColor="gray.200">
        <Thead>
          <Tr>
            <Th border="1px" borderColor="gray.200">
              Key
            </Th>
            <Th border="1px" borderColor="gray.200">
              Value
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {Object.entries(mdata).map(([key, value]) => (
            <Tr key={key}>
              <Td border="1px" borderColor="gray.200">
                {key}
              </Td>
              <Td border="1px" borderColor="gray.200">
                {value}
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    );
  };

  function onDocumentLoadSuccess({ numPages }: { numPages: number }): void {
    setNumPages(numPages);
  }

  useEffect(() => {
    if (isOpen) {
      (async () => {
        setPdfURL(getPDFURL(document_uuid));
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
        <ModalHeader>Viewing Document: {document_object.name}</ModalHeader>
        <ModalCloseButton />
        <MetadataTable mdata={document_object.mdata}></MetadataTable>
        <ModalBody>
          <Grid templateColumns="repeat(2, 1fr)" gap={6}>
            {isPDF && (
              <GridItem>
                <a href={pdfURL} target="_blank">
                  <Button>Download PDF</Button>
                </a>
                <PDFViewer file={pdfURL}></PDFViewer>
              </GridItem>
            )}
            <GridItem>
              <a href={"/api/files/markdown/" + document_uuid} target="_blank">
                <Button>Get Raw Text</Button>
              </a>
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

const ViewDocumentButton: React.FC<{ document_object: FileType }> = ({
  document_object,
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
          document_object={document_object}
          onClose={handleCloseModal}
        />
      )}
    </>
  );
};

// PDF Viewer Coming soon
// <Document
//   file="https://raw.githubusercontent.com/mycorrhizainc/examples/main/CO%20Clean%20Energy%20Plan%20Info%20Sheet.pdf"
//   onLoadSuccess={onDocumentLoadSuccess}
// >
//   <Page pageNumber={pageNumber} />
// </Document>
export default ViewDocumentButton;
