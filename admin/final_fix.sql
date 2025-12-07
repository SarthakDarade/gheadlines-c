-- FINAL FIX: ROBUST ADMIN ACCESS
-- Using SECURITY DEFINER function to bypass RLS recursion headaches completely.

-- 1. Create a secure function to check admin status
-- This function runs with "superuser" privileges (SECURITY DEFINER) 
-- so it can always read the admins table regardless of RLS.
CREATE OR REPLACE FUNCTION public.is_admin()
RETURNS BOOLEAN AS $$
BEGIN
  RETURN EXISTS (
    SELECT 1 FROM public.admins WHERE id = auth.uid()
  );
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- 2. Update Policies to use this function
-- This guarantees that if you are in the admins table, you see EVERYTHING.

-- PROFILES
DROP POLICY IF EXISTS "Admins can do everything on profiles" ON public.profiles;
CREATE POLICY "Admins can do everything on profiles"
    ON public.profiles FOR ALL
    USING (public.is_admin());

-- ARTICLES
DROP POLICY IF EXISTS "Admins can do everything on articles" ON public.articles;
CREATE POLICY "Admins can do everything on articles"
    ON public.articles FOR ALL
    USING (public.is_admin());

-- CATEGORIES
DROP POLICY IF EXISTS "Admins can do everything on categories" ON public.categories;
CREATE POLICY "Admins can do everything on categories"
    ON public.categories FOR ALL
    USING (public.is_admin());

-- TRENDING NEWS
DROP POLICY IF EXISTS "Admins can do everything on trending_news" ON public.trending_news;
CREATE POLICY "Admins can do everything on trending_news"
    ON public.trending_news FOR ALL
    USING (public.is_admin());

-- LIVE UPDATES
DROP POLICY IF EXISTS "Admins can do everything on live_updates" ON public.live_updates;
CREATE POLICY "Admins can do everything on live_updates"
    ON public.live_updates FOR ALL
    USING (public.is_admin());

-- USERS (Wait, profiles is users. But if you have other tables like 'newsletter_subscribers')
DROP POLICY IF EXISTS "Admins can do everything on newsletter_subscribers" ON public.newsletter_subscribers;
CREATE POLICY "Admins can do everything on newsletter_subscribers"
    ON public.newsletter_subscribers FOR ALL
    USING (public.is_admin());

DROP POLICY IF EXISTS "Admins can do everything on contact_messages" ON public.contact_messages;
CREATE POLICY "Admins can do everything on contact_messages"
    ON public.contact_messages FOR ALL
    USING (public.is_admin());


-- 3. Insert Mock Data (Only if empty) to ensure you see SOMETHING
-- Articles
INSERT INTO public.articles (title, description, content, slug, published, created_at)
SELECT 'Welcome to TroyGH Admin', 'This is your first article. Edit or delete it.', '<p>Welcome.</p>', 'welcome-admin', true, NOW()
WHERE NOT EXISTS (SELECT 1 FROM public.articles);

-- Trending
INSERT INTO public.trending_news (title, summary, category)
SELECT 'System Operational', 'The admin panel is now fully connected.', 'System'
WHERE NOT EXISTS (SELECT 1 FROM public.trending_news);
