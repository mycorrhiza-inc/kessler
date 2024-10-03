"use client";
import React from "react";
import Image from "next/image";
import { motion } from "framer-motion";
type Feature = {
  id: number;
  icon: string;
  title: string;
  description: string;
};

const featuresData: Feature[] = [
  {
    id: 1,
    icon: "/landing/images/icon/icon-01.svg",
    title: "Extensive Government Document Database",
    description:
      "Search through a vast collection of over 10,000 government documents, from safety reports to public comments on zoning changes.",
  },
  {
    id: 2,
    icon: "/landing/images/icon/icon-02.svg",
    title: "Advanced Search and Filtering",
    description:
      "Utilize advanced search capabilities identical to most government dockets, including search by title and filter by author.",
  },
  {
    id: 3,
    icon: "/landing/images/icon/icon-03.svg",
    title: "Modern User Interface",
    description:
      "Experience a user interface that is intuitive and modern, far from the .ASP forms that havent been updated since 2003.",
  },
  {
    id: 4,
    icon: "/landing/images/icon/icon-03.svg",
    title: "Multi-format Ingestion",
    description:
      "Ingest a variety of document formats such as PDFs, audio recordings, Word documents, and eBooks.Upload your own documents in bulk and integrate them seamlessly with all existing features.",
  },
  {
    id: 5,
    icon: "/landing/images/icon/icon-04.svg",
    title: "Full-text Search",
    description:
      "Perform full-text searches across all ingested documents for comprehensive data retrieval.",
  },
  {
    id: 6,
    icon: "/landing/images/icon/icon-05.svg",
    title: "AI-powered Interaction",
    description:
      "Engage with the docket using a chatbot that leverages Large Language Models (LLMs) and Retrieval-Augmented Generation (RAG).",
  },
  {
    id: 7,
    icon: "/landing/images/icon/icon-06.svg",
    title: "Continuous Improvements",
    description:
      " In addition to adding new features our goal is to eventually include every publicly accessible government document in our database. Stay tuned for more updates! ",
  },
];

const futureFeaturesData: Feature[] = [
  {
    id: 9,
    icon: "/landing/images/icon/icon-09.svg",
    title: "Document Sharing and Revenue Split",
    description:
      "Share your government documents with other users and earn revenue from making that data accessible.",
  },
  {
    id: 10,
    icon: "/landing/images/icon/icon-10.svg",
    title: "Author Submission Correlation",
    description:
      "Correlate submissions by a single author across multiple document databases and jurisdictions.",
  },
  {
    id: 11,
    icon: "/landing/images/icon/icon-11.svg",
    title: "Smart Notifications",
    description:
      "Receive notifications only for important new submissions, ensuring you stay updated without the noise.",
  },
  {
    id: 12,
    icon: "/landing/images/icon/icon-12.svg",
    title: "Sentiment Analysis",
    description:
      "Easily identify whether documents support or oppose a certain proposal with our sentiment analysis feature.",
  },
];

// <div className="card bg-base-100 w-96 shadow-xl">
//   <figure className="px-10 pt-10">
//     <img
//       src="https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"
//       alt="Shoes"
//       className="rounded-xl" />
//   </figure>
//   <div className="card-body items-center text-center">
//     <h2 className="card-title">Shoes!</h2>
//     <p>If a dog chews shoes whose shoes does he choose?</p>
//     <div className="card-actions">
//       <button className="btn btn-primary">Buy Now</button>
//     </div>
//   </div>
// </div>

//const SingleFeature = ({ feature }: { feature: Feature }) => {
//  const { icon, title, description } = feature;
//
//  return (
//    <>
//      <motion.div
//        variants={{
//          hidden: {
//            opacity: 0,
//            y: -10,
//          },
//
//          visible: {
//            opacity: 1,
//            y: 0,
//          },
//        }}
//        initial="hidden"
//        whileInView="visible"
//        transition={{ duration: 0.5 }}
//        viewport={{ once: true }}
//        className="animate_top z-40 rounded-lg border border-white bg-white p-7.5 shadow-solid-3 transition-all hover:shadow-solid-4 dark:border-strokedark dark:bg-blacksection dark:hover:bg-hoverdark xl:p-12.5"
//      >
//        <div className="relative flex h-16 w-16 items-center justify-center rounded-[4px] bg-primary">
//          <Image src={icon} width={36} height={36} alt="title" />
//        </div>
//        <h3 className="mb-5 mt-7.5 text-xl font-semibold text-black dark:text-white xl:text-itemtitle">
//          {title}
//        </h3>
//        <p>{description}</p>
//      </motion.div>
//    </>
//  );
//};
const SingleFeature = ({ feature }: { feature: Feature }) => {
  const { icon, title, description } = feature;

  return (
    <div className="card bg-primary w-96 shadow-xl text-neutral-content">
      <figure className="px-10 pt-10">
        <Image
          src={icon}
          width={36}
          height={36}
          alt={title}
          className="rounded-xl"
        />
      </figure>
      <div className="card-body items-center text-center text-primary-content">
        <h2 className="card-title">{title}</h2>
        <p>{description}</p>
      </div>
    </div>
  );
};

const Feature = () => {
  return (
    <>
      {/* <!-- ===== Features Start ===== --> */}
      <section id="features" className="py-20 lg:py-25 xl:py-30">
        <div className="mx-auto max-w-c-1315 px-4 md:px-8 xl:px-0">
          {/* <!-- Section Title Start --> */}
          <div className="text-center">
            <h2 className="text-4xl font-bold">SOLID FEATURES</h2>
            <h3 className="text-xl mt-4">Core Features of Kessler</h3>
            <p className="mt-2">
              Kessler is still in beta, but this is what we have implemented so
              far.
            </p>
          </div>
          {/* <!-- Section Title End --> */}

          <div className="mt-12.5 grid grid-cols-1 gap-7.5 md:grid-cols-2 lg:mt-15 lg:grid-cols-3 xl:mt-20 gap-5">
            {/* <!-- Features item Start --> */}

            {featuresData.map((feature, key) => (
              <SingleFeature feature={feature} key={key} />
            ))}
            {/* <!-- Features item End --> */}
          </div>
        </div>
      </section>

      {/* <!-- ===== Features End ===== --> */}
    </>
  );
};

export default Feature;
