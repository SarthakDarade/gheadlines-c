import DashboardLayout from '@/layouts/DashboardLayout'
import { useState, useEffect } from 'react'
import { supabase } from '@/lib/supabaseClient'
import { Mail, Trash2, Search, MessageSquare, Clock } from 'lucide-react'
import { toast } from 'react-hot-toast'

interface ContactMessage {
    id: string
    name: string
    email: string
    subject: string | null
    message: string
    created_at: string
    // read: boolean // Assuming field exists, if not we'll ignore or add locally
}

export default function CommunicationsPage() {
    const [messages, setMessages] = useState<ContactMessage[]>([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    useEffect(() => {
        fetchMessages()
    }, [])

    const fetchMessages = async () => {
        const { data, error } = await supabase
            .from('contact_messages')
            .select('*')
            .order('created_at', { ascending: false })

        if (error) {
            console.error(error)
            toast.error("Could not load messages")
        } else {
            setMessages(data || [])
        }
        setLoading(false)
    }

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure you want to delete this message?")) return;

        const { error } = await supabase.from('contact_messages').delete().eq('id', id)
        if (error) {
            toast.error("Failed to delete")
        } else {
            setMessages(prev => prev.filter(m => m.id !== id))
            toast.success("Message deleted")
        }
    }

    const filteredMessages = messages.filter(msg =>
        msg.name.toLowerCase().includes(search.toLowerCase()) ||
        msg.email.toLowerCase().includes(search.toLowerCase()) ||
        (msg.subject && msg.subject.toLowerCase().includes(search.toLowerCase()))
    )

    return (
        <DashboardLayout title="Communications">
            <div className="mb-8">
                <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-pink-500 to-rose-500">
                    Communications
                </h1>
                <p className="text-gray-500 mt-1">Manage contact form inquiries and messages.</p>
            </div>

            <div className="glass-card overflow-hidden border border-white/10">
                <div className="p-4 border-b border-white/10 bg-gray-50/50 dark:bg-white/5 flex justify-between items-center">
                    <div className="relative max-w-sm w-full">
                        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                        <input
                            type="text"
                            placeholder="Search sender..."
                            className="glass-input pl-10 h-10"
                            value={search}
                            onChange={e => setSearch(e.target.value)}
                        />
                    </div>
                    <div className="text-sm text-gray-500">
                        {messages.length} Messages
                    </div>
                </div>

                {loading ? (
                    <div className="p-10 text-center text-gray-500">Loading messages...</div>
                ) : filteredMessages.length === 0 ? (
                    <div className="p-10 text-center flex flex-col items-center text-gray-400">
                        <Mail size={48} className="mb-3 opacity-20" />
                        <p>No messages found.</p>
                    </div>
                ) : (
                    <div className="divide-y divide-gray-100 dark:divide-white/5">
                        {filteredMessages.map(msg => (
                            <div key={msg.id} className="p-6 hover:bg-white/50 dark:hover:bg-white/5 transition-colors group">
                                <div className="flex flex-col sm:flex-row gap-4 justify-between items-start">
                                    <div className="flex-1 space-y-2">
                                        <div className="flex items-center gap-2">
                                            <h3 className="font-bold text-gray-900 dark:text-gray-100">{msg.name}</h3>
                                            <span className="text-xs text-gray-400">&lt;{msg.email}&gt;</span>
                                        </div>
                                        <div className="text-sm font-medium text-gray-700 dark:text-gray-300">
                                            {msg.subject || 'No Subject'}
                                        </div>
                                        <p className="text-gray-600 dark:text-gray-400 text-sm leading-relaxed bg-gray-50/50 dark:bg-white/5 p-3 rounded-lg border border-gray-100 dark:border-white/5">
                                            {msg.message}
                                        </p>
                                        <div className="flex items-center gap-2 text-xs text-gray-400 mt-2">
                                            <Clock size={12} />
                                            {new Date(msg.created_at).toLocaleString()}
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity self-start sm:self-center">
                                        <a href={`mailto:${msg.email}`} className="p-2 text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-500/10 rounded-lg" title="Reply">
                                            <Mail size={18} />
                                        </a>
                                        <button
                                            onClick={() => handleDelete(msg.id)}
                                            className="p-2 text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg"
                                            title="Delete"
                                        >
                                            <Trash2 size={18} />
                                        </button>
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
