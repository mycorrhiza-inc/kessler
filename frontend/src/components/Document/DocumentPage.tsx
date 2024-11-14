import DocumentModalBody from "./DocumentModalBody";

const DocumentPage = ({ objectId }: { objectId: string }) => {
  const open = true;
  const title = "Test Document";
  return <DocumentModalBody open={open} objectId={objectId} title={title} />;
};
export default DocumentPage;
