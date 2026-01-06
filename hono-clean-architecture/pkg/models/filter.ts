export interface Search {
  keyword?: string;
  column?: string;
  text?: string;
}

export function createSearch(keyword?: string, column?: string, text?: string): Search {
  return {
    keyword: keyword || undefined,
    column: column || undefined,
    text: text || undefined,
  };
}

export function getSearchString(search: Search): string {
  return `keyword=${search.keyword || ""}&column=${search.column || ""}&text=${search.text || ""}`;
}
