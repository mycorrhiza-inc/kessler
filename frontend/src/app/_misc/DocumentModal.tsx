"use client";
import { apiURL } from "@/lib/env_variables";
import { completeFileSchemaGet } from "@/lib/requests/search";
import useSWRImmutable from "swr";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import { DocumentMainTabs } from "@/components/Document/DocumentBody";
import { CompleteFileSchema } from "@/lib/types/backend_schemas";

import Modal from "@/components/styled-components/Modal";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
type ModalProps = {
  objectId: string;
  children?: React.ReactNode;
  isPage: boolean;
};
const DocumentModalBody = ({ objectId, isPage }: ModalProps) => {
  const semiCompleteFileUrl = `${apiURL}/v2/public/files/${objectId}`;
  const { data, error, isLoading } = useSWRImmutable(
    semiCompleteFileUrl,
    completeFileSchemaGet,
  );
  if (isLoading) {
    return <LoadingSpinner loadingText="Loading Document" />;
  }
  if (error) {
    return (
      <p>Encountered an error getting text from the server: {String(error)}</p>
    );
  }
  const docObj = data as CompleteFileSchema;
  return <DocumentMainTabs documentObject={docObj} isPage={isPage} />;
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
  const nextJSRouter = useRouter();
  const [previousOpenState, setPreviousOpenState] = useState(open);
  useEffect(() => {
    if (previousOpenState != open) {
      if (open) {
        nextJSRouter.push(`/file/${objectId}`);
      } else {
        nextJSRouter.back();
      }
    }
    setPreviousOpenState(open);
  }, [open]);
  return (
    <Modal open={open} setOpen={setOpen}>
      <DocumentModalBody objectId={objectId} isPage={false} />
    </Modal>
  );
};

export default DocumentModal;
