-- ============================================
-- FEED RAW ADMIN DATA SCRIPT
-- ============================================

-- INSTRUCTIONS:
-- 1. Create a User in Supabase Auth if needed (or use existing).
-- 2. COPY the 'User UID' from Authentication > Users.
-- 3. PASTE it below replacing 'PASTE_YOUR_UID_HERE'.
-- 4. RUN this script in the SQL Editor.

-- STEP 1: Enable RLS (Should be on, but ensuring)
ALTER TABLE public.admins ENABLE ROW LEVEL SECURITY;

-- STEP 2: TEMPORARILY disable policy check for this insertion to "Feed" data without permissions loop
-- (Or just use the correct non-recursive policy we fixed earlier)
DROP POLICY IF EXISTS "Admins can view other admins" ON public.admins;
DROP POLICY IF EXISTS "Admins can view own record" ON public.admins;

CREATE POLICY "Admins can view own record" 
    ON public.admins FOR SELECT 
    USING (auth.uid() = id);

-- STEP 3: INSERT THE RAW DATA
-- Replace 'PASTE_YOUR_UID_HERE' with your actual UUID like 'a0eebc99-9c0b...'
INSERT INTO public.admins (id, email, role)
VALUES 
    ('PASTE_YOUR_UID_HERE', 'admin@troygh.com', 'superadmin')
ON CONFLICT (id) DO UPDATE 
SET role = 'superadmin';

-- STEP 4: VERIFY
SELECT * FROM public.admins;
