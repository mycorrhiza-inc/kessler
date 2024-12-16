import { BreadcrumbValues } from "../SitemapUtils";
import PageContainer from "../Page/PageContainer";
import { completeFileSchemaGet } from "@/lib/requests/search";
import { DocumentMainTabs } from "./DocumentBody";
import { internalAPIURL } from "@/lib/env_variables";
import { CompleteFileSchema } from "@/lib/types/backend_schemas";

export const generateDocumentPageData = async (objectId: string) => {
  const semiCompleteFileUrl = `${internalAPIURL}/v2/public/files/${objectId}`;
  const fileObj = await completeFileSchemaGet(semiCompleteFileUrl);
  const breadcrumbs: BreadcrumbValues = {
    breadcrumbs: [
      { title: "Files", value: "files" },
      { title: fileObj.name, value: objectId },
    ],
  };
  return { breadcrumbs: breadcrumbs, fileObj: fileObj };
};

const DocumentPage = ({
  breadcrumbs,
  fileObj,
}: {
  breadcrumbs: BreadcrumbValues;
  fileObj: CompleteFileSchema;
}) => {
  // PROD API URL, since substituting in localhost doesnt work if you try to fetch from within a docker container
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <DocumentMainTabs documentObject={fileObj} isPage={true} />
    </PageContainer>
  );
};
export default DocumentPage;
