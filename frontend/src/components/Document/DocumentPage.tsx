"use client";
import { useState } from "react";
import DocumentModalBody from "./DocumentModalBody";
import { BreadcrumbValues } from "../SitemapUtils";
import { User } from "@supabase/supabase-js";
import PageContainer from "../Page/PageContainer";

const DocumentPage = ({
  objectId,
  state,
  user,
}: {
  objectId: string;
  state?: string;
  user: User | null;
}) => {
  const open = true;
  const [title, setTitle] = useState("Loading Document...");
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [
      { title: "Files", value: "files" },
      { title: title, value: objectId },
    ],
  };
  return (
    <PageContainer user={user} breadcrumbs={breadcrumbs}>
      <DocumentModalBody
        open={open}
        objectId={objectId}
        setTitle={setTitle}
        isPage={true}
      />
    </PageContainer>
  );
};
export default DocumentPage;
