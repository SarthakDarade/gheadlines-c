-- FIX SQL: Run this in Supabase SQL Editor to fix login permissions

-- 1. Drop the potentially problematic recursive policy
DROP POLICY IF EXISTS "Admins can view other admins" ON public.admins;

-- 2. Create a simpler policy: Users can ALWAYS see their OWN admin record
-- This is necessary for the initial check "Am I an admin?"
CREATE POLICY "Admins can view own record" 
    ON public.admins FOR SELECT 
    USING (auth.uid() = id);

-- 3. Ensure RLS is enabled
ALTER TABLE public.admins ENABLE ROW LEVEL SECURITY;

-- 4. Grant permissions
GRANT SELECT ON public.admins TO authenticated;
GRANT SELECT ON public.admins TO anon;

-- 5. Verification (Optional - run this to check if you see your row)
-- SELECT * FROM public.admins;
