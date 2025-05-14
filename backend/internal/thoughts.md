# Ideas for Hash-based Markdown Retrieval Endpoint

- Utilize existing `HashGetUUIDsFile` function to map file hash to UUID(s).
- For simplicity, select the first matching UUID if multiple are returned.
- Reuse `GetSpecificFileText` to generate markdown output from the file texts.
- Support optional query parameters:
  - `match_lang`: filter for a specific language text.
  - `original_lang`: boolean flag to retrieve the original text instead of translations.
- Implement a new HTTP handler `FileMarkdownByHashHandler` in `file_read_handler.go`.
- Keep code minimal and non-disruptive to existing infrastructure.
- Return `404 Not Found` if no file UUIDs or markdown text is found.
- Return `500 Internal Server Error` for unexpected errors.
- Document endpoint design and usage in `hash-text-retrival.md`.
- (Future) Wire up the new handler to the router under path `/files/markdown/hash/{hash}`.
