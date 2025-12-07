import { createContext, useContext, useEffect, useState } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { Session, User } from '@supabase/supabase-js'
import { useRouter } from 'next/router'

type AdminRole = 'superadmin' | 'editor' | 'reporter'

interface AdminProfile {
    id: string
    email: string
    role: AdminRole
    created_at: string
}

interface AuthContextType {
    user: User | null
    session: Session | null
    profile: AdminProfile | null
    isLoading: boolean
    isAdmin: boolean
    signOut: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null)
    const [session, setSession] = useState<Session | null>(null)
    const [profile, setProfile] = useState<AdminProfile | null>(null)
    const [isLoading, setIsLoading] = useState(true)
    const router = useRouter()

    useEffect(() => {
        // Check active session
        supabase.auth.getSession().then(({ data: { session } }) => {
            setSession(session)
            setUser(session?.user ?? null)
            if (session?.user) {
                fetchAdminProfile(session.user)
            } else {
                setIsLoading(false)
            }
        })

        const { data: { subscription } } = supabase.auth.onAuthStateChange((_event, session) => {
            setSession(session)
            setUser(session?.user ?? null)
            if (session?.user) {
                fetchAdminProfile(session.user)
            } else {
                setProfile(null)
                setIsLoading(false)
            }
        })

        return () => subscription.unsubscribe()
    }, [])

    const fetchAdminProfile = async (user: User) => {
        try {
            // Check if user exists in 'admins' table
            const { data, error } = await supabase
                .from('admins')
                .select('*')
                .eq('id', user.id) // Assuming auth.uid matches admins.id
                .single()

            if (error || !data) {
                console.error('Not an admin or error fetching profile:', error)
                // If not admin, sign out? Or just deny access?
                // User said "Admins only".
                // We'll let the UI handle the redirect, or do it here.
                setProfile(null)
            } else {
                setProfile(data as AdminProfile)
            }
        } catch (err) {
            console.error('Unexpected error fetching admin profile:', err)
            setProfile(null)
        } finally {
            setIsLoading(false)
        }
    }

    const signOut = async () => {
        await supabase.auth.signOut()
        router.push('/login')
    }

    const value = {
        user,
        session,
        profile,
        isLoading,
        isAdmin: !!profile,
        signOut,
    }

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
