import React, { useState } from "react";
import {
  Box,
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

const UploadFileButton: React.FC = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [selectedFiles, setSelectedFiles] = useState<FileList | null>(null);
  const { register, handleSubmit, reset } = useForm();
  const toast = useToast();

  const onFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const onFolderChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFiles(event.target.files);
    }
  };

  const uploadFile = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    try {
      await axios.post("/api/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });
      toast({
        title: "File uploaded",
        description: `Your file ${file.name} has been uploaded successfully.`,
        status: "success",
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: "Upload failed",
        description: `There was an error uploading your file ${file.name}.`,
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    }
  };
  const onSubmit = async () => {
    if (!selectedFile && !selectedFiles) {
      toast({
        title: "No file or folder selected",
        description: "Please select a file or folder to upload.",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
      return;
    }

    if (selectedFile) {
      await uploadFile(selectedFile);
    }

    if (selectedFiles) {
      for (let i = 0; i < selectedFiles.length; i++) {
        await uploadFile(selectedFiles[i]);
      }
    }

    reset();
    setSelectedFile(null);
    setSelectedFiles(null);
    onClose();
  };

  const onUrlSubmit = async (data: any) => {
    try {
      const response = await fetch("/api/files/add_url?process=true", {
        method: "POST",
        headers: {
          accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          url: data.url,
          metadata: {
            source: "Personal",
          },
        }),
      });

      if (!response.ok) {
        throw new Error("Network response was not ok");
      }

      toast({
        title: "URL processed",
        description: `The URL ${data.url} has been processed.`,
        status: "success",
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: "Error processing URL",
        description: `There was an error processing the URL`,
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    } finally {
      reset();
      onClose();
    }
  };

  return (
    <>
      <Button onClick={onOpen}>Upload File</Button>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Upload File</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Tabs>
              <TabList>
                <Tab>File/Folder Upload</Tab>
                <Tab>Process URL</Tab>
              </TabList>
              <TabPanels>
                <TabPanel>
                  <form onSubmit={handleSubmit(onSubmit)}>
                    <FormControl>
                      <FormLabel>File</FormLabel>
                      <Input
                        type="file"
                        accept="*"
                        {...register("file")}
                        onChange={(e) =>
                          setSelectedFile(e.target.files?.[0] || null)
                        }
                      />
                    </FormControl>
                    <Button type="submit">Submit</Button>
                  </form>
                </TabPanel>
                <TabPanel>
                  <form onSubmit={handleSubmit(onUrlSubmit)}>
                    <FormControl>
                      <FormLabel>URL</FormLabel>
                      <Input type="text" {...register("url")} />
                    </FormControl>
                    <Button type="submit">Process URL</Button>
                  </form>
                </TabPanel>
              </TabPanels>
            </Tabs>
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
export default UploadFileButton;
