import { useState } from "react";
import DocumentModalBody from "./DocumentModalBody";
import { BreadcrumbValues } from "../SitemapUtils";
import { User } from "@supabase/supabase-js";
import Navbar from "../Navbar";

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
    <>
      <Navbar user={user} breadcrumbs={breadcrumbs} />
      <div className="w-full h-full">
        <DocumentModalBody
          open={open}
          objectId={objectId}
          setTitle={setTitle}
        />
      </div>
    </>
  );
};
export default DocumentPage;
