-- Create breaking_news table
CREATE TABLE IF NOT EXISTS public.breaking_news (
  id uuid not null default gen_random_uuid (),
  headline text not null,
  url text null,
  category text null,
  created_at timestamp with time zone null default now(),
  constraint breaking_news_pkey primary key (id)
);

-- Enable RLS
ALTER TABLE public.breaking_news ENABLE ROW LEVEL SECURITY;

-- Policies for Admins (using the is_admin function we defined earlier for robust access)
DROP POLICY IF EXISTS "Admins can do everything on breaking_news" ON public.breaking_news;
CREATE POLICY "Admins can do everything on breaking_news"
    ON public.breaking_news FOR ALL
    USING (public.is_admin());

-- Policy for Public Read (if needed for the main site)
DROP POLICY IF EXISTS "Breaking news viewable by everyone" ON public.breaking_news;
CREATE POLICY "Breaking news viewable by everyone"
    ON public.breaking_news FOR SELECT
    USING (true);
