import { AuthorInformation } from "@/lib/types/backend_schemas";
import Link from "next/link";
import { FaFilePdf, FaFileWord, FaFileExcel, FaHtml5 } from "react-icons/fa";
import { AiOutlineFileUnknown } from "react-icons/ai";
import { FileExtension } from "./FileExtension";
import { ReactNode, Dispatch, SetStateAction } from "react";

// Color generation utilities
const oklchHueSubdivide = (colorNum: number) => {
  const hue = (colorNum % HUE_DIVISONS) * (360 / HUE_DIVISONS);
  return `${hue})`;
};

const HUE_DIVISONS = 18;

export const subdividedHueFromSeed = (seed: string) => {
  const seed_integer = Math.abs(
    seed
      .split("")
      .reduce(
        (acc, char) => (1 + acc * 31 + 224 * char.charCodeAt(0)) % 217,
        0,
      ),
  );
  return oklchHueSubdivide(seed_integer);
}

export const subdividedColorFromSeed = (seed?: string): string => {
  if (seed === undefined) {
    return "oklch(80% 0.16 320)";
  }
  return `oklch(83 % 0.123 ${subdividedHueFromSeed(seed)})`
};

// File type colors
export const fileTypeColor: Record<FileExtension, string> = {
  [FileExtension.PDF]: "oklch(57.53% 0.1831 25.02)", // Adobe red
  [FileExtension.DOCX]: "oklch(70.55% 0.13 240)",
  [FileExtension.XLSX]: "oklch(75.55% 0.17 140)",
  [FileExtension.HTML]: "oklch(80.55% 0.08 60)",
  [FileExtension.UNKNOWN]: "oklch(60% 0.24 0)",
};

