import DashboardLayout from '@/layouts/DashboardLayout'
import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { motion } from 'framer-motion'
import { Search, Shield, Trash2, Ban, CheckCircle, MoreVertical, Mail } from 'lucide-react'
import { Toaster, toast } from 'react-hot-toast'

interface Profile {
    id: string
    full_name: string | null
    email?: string // Not always available in profiles, might need join with auth or just rely on what we have. 
    // RLS might limit viewing 'auth.users' directly, but we can view 'public.profiles'.
    // Ideally, profiles table should have email synced or we just show name.
    // Wait, the user asked for "User Management". Usually requires Admin API to delete users from Auth.
    // Viewing public.profiles is easy. Banning requires logic.
    website: string | null
    created_at: string
}

export default function UsersPage() {
    const [users, setUsers] = useState<Profile[]>([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    useEffect(() => {
        fetchUsers()
    }, [])

    const fetchUsers = async () => {
        // Fetch profiles.
        // NOTE: real emails are in auth.users, which is protected. 
        // If you haven't synced emails to public.profiles, you won't see them here unless you use a Postgres Function or Edge Function with Service Role.
        // For this UI demo, I'll attempt to fetch profiles. If email missing, I'll show placeholder or ID.
        const { data, error } = await supabase
            .from('profiles')
            .select('*')
            .order('created_at', { ascending: false })
            .limit(50)

        if (error) {
            toast.error('Failed to load users')
        } else {
            setUsers(data || [])
        }
        setLoading(false)
    }

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this user profile? This does NOT delete their login account (requires Admin API).')) return;

        // We can delete the profile from public table
        const { error } = await supabase.from('profiles').delete().eq('id', id)
        if (error) {
            toast.error('Error deleting profile')
        } else {
            toast.success('Profile deleted')
            setUsers(users.filter(u => u.id !== id))
        }
    }

    return (
        <DashboardLayout title="Users">
            <div className="mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-cyan-500">
                        User Management
                    </h1>
                    <p className="text-gray-500 mt-1">View and manage registered users.</p>
                </div>

                <div className="bg-white/10 p-1 rounded-xl flex gap-1">
                    <span className="px-3 py-1.5 rounded-lg bg-white shadow-sm text-xs font-bold text-gray-800">All Users</span>
                    <span className="px-3 py-1.5 rounded-lg text-xs font-bold text-gray-500 hover:bg-white/5 cursor-pointer">Admins</span>
                    <span className="px-3 py-1.5 rounded-lg text-xs font-bold text-gray-500 hover:bg-white/5 cursor-pointer">Banned</span>
                </div>
            </div>

            <div className="glass-card overflow-hidden">
                {/* Toolbar */}
                <div className="p-4 border-b border-white/10 flex gap-4 items-center bg-white/30 dark:bg-black/20">
                    <div className="relative flex-1 max-w-md">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                        <input
                            type="text"
                            placeholder="Search users..."
                            className="glass-input pl-10"
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                        />
                    </div>
                    <div className="flex-1"></div>
                    {/* Action Buttons */}
                </div>

                {/* Table */}
                <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse">
                        <thead>
                            <tr className="border-b border-white/10 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                                <th className="px-6 py-4 bg-gray-50/50 dark:bg-white/5">User</th>
                                <th className="px-6 py-4 bg-gray-50/50 dark:bg-white/5">Details</th>
                                <th className="px-6 py-4 bg-gray-50/50 dark:bg-white/5">Joined</th>
                                <th className="px-6 py-4 bg-gray-50/50 dark:bg-white/5 text-right">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-white/10">
                            {loading ? (
                                [1, 2, 3].map(i => (
                                    <tr key={i} className="animate-pulse">
                                        <td className="px-6 py-4"><div className="h-10 w-32 bg-gray-200 dark:bg-white/10 rounded"></div></td>
                                        <td className="px-6 py-4"><div className="h-4 w-24 bg-gray-200 dark:bg-white/10 rounded"></div></td>
                                        <td className="px-6 py-4"><div className="h-4 w-20 bg-gray-200 dark:bg-white/10 rounded"></div></td>
                                        <td className="px-6 py-4"></td>
                                    </tr>
                                ))
                            ) : users.length === 0 ? (
                                <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400">No users found.</td></tr>
                            ) : (
                                users.map((user) => (
                                    <motion.tr
                                        key={user.id}
                                        initial={{ opacity: 0 }}
                                        animate={{ opacity: 1 }}
                                        className="group hover:bg-white/40 dark:hover:bg-white/5 transition-colors"
                                    >
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-bold text-sm shadow-md">
                                                    {user.full_name?.[0]?.toUpperCase() || 'U'}
                                                </div>
                                                <div>
                                                    <div className="font-bold text-gray-900 dark:text-gray-100">{user.full_name || 'Unnamed User'}</div>
                                                    <div className="text-xs text-gray-500 font-mono mt-0.5 max-w-[120px] truncate" title={user.id}>ID: {user.id}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-500">
                                            {user.website ? (
                                                <a href={user.website} target="_blank" className="flex items-center gap-1 text-blue-500 hover:underline">
                                                    {user.website.replace(/^https?:\/\//, '')}
                                                </a>
                                            ) : <span className="text-gray-400 italic">No website</span>}
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-500">
                                            {new Date(user.created_at).toLocaleDateString(undefined, { dateStyle: 'medium' })}
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                                                <button className="p-2 hover:bg-blue-500/10 text-blue-500 rounded-lg" title="Email User">
                                                    <Mail size={16} />
                                                </button>
                                                <button
                                                    onClick={() => handleDelete(user.id)}
                                                    className="p-2 hover:bg-red-500/10 text-red-500 rounded-lg"
                                                    title="Delete Profile"
                                                >
                                                    <Trash2 size={16} />
                                                </button>
                                            </div>
                                        </td>
                                    </motion.tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </DashboardLayout>
    )
}
