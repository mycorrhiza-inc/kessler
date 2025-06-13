// Figure out way to statically render this markdown without needing a UseState

import MarkdownRenderer from "../misc/MarkdownRenderer";


const Pricing = ({ className }: { className: string }) => {
  const pricing_tiers = [
    {
      key: "community",
      message: `
## $0 /month
**Community Tier**

- Access to search functionality
- Access to a limited set of open access government documents
`,
      buttonLink: "https://kessler.xyz/demo",
      buttonText: "Try Now",
      indicator: null,
    },
    {
      key: "professional",
      message: `
## ~~$50 /month~~ Free while in beta
**Professional Tier**

- Access to all our government documents
- Ability to upload and process your own documents
- Use of frontier level Large Language Models for rag functionality (Llama 405B, GPT-4o)
- Intended for Individuals or Non-Profits doing community work.
- Access to the source code under an MIT License!
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
    <div className={className}>
      <section className="overflow-hidden pb-20 pt-15 lg:pb-25 xl:pb-30">
        <div className="">
          <div className="animate_top mx-auto text-center">
            <div className="text-center">
              <h2 className="text-2xl font-bold">PRICING PLANS ()</h2>
              <p className="mt-4 text-xl">
                However, once we have built out the features and dataset, we are
                planning on charging to continue to support development.
              </p>
            </div>
          </div>
        </div>
        {/* This should be dryified quite a bit */}
        <div className="flex justify-center">
          <div className="flex flex-wrap justify-center gap-6 mt-15 max-w-[1207px] px-4 md:px-8 xl:mt-20 xl:px-0">
            {pricing_tiers.map(
              ({ key, message, buttonLink, buttonText, indicator }) => (
                <div
                  key={key}
                  className={`card grow border-secondary border-4  outline-secondary w-96 shadow-xl ${indicator ? "indicator" : ""}`}
                >
                  {indicator && (
                    <span className="indicator-item badge h-auto badge-accent mr-10 p-2">
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
        </div>
      </section>
    </div>
  );
};

export default Pricing;
