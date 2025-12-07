import DashboardLayout from '@/layouts/DashboardLayout'
import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { Briefcase, Download, Mail, Phone, Calendar, Search } from 'lucide-react'
import { toast } from 'react-hot-toast'

interface Application {
    id: string
    name: string
    email: string
    phone: string | null
    resume_url: string | null
    message: string | null
    created_at: string
}

export default function CareersPage() {
    const [applications, setApplications] = useState<Application[]>([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    useEffect(() => {
        fetchApplications()
    }, [])

    const fetchApplications = async () => {
        const { data, error } = await supabase
            .from('careers_applications')
            .select('*')
            .order('created_at', { ascending: false })

        if (error) {
            console.error(error)
            toast.error("Could not load applications")
        } else {
            setApplications(data || [])
        }
        setLoading(false)
    }

    const filteredApps = applications.filter(app =>
        app.name.toLowerCase().includes(search.toLowerCase()) ||
        app.email.toLowerCase().includes(search.toLowerCase())
    )

    return (
        <DashboardLayout title="Careers">
            <div className="mb-8">
                <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-green-500 to-teal-500">
                    Careers Management
                </h1>
                <p className="text-gray-500 mt-1">Review job applications and resumes.</p>
            </div>

            <div className="glass-card overflow-hidden border border-white/10">
                <div className="p-4 border-b border-white/10 bg-gray-50/50 dark:bg-white/5 flex justify-between items-center">
                    <div className="relative max-w-sm w-full">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                        <input
                            type="text"
                            placeholder="Search applicant..."
                            className="glass-input pl-10 h-10"
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                        />
                    </div>
                    <div className="text-sm text-gray-500">
                        Total: {applications.length}
                    </div>
                </div>

                {loading ? (
                    <div className="p-10 text-center text-gray-500">Loading applications...</div>
                ) : filteredApps.length === 0 ? (
                    <div className="p-10 text-center flex flex-col items-center text-gray-400">
                        <Briefcase size={48} className="mb-3 opacity-20" />
                        <p>No applications found.</p>
                    </div>
                ) : (
                    <div className="divide-y divide-gray-100 dark:divide-white/5">
                        {filteredApps.map(app => (
                            <div key={app.id} className="p-6 hover:bg-white/50 dark:hover:bg-white/5 transition-colors group">
                                <div className="flex flex-col md:flex-row justify-between gap-4">
                                    <div className="flex-1">
                                        <div className="flex items-center gap-2 mb-1">
                                            <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100">{app.name}</h3>
                                            <span className="text-xs bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 px-2 py-0.5 rounded-full">New</span>
                                        </div>

                                        <div className="flex flex-wrap gap-4 text-sm text-gray-500 mt-2">
                                            <div className="flex items-center gap-1.5">
                                                <Mail size={14} />
                                                <a href={`mailto:${app.email}`} className="hover:text-blue-500">{app.email}</a>
                                            </div>
                                            {app.phone && (
                                                <div className="flex items-center gap-1.5">
                                                    <Phone size={14} />
                                                    <span>{app.phone}</span>
                                                </div>
                                            )}
                                            <div className="flex items-center gap-1.5">
                                                <Calendar size={14} />
                                                <span>{new Date(app.created_at).toLocaleDateString()}</span>
                                            </div>
                                        </div>

                                        {app.message && (
                                            <p className="mt-3 text-sm text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-black/20 p-3 rounded-lg border border-gray-100 dark:border-white/5">
                                                "{app.message}"
                                            </p>
                                        )}
                                    </div>

                                    <div className="flex items-start">
                                        {app.resume_url ? (
                                            <a
                                                href={app.resume_url}
                                                target="_blank"
                                                rel="noreferrer"
                                                className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors shadow-sm"
                                            >
                                                <Download size={16} />
                                                Download Resume
                                            </a>
                                        ) : (
                                            <span className="text-xs text-gray-400 italic px-4 py-2">No Resume</span>
                                        )}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </DashboardLayout>
    )
}
