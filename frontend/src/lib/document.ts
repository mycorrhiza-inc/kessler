import axios from "axios";

export interface DocumentInterface {
  docid: string;
  title?: string;
  docText?: string;
  pdfUrl?: string;
  docMetadata?: any;
}

export type DocidMap = {
  [key: string]: Document;
};

export class Document implements DocumentInterface {
  docid: string = "";
  title?: string;
  docText?: string;
  pdfUrl?: string;
  docMetadata?: any;
  loaded: boolean = false;

  constructor({
    docid,
    title,
    docText,
    pdfUrl,
    docMetadata,
  }: DocumentInterface) {
    this.docid = docid;
    if (title) {
      this.title = title;
    }
    if (docText) {
      this.docText = docText;
    }
    if (pdfUrl) {
      this.pdfUrl = pdfUrl;
    }
    if (docMetadata) {
      this.docMetadata = docMetadata;
    }
  }

  getDocumentMetadata = async () => {
    const response = await axios.get(`/api/v1/files/metadata/${this.docid}`);
    this.docMetadata = response.data;
    console.log(this.docMetadata);
  };

  getDocumentText = async () => {
    const response = await axios.get(`/api/v1/files/markdown/${this.docid}`);
    this.docText = response.data;
  };

  getPdfUrl = async () => {
    this.pdfUrl = `/api/v1/files/raw/${this.docid}`;
  };

  loadDocument = async () => {
    await this.getDocumentMetadata();
    await this.getDocumentText();
    await this.getPdfUrl();
    this.loaded = true;
  };
}
