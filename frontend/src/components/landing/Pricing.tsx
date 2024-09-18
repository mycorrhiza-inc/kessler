// Figure out way to statically render this markdown without needing a UseState
import Image from "next/image";

import MarkdownRenderer from "../MarkdownRenderer";

const Pricing = () => {
  const isUser = true; // Example boolean, replace with actual logic
  const message = {
    community: `
## $0 /month
**Community Tier**

- Access to search functionality
- Access to a limited set of open access government documents
- [Run your own instance with our open source code](https://git.mycor.io/mycorrhiza/kessler)
`,
    professional: `
## $30 /month
**Professional Tier**

- Access to all our government documents
- Ability to upload and process your own documents
- Use of frontier level Large Language Models for rag functionality (Llama 405B, GPT-4o)
- Intended for Individuals or Non-Profits doing community work.
`,

    enterprise: `
## Contact Us
**Enterprise Tier**

- Priority support for adding your custom datasets
- On-Premises Deployment
- Access to our raw government datasets
- Intended for Large Companies or National/International NGO's
`,
  };

  return (
    <>
      <section className="overflow-hidden pb-20 pt-15 lg:pb-25 xl:pb-30">
        <div className="mx-auto max-w-c-1315 px-4 md:px-8 xl:px-0">
          <div className="animate_top mx-auto text-center">
            <div className="text-center">
              <h2 className="text-2xl font-bold">PRICING PLANS</h2>
              <p className="mt-4 text-gray-600">
                Pricing that fits your needs and helps fund future development.
              </p>
            </div>
          </div>
        </div>

        <div className="relative mx-auto mt-15 max-w-[1207px] px-4 md:px-8 xl:mt-20 xl:px-0">
          <div className="absolute -bottom-15 -z-1 h-full w-full">
            <Image
              fill
              src="./images/shape/shape-dotted-light.svg"
              alt="Dotted"
              className="dark:hidden"
            />
          </div>

          <div className="flex flex-wrap justify-center gap-7.5 lg:flex-nowrap xl:gap-12.5">
            <div className="card bg-base-100 w-96 shadow-xl">
              <div className="card-body">
                <MarkdownRenderer>{message.community}</MarkdownRenderer>
                <div className="card-actions justify-end">
                  <a href="https://app.kessler.xyz" className="btn btn-primary">
                    Try Now
                  </a>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 w-96 shadow-xl">
              <div className="card-body">
                <MarkdownRenderer>{message.professional}</MarkdownRenderer>
                <div className="card-actions justify-end">
                  <a href="/payment" className="btn btn-primary">
                    Purchase
                  </a>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 w-96 shadow-xl">
              <div className="card-body">
                <MarkdownRenderer>{message.enterprise}</MarkdownRenderer>
                <div className="card-actions justify-end">
                  <a href="/contact" className="btn btn-primary">
                    Contact Us
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </>
  );
};

export default Pricing;
