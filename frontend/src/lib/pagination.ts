export const queryStringFromLimitOffset = (limit: number, offset: number) => {
  return `?limit=${limit}&offset=${offset}`;
};

export const queryStringFromPageMaxHits = (page: number, maxHits?: number) => {
  const limit = maxHits || 40;
  const offset = page * limit;
  return `?limit=${limit}&offset=${offset}`;
};
