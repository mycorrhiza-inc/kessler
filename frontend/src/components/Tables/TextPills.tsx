import { AuthorInformation } from "@/lib/types/backend_schemas";
import Link from "next/link";

const oklchSubdivide = (colorNum: number, divisions?: number) => {
  const defaultDivisions = divisions || 15;
  const hue = (colorNum % defaultDivisions) * (360 / defaultDivisions);
  return `oklch(73% 0.123 ${hue})`;
};

const subdivide15 = [
  "oklch(73% 0.123 0)",
  "oklch(73% 0.123 30)",
  "oklch(73% 0.123 60)",
  "oklch(73% 0.123 90)",
  "oklch(73% 0.123 120)",
  "oklch(73% 0.123 150)",
  "oklch(73% 0.123 180)",
  "oklch(73% 0.123 210)",
  "oklch(73% 0.123 240)",
  "oklch(73% 0.123 270)",
  "oklch(73% 0.123 300)",
  "oklch(73% 0.123 330)",
];

type FileColor = {
  pdf: string;
  doc: string;
  xlsx: string;
};

const fileTypeColor = {
  pdf: "oklch(65.55% 0.133 0)",
  doc: "oklch(60.55% 0.13 240)",
  xlsx: "oklch(75.55% 0.133 140)",
};

const IsFiletypeColor = (key: string): key is keyof FileColor => {
  return key in fileTypeColor;
};

export const subdividedHueFromSeed = (seed: string): string => {
  const seed_integer =
    Math.abs(
      seed
        .split("")
        .reduce((acc, char) => (acc * 3 + 2 * char.charCodeAt(0)) % 27, 0),
    ) % 9;
  return oklchSubdivide(seed_integer, 15);
};

export const TextPill = ({
  text,
  href,
  seed,
}: {
  text?: string;
  href?: string;
  seed?: string;
}) => {
  const textDefined: string = text || "Unknown";
  const actualSeed = seed || textDefined;
  var pillColor = "";
  if (IsFiletypeColor(textDefined)) {
    pillColor = fileTypeColor[textDefined];
  } else {
    pillColor = subdividedHueFromSeed(actualSeed);
  }
  // btn-[${pillColor}]
  if (href) {
    return (
      <Link
        style={{ backgroundColor: pillColor }}
        className={`btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty	`}
        href={href}
      >
        {text}
      </Link>
    );
  }
  return (
    <button
      style={{ backgroundColor: pillColor }}
      className={`btn  btn-xs m-1 h-auto pb-1 no-animation text-black mt-2 mb-2 noclick text-pretty`}
    >
      {text}
    </button>
  );
};

export const AuthorInfoPill = ({
  author_info,
}: {
  author_info: AuthorInformation;
}) => (
  <TextPill
    text={author_info.author_name}
    seed={author_info.author_id}
    href={`/orgs/${author_info.author_id}`}
  />
);

export const DocketPill = ({
  docket_named_id,
}: {
  docket_named_id: string;
}) => <TextPill text={docket_named_id} href={`/dockets/${docket_named_id}`} />;
