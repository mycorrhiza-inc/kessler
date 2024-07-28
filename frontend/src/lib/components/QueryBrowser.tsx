import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
} from "@chakra-ui/react";
import { useState } from "react";
import FilePageBrowser from "./FilePageBrowser";
import UploadFileButton from "./FileUploader";

const QueryBrowser: React.FC = () => {
  const [formData, setFormData] = useState({
    match_name: "",
    match_source: "",
    match_doctype: "",
    match_stage: "",
    match_docket_id: "",
    match_document_class: "",
    match_author: "",
  });
  const [queryData, setQueryData] = useState<any | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const data = {
      match_name: formData.match_name || null,
      match_source: formData.match_source || null,
      match_doctype: formData.match_doctype || null,
      match_stage: formData.match_stage || null,
      match_mdata: {
        docket_id: formData.match_docket_id || null,
        document_class: formData.match_document_class || null,
        author: formData.match_author || null,
      },
    };
    setQueryData(data);
  };

  return (
    <Box>
      <UploadFileButton></UploadFileButton>
      <form onSubmit={handleSubmit}>
        <VStack spacing={4}>
          <FormControl>
            <FormLabel>Name</FormLabel>
            <Input
              name="match_name"
              value={formData.match_name}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Source</FormLabel>
            <Input
              name="match_source"
              value={formData.match_source}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Author</FormLabel>
            <Input
              name="match_author"
              value={formData.match_source}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Docket ID</FormLabel>
            <Input
              name="match_docket_id"
              value={formData.match_source}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Document Class (Public Comment/Decision)</FormLabel>
            <Input
              name="match_source"
              value={formData.match_document_class}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>File Extension</FormLabel>
            <Input
              name="match_doctype"
              value={formData.match_doctype}
              onChange={handleChange}
            />
          </FormControl>
          <FormControl>
            <FormLabel>Procesing Stage</FormLabel>
            <Input
              name="match_stage"
              value={formData.match_stage}
              onChange={handleChange}
            />
          </FormControl>
          <Button type="submit">Search</Button>
        </VStack>
      </form>
      {queryData && (
        <FilePageBrowser fileUrl="/api/files/query/paginate" data={queryData} />
      )}
    </Box>
  );
};

export default QueryBrowser;
