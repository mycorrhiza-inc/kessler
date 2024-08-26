export interface article {
  title: string;
  summary: string;
  text: string;
  sources: documentInfo[];
}

export interface documentInfo {
  title: string;
  display_text: string;
  summary: string;
  metadata: Object;
}
export const exampleArticle: article = {
  title: "The Evolution of AI",
  summary: "An overview of the advancements in artificial intelligence.",
  text: "Artificial intelligence has rapidly evolved over the past few decades, transforming various industries and creating new opportunities...",
  sources: [
    {
      title: "The Rise of Machine Learning",
      display_text: "Machine Learning advancements",
      summary:
        "An in-depth look at how machine learning has become a cornerstone of AI development.",
      metadata: {
        author: "Jane Doe",
        publicationDate: "2021-05-15",
        journal: "Tech Journal",
      },
    },
    {
      title: "Neural Networks Explained",
      display_text: "Understanding Neural Networks",
      summary:
        "A comprehensive guide to how neural networks function and their applications.",
      metadata: {
        author: "John Smith",
        publicationDate: "2020-11-01",
        journal: "AI Review",
      },
    },
  ],
};

export const exampleDocumentInfo: documentInfo = {
  title: "Deep Learning Innovations",
  display_text: "Innovations in Deep Learning",
  summary:
    "Exploring the latest innovations in deep learning techniques and their impact.",
  metadata: {
    author: "Alice Johnson",
    publicationDate: "2022-01-10",
    journal: "Deep Learning Today",
  },
};
