import DashboardLayout from '@/layouts/DashboardLayout'
import { useState } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { useAuth } from '@/hooks/useAuth'
import { User, Lock, Save, Globe, Bell } from 'lucide-react'
import { toast } from 'react-hot-toast'

export default function Settings() {
    const { user, signOut } = useAuth()
    const [loading, setLoading] = useState(false)

    // Placeholder mock settings until specific requirements
    const [settings, setSettings] = useState({
        siteName: 'TroyGH News',
        adminEmail: user?.email || '',
        notifications: true,
        maintenanceMode: false
    })

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        // Simulate save
        await new Promise(r => setTimeout(r, 1000))
        toast.success("Settings saved successfully")
        setLoading(false)
    }

    const handlePasswordReset = async () => {
        if (!user?.email) return
        const { error } = await supabase.auth.resetPasswordForEmail(user.email, {
            redirectTo: window.location.origin + '/reset-password',
        })
        if (error) toast.error(error.message)
        else toast.success("Password reset email sent!")
    }

    return (
        <DashboardLayout title="Settings">
            <div className="mb-8">
                <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-gray-700 to-black dark:from-white dark:to-gray-400">
                    System Settings
                </h1>
                <p className="text-gray-500 mt-1">Configure global application preferences.</p>
            </div>

            <div className="max-w-4xl">
                <form onSubmit={handleSave} className="space-y-6">

                    {/* General Settings */}
                    <div className="glass-card p-6 border border-white/10">
                        <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                            <Globe size={20} className="text-blue-500" />
                            General Configuration
                        </h3>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <label className="block text-sm font-medium mb-1">Site App Name</label>
                                <input
                                    className="glass-input"
                                    value={settings.siteName}
                                    onChange={e => setSettings({ ...settings, siteName: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium mb-1">Admin Contact Email</label>
                                <input
                                    className="glass-input bg-gray-100 dark:bg-white/5 opacity-70"
                                    value={settings.adminEmail}
                                    disabled
                                />
                            </div>
                        </div>
                    </div>

                    {/* Security & Account */}
                    <div className="glass-card p-6 border border-white/10">
                        <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                            <Lock size={20} className="text-purple-500" />
                            Security & Account
                        </h3>
                        <div className="flex flex-col md:flex-row gap-4 items-center justify-between p-4 bg-gray-50 dark:bg-white/5 rounded-xl border border-dashed border-gray-300 dark:border-gray-700">
                            <div>
                                <h4 className="font-medium">Password Reset</h4>
                                <p className="text-sm text-gray-500">Send a password reset link to your email.</p>
                            </div>
                            <button
                                type="button"
                                onClick={handlePasswordReset}
                                className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm hover:bg-gray-100 dark:hover:bg-white/10 transition-colors"
                            >
                                Send Reset Link
                            </button>
                        </div>

                        <div className="mt-4 flex flex-col md:flex-row gap-4 items-center justify-between p-4 bg-red-50 dark:bg-red-900/10 rounded-xl border border-dashed border-red-200 dark:border-red-800">
                            <div>
                                <h4 className="font-medium text-red-600 dark:text-red-400">Logout Session</h4>
                                <p className="text-sm text-gray-500">End your current admin session securely.</p>
                            </div>
                            <button
                                type="button"
                                onClick={signOut}
                                className="px-4 py-2 bg-red-500 text-white rounded-lg text-sm hover:bg-red-600 transition-colors shadow-sm"
                            >
                                Sign Out
                            </button>
                        </div>
                    </div>

                    {/* Toggles */}
                    <div className="glass-card p-6 border border-white/10">
                        <h3 className="text-lg font-bold mb-4 flex items-center gap-2">
                            <Bell size={20} className="text-yellow-500" />
                            System Toggles
                        </h3>
                        <div className="space-y-4">
                            <label className="flex items-center justify-between cursor-pointer">
                                <span className="font-medium text-gray-700 dark:text-gray-300">Enable Maintenance Mode</span>
                                <div className="relative inline-flex items-center cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={settings.maintenanceMode}
                                        onChange={e => setSettings({ ...settings, maintenanceMode: e.target.checked })}
                                        className="sr-only peer"
                                    />
                                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
                                </div>
                            </label>
                            <label className="flex items-center justify-between cursor-pointer">
                                <span className="font-medium text-gray-700 dark:text-gray-300">Email Notifications for New Users</span>
                                <div className="relative inline-flex items-center cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={settings.notifications}
                                        onChange={e => setSettings({ ...settings, notifications: e.target.checked })}
                                        className="sr-only peer"
                                    />
                                    <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-600"></div>
                                </div>
                            </label>
                        </div>
                    </div>

                    <div className="flex justify-end pt-4">
                        <button
                            type="submit"
                            disabled={loading}
                            className="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-70 flex items-center gap-2 font-bold"
                        >
                            <Save size={20} />
                            {loading ? 'Saving Changes...' : 'Save Configuration'}
                        </button>
                    </div>

                </form>
            </div>
        </DashboardLayout>
    )
}
