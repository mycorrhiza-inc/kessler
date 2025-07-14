export interface TypedFilter {
  org_id?: string;
  convo_id?: string;
}

export function toFilterRecords(
  typed_filter: TypedFilter,
): Record<string, string> {
  return Object.fromEntries(
    Object.entries(typed_filter).filter(
      ([_, value]) => value !== undefined && value !== "",
    ),
  );
}
