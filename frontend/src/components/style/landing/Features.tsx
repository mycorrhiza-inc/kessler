"use client";
import React from "react";
import { Meteors } from "../aceternity/meteors";
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
  // {
  //   id: 9,
  //   icon: "/landing/images/icon/icon-09.svg",
  //   title: "Document Sharing and Revenue Split",
  //   description:
  //     "Share your government documents with other users and earn revenue from making that data accessible.",
  // },
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
    <div className="">
      <div className=" w-full relative max-w-xs">
        <div className="absolute inset-0 h-full w-full bg-linear-to-r from-blue-500 to-teal-500 transform scale-[0.80] bg-red-500 rounded-full blur-3xl" />
        <div className="relative shadow-xl bg-gray-900 border border-gray-800  px-4 py-8 h-full overflow-hidden rounded-2xl flex flex-col justify-end items-start">
          <h1 className="font-bold text-xl text-white mb-4 relative z-50">
            {title}
          </h1>

          <p className="font-normal text-base text-slate-500 mb-4 relative z-50">
            {description}
          </p>

          {/* Meaty part - Meteor effect */}
          <Meteors number={20} />
        </div>
      </div>
    </div>
  );
};
// const SingleFeature = ({ feature }: { feature: Feature }) => {
//   const { icon, title, description } = feature;
//   return (
//     <div className="rounded-box  w-auto p-10 m-5 bg-primary shadow-xl text-neutral-content">
//       <h2 className="card-title">{title}</h2>
//       <p>{description}</p> <Meteors number={20} />
//     </div>
//   );
// };

const Feature = ({ className }: { className: string }) => {
  return (
    <div className={className}>
      {/* <!-- ===== Features Start ===== --> */}
      <section id="features">
        <div className="">
          {/* <!-- Section Title Start --> */}
          <div className="text-center">
            <h2 className="text-4xl font-bold">
              <br />
              FEATURES
            </h2>
            {/* <h3 className="text-xl mt-4">Core Features of Kessler</h3> */}
            <p className="mt-2">
              Kessler is still in beta, but this is what we have implemented so
              far.
              <br />
              <br />
              <br />
            </p>
          </div>
          {/* <!-- Section Title End --> */}

          <div className="grid grid-cols-1  lg:grid-cols-3 justify-items-center gap-x-1 gap-y-9">
            {/* <!-- Features item Start --> */}

            {featuresData.map((feature, key) => (
              <SingleFeature feature={feature} key={key} />
            ))}
            {/* <!-- Features item End --> */}
          </div>
        </div>
      </section>

      {/* <!-- ===== Features End ===== --> */}
    </div>
  );
};

export default Feature;
