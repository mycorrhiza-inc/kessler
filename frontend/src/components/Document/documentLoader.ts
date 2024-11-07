import axios from "axios";

type DocumentData = {
  text: string;
  metadata: any;
  pdfUrl: string;
};

function wrapPromise(promise: Promise<any>) {
  let status = "pending";
  let result: any;
  let suspender = promise.then(
    (r) => {
      status = "success";
      result = r;
    },
    (e) => {
      status = "error";
      result = e;
    },
  );

  return {
    read() {
      if (status === "pending") {
        throw suspender;
      } else if (status === "error") {
        throw result;
      }
      return result;
    },
  };
}

export function fetchDocumentData(objectId: string, overridePDFUrl?: string) {
  const documentPromise = Promise.all([
    axios.get(`https://api.kessler.xyz/v2/public/files/${objectId}/markdown`),
    axios.get(`https://api.kessler.xyz/v2/public/files/${objectId}`),
  ]).then(([textResponse, metadataResponse]) => ({
    text: textResponse.data,
    metadata: metadataResponse.data,
    pdfUrl: overridePDFUrl || `/api/v1/files/${objectId}/raw`,
  }));

  return wrapPromise(documentPromise);
}
