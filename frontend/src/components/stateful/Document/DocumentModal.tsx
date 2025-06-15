"use client";
import { completeFileSchemaGet } from "@/lib/requests/search";
import { DocumentMainTabs } from "./DocumentBody";
import { CompleteFileSchema } from "@/lib/types/backend_schemas";

import { CLIENT_API_URL } from "@/lib/env_variables";
import LoadingSpinner from "@/componenets/style/misc/LoadingSpinner";
import Modal from "@/componenets/style/misc/Modal";
type ModalProps = {
  objectId: string;
  children?: React.ReactNode;
  isPage: boolean;
};
const DocumentModalBody = ({ objectId, isPage }: ModalProps) => {
  const semiCompleteFileUrl = `${CLIENT_API_URL}/v2/public/files/${objectId}`;
  // const { data, error, isLoading } = useSWRImmutable(
  //   semiCompleteFileUrl,
  //   completeFileSchemaGet,
  // );
  // if (isLoading) {
  //   return <LoadingSpinner loadingText="Loading Document" />;
  // }
  // if (error) {
  //   return (
  //     <p>Encountered an error getting text from the server: {String(error)}</p>
  //   );
  // }
  // const docObj = data as CompleteFileSchema;
  // return <DocumentMainTabs documentObject={docObj} isPage={isPage} />;
  return <p>NOT IMPLEMENTED</p>
};

const DocumentModal = ({
  objectId,
  open,
  setOpen,
}: {
  objectId: string;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
}) => {
  return (
    <Modal open={open} setOpen={setOpen}>
      <DocumentModalBody objectId={objectId} isPage={false} />
    </Modal>
  );
};

export default DocumentModal;
