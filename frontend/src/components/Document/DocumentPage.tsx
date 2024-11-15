import DocumentModalBody from "./DocumentModalBody";

const DocumentPage = ({ objectId }: { objectId: string }) => {
  const open = true;
  const title = "Test Document";
  return (
    <div className="w-full h-full">
      <DocumentModalBody open={open} objectId={objectId} title={title} />
    </div>
  );
};
export default DocumentPage;
