-- Sql command for determining how many processed documents since date.
SELECT
    count(id)
FROM
    public.file
WHERE
    verified = TRUE
    AND updated_at >= '2024-12-01';

-- Sql command for seeing the most recent documents.
SELECT
    *
FROM
    public.file
WHERE
    verified = TRUE
    AND updated_at >= '2024-12-01'
ORDER BY
    updated_at DESC;