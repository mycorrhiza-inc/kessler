import { BreadcrumbValues } from "../SitemapUtils";
import { User } from "@supabase/supabase-js";
import PageContainer from "../Page/PageContainer";
import { completeFileSchemaGet } from "@/lib/requests/search";
import { prodAPIURL } from "@/lib/env_variables";
import { CompleteFileSchema } from "@/lib/types/backend_schemas";
import { DocumentMainTabs } from "./DocumentModalBody";

const DocumentPage = async ({
  objectId,
  state,
}: {
  objectId: string;
  state?: string;
  user: User | null;
}) => {
  const semiCompleteFileUrl = `${prodAPIURL}/v2/public/files/${objectId}`;
  const fileObj = await completeFileSchemaGet(semiCompleteFileUrl);
  const breadcrumbs: BreadcrumbValues = {
    state: state,
    breadcrumbs: [
      { title: "Files", value: "files" },
      { title: fileObj.name, value: objectId },
    ],
  };
  // PROD API URL, since substituting in localhost doesnt work if you try to fetch from within a docker container
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <DocumentMainTabs documentObject={fileObj} isPage={true} />
    </PageContainer>
  );
};
export default DocumentPage;