// Icon utilities
export const getExtensionIcon = (ext: FileExtension) => {
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

// Style constants
const PillButtonStyle = "inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs m-1";
const PillLinkStyle = "inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs m-1";

// Pill variants
type PillVariant = 'default' | 'file' | 'author' | 'docket' | 'custom';

// Pill size options
type PillSize = 'xs' | 'sm' | 'md';

const getSizeClasses = (size: PillSize) => {
  switch (size) {
    case 'xs':
      return "px-2 py-1 text-xs";
    case 'sm':
      return "px-3 py-1.5 text-sm";
    case 'md':
      return "px-4 py-2 text-base";
    default:
      return "px-2 py-1 text-xs";
  }
};

// Main Pill Props
interface BasePillProps {
  children: ReactNode;
  href?: string;
  onClick?: () => void;
  onRemove?: () => void;
  color?: string;
  textColor?: 'black' | 'white' | 'auto';
  size?: PillSize;
  className?: string;
  disabled?: boolean;
  removable?: boolean;
  variant?: PillVariant;
  'aria-label'?: string;
}

// Text-specific props
interface TextPillProps extends Omit<BasePillProps, 'children' | 'color'> {
  text?: string;
  seed?: string;
  placeholder?: string;
}

// File-specific props
interface FilePillProps extends Omit<BasePillProps, 'children' | 'color'> {
  extension: FileExtension;
  showIcon?: boolean;
}

// Author-specific props
interface AuthorPillProps extends Omit<BasePillProps, 'children' | 'color' | 'href'> {
  author: AuthorInformation;
  baseUrl?: string;
}

// Docket-specific props
export interface ConvoInfo {
  convo_id: string,
  convo_name: string,
  convo_number: string,
}
interface ConvoPillProps extends Omit<BasePillProps, 'children' | 'color' | 'href'> {
  convo_info: ConvoInfo;
  baseUrl?: string;
}


// Utility to determine text color based on background
const getAutoTextColor = (backgroundColor: string): 'black' | 'white' => {
  // Extract lightness from oklch color
  const lightnessMatch = backgroundColor.match(/oklch\((\d+(?:\.\d+)?)%/);
  if (lightnessMatch) {
    const lightness = parseFloat(lightnessMatch[1]);
    return lightness > 60 ? 'black' : 'white';
  }
  return 'black'; // Default fallback
};

// Base Pill Component
export const BasePill = ({
  children,
  href,
  onClick,
  onRemove,
  color = "oklch(75% 0.16 000)",
  textColor = 'auto',
  size = 'xs',
  className = '',
  disabled = false,
  removable = false,
  variant = 'default',
  'aria-label': ariaLabel,
  ...props
}: BasePillProps) => {
  const actualTextColor = textColor === 'auto' ? getAutoTextColor(color) : textColor;
  const sizeClasses = getSizeClasses(size);
  const baseClasses = `inline-flex items-center rounded-full transition-colors ${sizeClasses}`;

  const pillContent = (
    <span className={`inline-flex items-center ${removable ? 'gap-1' : ''}`}>
      {children}
      {removable && onRemove && (
        <span
          className="rounded-full p-0.5 transition-colors cursor-pointer hover:opacity-80 ml-1"
          style={{
            backgroundColor: 'transparent'
          }}
          onMouseEnter={(e) => {
            // Create slightly more saturated hover color
            const hueMatch = color.match(/oklch\(83% 0\.123 (\d+)\)/);
            if (hueMatch) {
              const hue = hueMatch[1];
              e.currentTarget.style.backgroundColor = `oklch(75% 0.18 ${hue})`;
            }
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = 'transparent';
          }}
          onClick={(e) => {
            e.stopPropagation();
            e.preventDefault();
            onRemove();
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') {
              e.preventDefault();
              e.stopPropagation();
              onRemove();
            }
          }}
          role="button"
          tabIndex={0}
          aria-label="Remove"
        >
          <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </span>
      )}
    </span>
  );

  const commonStyle = {
    backgroundColor: color,
    color: actualTextColor
  };

  const commonClasses = `${baseClasses} ${className} ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer hover:opacity-90'}`;

  if (href && !disabled) {
    return (
      <Link
        href={href}
        className={commonClasses}
        style={commonStyle}
        aria-label={ariaLabel}
        {...props}
      >
        {pillContent}
      </Link>
    );
  }

  if (onClick && !disabled) {
    return (
      <button
        type="button"
        onClick={onClick}
        className={commonClasses}
        style={commonStyle}
        disabled={disabled}
        aria-label={ariaLabel}
        {...props}
      >
        {pillContent}
      </button>
    );
  }

  return (
    <span
      className={commonClasses}
      style={commonStyle}
      aria-label={ariaLabel}
      {...props}
    >
      {pillContent}
    </span>
  );
};

// Text Pill Component
export const TextPill = ({
  text,
  seed,
  placeholder = "Unknown",
  ...props
}: TextPillProps) => {
  const textDefined = text || placeholder;
  const actualSeed = seed || textDefined;
  const pillColor = subdividedColorFromSeed(actualSeed);

  return (
    <BasePill color={pillColor} variant="default" {...props}>
      {textDefined}
    </BasePill>
  );
};

// File Extension Pill Component
export const FilePill = ({
  extension,
  showIcon = true,
  textColor = 'white',
  ...props
}: FilePillProps) => {
  const color = fileTypeColor[extension];
  const icon = getExtensionIcon(extension);

  return (
    <BasePill color={color} textColor={textColor} variant="file" {...props}>
      <span className="flex items-center">
        {showIcon && <span className="mr-2">{icon}</span>}
        <span>{FileExtension[extension].toUpperCase()}</span>
      </span>
    </BasePill>
  );
};

// Author Pill Component
export const AuthorPill = ({
  author,
  baseUrl = "/orgs",
  ...props
}: AuthorPillProps) => {
  const href = `${baseUrl}/${author.author_id}`;

  return (
    <TextPill
      text={author.author_name}
      seed={author.author_id}
      href={href}
      variant="author"
      aria-label={`View ${author.author_name}'s profile`}
      {...props}
    />
  );
};

// Docket Pill Component
export const ConversationPill = ({
  convo_info,
  baseUrl = "/dockets",
  ...props
}: ConvoPillProps) => {
  const displayText = convo_info.convo_number || convo_info.convo_name || convo_info.convo_id || ""
  const href = `${baseUrl}/${convo_info.convo_id}`;

  return (
    <TextPill
      text={displayText}
      seed={convo_info.convo_id}
      href={href}
      variant="docket"
      aria-label={`View docket ${displayText}`}
      {...props}
    />
  );
};

// Legacy component aliases for backward compatibility
export const RawPill = BasePill;
export const ExtensionPill = FilePill;
export const AuthorInfoPill = ({ author_info, ...props }: { author_info: AuthorInformation } & Omit<AuthorPillProps, 'author'>) => (
  <AuthorPill author={author_info} {...props} />
);

// Utility function for creating custom colored pills
export const createColoredPill = (color: string) => (props: Omit<BasePillProps, 'color'>) => (
  <BasePill color={color} {...props} />
);

// Predefined color pills
export const RedPill = createColoredPill("oklch(70% 0.15 25)");
export const BluePill = createColoredPill("oklch(70% 0.15 240)");
export const GreenPill = createColoredPill("oklch(70% 0.15 140)");
export const YellowPill = createColoredPill("oklch(80% 0.12 85)");
export const PurplePill = createColoredPill("oklch(70% 0.15 300)");

// Bulk operations for multiple pills
export const PillGroup = ({
  children,
  className = "flex flex-wrap gap-1"
}: {
  children: ReactNode;
  className?: string;
}) => (
  <div className={className}>
    {children}
  </div>
);

// Hook for managing removable pills
// Removing for now, all these pills should be renderable on the server.
// export const useRemovablePills = <T extends { id: string | number }>(
//   initialItems: T[]
// ): {
//   items: T[];
//   removeItem: (id: string | number) => void;
//   addItem: (item: T) => void;
//   clearAll: () => void;
//   setItems: Dispatch<SetStateAction<T[]>>;
// } => {
//   const [items, setItems] = useState<T[]>(initialItems);
//
//   const removeItem = (id: string | number): void => {
//     setItems(prev => prev.filter(item => item.id !== id));
//   };
//
//   const addItem = (item: T): void => {
//     setItems(prev => [...prev, item]);
//   };
//
//   const clearAll = (): void => {
//     setItems([]);
//   };
//
//   return {
//     items,
//     removeItem,
//     addItem,
//     clearAll,
//     setItems
//   };
// };
