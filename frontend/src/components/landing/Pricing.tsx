// Figure out way to statically render this markdown without needing a UseState
import Image from "next/image";

import MarkdownRenderer from "../MarkdownRenderer";
import { BadgeIcon } from "lucide-react";

const Pricing = () => {
  const pricing_tiers = [
    {
      key: "community",
      message: `
## $0 /month
**Community Tier**

- Access to search functionality
- Access to a limited set of open access government documents
- [Run your own instance with our open source code](https://git.mycor.io/mycorrhiza/kessler)
`,
      buttonLink: "https://kessler.xyz",
      buttonText: "Try Now",
      indicator: null,
    },
    {
      key: "professional",
      message: `
## $50 /month
**Professional Tier**

- Access to all our government documents
- Ability to upload and process your own documents
- Use of frontier level Large Language Models for rag functionality (Llama 405B, GPT-4o)
- Intended for Individuals or Non-Profits doing community work.
`,
      buttonLink: "/payment",
      buttonText: "Purchase",
      indicator: "Most Popular!",
    },
    {
      key: "enterprise",
      message: `
## Contact Us
**Enterprise Tier**

- Priority support for adding your custom datasets
- On-Premises Deployment
- Access to our raw government datasets
- Intended for Large Companies or National/International NGO's
`,
      buttonLink: "/contact",
      buttonText: "Contact Us",
      indicator: null,
    },
  ];

  return (
    <>
      <section className="overflow-hidden pb-20 pt-15 lg:pb-25 xl:pb-30">
        <div className="mx-auto max-w-c-1315 px-4 md:px-8 xl:px-0">
          <div className="animate_top mx-auto text-center">
            <div className="text-center">
              <h2 className="text-2xl font-bold">PRICING PLANS</h2>
              <p className="mt-4">
                Pricing that fits your needs and helps fund future development.
              </p>
            </div>
          </div>
        </div>
        {/* This should be dryified quite a bit */}
        <div className="flex flex-wrap justify-center gap-6 mt-15 max-w-[1207px] px-4 md:px-8 xl:mt-20 xl:px-0">
          {pricing_tiers.map(
            ({ key, message, buttonLink, buttonText, indicator }) => (
              <div
                key={key}
                className={`card border-secondary border-2  outline-secondary w-96 shadow-xl ${indicator ? "indicator" : ""}`}
              >
                {indicator && (
                  <span className="indicator-item badge h-auto badge-accent mr-10">
                    {indicator}
                  </span>
                )}
                <div className="card-body">
                  <MarkdownRenderer>{message}</MarkdownRenderer>
                  <div className="card-actions justify-end">
                    <a href={buttonLink} className="btn btn-accent">
                      {buttonText}
                    </a>
                  </div>
                </div>
              </div>
            ),
          )}
        </div>
      </section>
    </>
  );
};

export default Pricing;
