-- ==========================================
-- SOURCES TABLE SETUP (UPDATED)
-- ==========================================

-- 1. Create or Alter table structure
-- We will handle both new creation and updating existing table
create table if not exists public.sources (
  id uuid default gen_random_uuid() primary key,
  name text not null,
  url text,
  category text,
  created_at timestamp with time zone default timezone('utc'::text, now()) not null,
  updated_at timestamp with time zone default timezone('utc'::text, now())
);

-- Safe migration: Add category if missing, drop logo_url if exists
do $$
begin
    if not exists (select 1 from information_schema.columns where table_name = 'sources' and column_name = 'category') then
        alter table public.sources add column category text;
    end if;
     if exists (select 1 from information_schema.columns where table_name = 'sources' and column_name = 'logo_url') then
        alter table public.sources drop column logo_url;
    end if;
end $$;


-- 2. Enable Realtime
begin;
  drop publication if exists supabase_realtime;
  create publication supabase_realtime for table public.sources;
commit;

-- 3. RLS
alter table public.sources enable row level security;

create policy "Enable read access for all users" 
on public.sources for select 
using (true);

create policy "Enable all access for authenticated users" 
on public.sources for all 
using (auth.role() = 'authenticated') 
with check (auth.role() = 'authenticated');
