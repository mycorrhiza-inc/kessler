import { AuthorInformation } from "@/lib/types/backend_schemas";
import Link from "next/link";
import { FaFilePdf, FaFileWord, FaFileExcel, FaHtml5 } from "react-icons/fa";
import { AiOutlineFileUnknown } from "react-icons/ai";

const oklchSubdivide = (colorNum: number, divisions?: number) => {
  const defaultDivisions = divisions || 18;
  const hue = (colorNum % defaultDivisions) * (360 / defaultDivisions);
  return `oklch(83% 0.123 ${hue})`;
};

const subdivide20 = [
  "oklch(83% 0.123 0)",
  "oklch(83% 0.123 18)",
  "oklch(83% 0.123 36)",
  "oklch(83% 0.123 54)",
  "oklch(83% 0.123 72)",
  "oklch(83% 0.123 90)",
  "oklch(83% 0.123 108)",
  "oklch(83% 0.123 126)",
  "oklch(83% 0.123 144)",
  "oklch(83% 0.123 162)",
  "oklch(83% 0.123 180)",
  "oklch(83% 0.123 198)",
  "oklch(83% 0.123 216)",
  "oklch(83% 0.123 234)",
  "oklch(83% 0.123 252)",
  "oklch(83% 0.123 270)",
  "oklch(83% 0.123 288)",
  "oklch(83% 0.123 306)",
  "oklch(83% 0.123 324)",
  "oklch(83% 0.123 342)",
];

type FileColor = {
  pdf: string;
  doc: string;
  xlsx: string;
};

const IsFiletypeColor = (key: string): key is keyof FileColor => {
  return key in fileTypeColor;
};

export const subdividedHueFromSeed = (seed?: string): string => {
  if (seed === undefined) {
    return "oklch(40% 0.33 0)";
  }
  const seed_integer = Math.abs(
    seed
      .split("")
      .reduce(
        (acc, char) => (1 + acc * 31 + 224 * char.charCodeAt(0)) % 217,
        0,
      ),
  );
  return oklchSubdivide(seed_integer, 18);
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
  const pillColor = subdividedHueFromSeed(actualSeed);
  // btn-[${pillColor}]
  return (
    <RawPill color={pillColor} href={href}>
      {textDefined}
    </RawPill>
  );
};

export const RawPill = ({
  children,
  color,
  href,
}: {
  children: React.ReactNode;
  color: string;
  href?: string;
}) => {
  if (href) {
    return (
      <Link
        style={{ backgroundColor: color }}
        className={`btn btn-xs m-1 h-auto pb-1 text-black noclick text-pretty	`}
        href={href}
      >
        {children}
      </Link>
    );
  }
  return (
    <button
      style={{ backgroundColor: color }}
      className={`btn btn-xs m-1 h-auto pb-1 no-animation text-black mt-2 mb-2 noclick text-pretty`}
    >
      {children}
    </button>
  );
};

const getIcon = (ext: FileExtension) => {
  switch (ext) {
    case FileExtension.PDF:
      return <FaFilePdf />;
    case FileExtension.DOCX:
      return <FaFileWord />;
    case FileExtension.XLSX:
      return <FaFileExcel />;
    case FileExtension.HTML:
      return <FaHtml5 />;
    case FileExtension.UNKNOWN:
    default:
      return <AiOutlineFileUnknown />;
  }
};
export enum FileExtension {
  PDF,
  XLSX,
  DOCX,
  HTML,
  UNKNOWN,
}
const fileTypeColor: Record<FileExtension, string> = {
  [FileExtension.PDF]: "oklch(65.55% 0.133 0)",
  [FileExtension.DOCX]: "oklch(70.55% 0.13 240)",
  [FileExtension.XLSX]: "oklch(75.55% 0.133 140)",
  [FileExtension.HTML]: "oklch(80.55% 0.08 60)",
  [FileExtension.UNKNOWN]: "oklch(60% 0.3 0)",
};

export const ExtensionPill = ({ ext }: { ext: FileExtension }) => {
  const colorString = fileTypeColor[ext];
  const icon = getIcon(ext);

  return (
    <RawPill color={colorString}>
      <span className="flex items-center">
        <span className="mr-2">{icon}</span>
        {FileExtension[ext].toUpperCase()}
      </span>
    </RawPill>
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
  text,
}: {
  docket_named_id: string;
  text?: string;
}) => (
  <TextPill
    text={text || docket_named_id}
    href={`/dockets/${docket_named_id}`}
  />
);
