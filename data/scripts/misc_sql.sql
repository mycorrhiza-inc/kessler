-- Sql command for determining how many processed documents since date.
SELECT
    count(id)
FROM
    public.file
WHERE
    verified = TRUE
    AND updated_at >= `2024-12-01`;
