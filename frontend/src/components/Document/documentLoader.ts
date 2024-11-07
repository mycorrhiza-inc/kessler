import axios from "axios";

type DocumentData = {
  text: string;
  metadata: any;
  pdfUrl: string;
};

export const fetchTextDataFromURL = async (url: string) => {
  const response = await axios.get(url);
  if (response.status !== 200) {
    console.log("Error fetching data");
    console.log(response);
    return "Error Retriving Data";
  }
  if (typeof response.data !== "string") {
    console.log("Did not return string");
    console.log(response.data);
    return "Endpoint did not return text";
  }
  return String(response.data);
};

export const fetchObjectDataFromURL = async (url: string) => {
  const response = await axios.get(url);
  if (response.status !== 200) {
    console.log("Error fetching data");
    console.log(response);
    return {};
  }
  return response.data;
};
