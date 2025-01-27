-- +goose Up
CREATE TABLE IF NOT EXISTS public.jobs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	job_priority INT NOT NULL,
	job_name VARCHAR(255) NOT NULL,
	job_status VARCHAR(255) NOT NULL,
	job_type VARCHAR(255) NOT NULL,
	job_data JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS public.jobs_log (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	job_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(255) NOT NULL,
	message TEXT,
	FOREIGN KEY (job_id) REFERENCES public.jobs(id) ON DELETE CASCADE
);

CREATE INDEX idx_jobs_priority ON public.jobs(job_priority);


-- +goose Down
DROP TABLE IF EXISTS public.jobs CASCADE;
DROP TABLE IF EXISTS public.jobs_log CASCADE;
DROP INDEX IF EXISTS idx_jobs_priority CASCADE;
