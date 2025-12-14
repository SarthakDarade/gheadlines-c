import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { supabase } from '@/lib/supabaseClient'
import { Globe, Trash2, Plus, Loader2, Link as LinkIcon, Image as ImageIcon, Edit2, X, Save, RefreshCw } from 'lucide-react'
import { toast } from 'react-hot-toast'

interface Source {
    id: string
    channel_name: string
    rss_link: string
    url?: string
    category?: string
    created_at: string
}

export default function SourcesModule() {
    const [sources, setSources] = useState<Source[]>([])
    const [loading, setLoading] = useState(true)
    const [submitting, setSubmitting] = useState(false)

    // Form State
    const [name, setName] = useState('')
    const [rssLink, setRssLink] = useState('')
    const [url, setUrl] = useState('')
    const [category, setCategory] = useState('')
    const [editingId, setEditingId] = useState<string | null>(null)

    useEffect(() => {
        fetchSources()

        // Realtime subscription as a backup/sync mechanism
        const subscription = supabase
            .channel('sources_changes')
            .on('postgres_changes', { event: '*', schema: 'public', table: 'sources' }, () => {
                fetchSources()
            })
            .subscribe()

        return () => {
            subscription.unsubscribe()
        }
    }, [])

    const fetchSources = async () => {
        try {
            const { data, error } = await supabase
                .from('sources')
                .select('*')
                .order('created_at', { ascending: false })

            if (error) throw error
            if (data) setSources(data)
        } catch (error) {
            console.error('Error fetching sources:', error)
            // Quietly fail or show toast if needed, but don't spam
        } finally {
            setLoading(false)
        }
    }

    const resetForm = () => {
        setName('')
        setRssLink('')
        setUrl('')
        setCategory('')
        setEditingId(null)
    }

    const handleEdit = (source: Source) => {
        setName(source.channel_name)
        setRssLink(source.rss_link || '')
        setUrl(source.url || '')
        setCategory(source.category || '')
        setEditingId(source.id)
        toast.dismiss()
        toast('Editing ' + source.channel_name, { icon: '✏️' })
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!name) return

        setSubmitting(true)
        const toastId = toast.loading(editingId ? 'Updating source...' : 'Adding source...')

        try {
            if (editingId) {
                const { error } = await supabase
                    .from('sources')
                    .update({
                        channel_name: name,
                        rss_link: rssLink,
                        url: url || null,
                        category: category || 'General',
                        updated_at: new Date().toISOString()
                    })
                    .eq('id', editingId)

                if (error) throw error

                toast.success('Source updated successfully!', { id: toastId })
            } else {
                const { error } = await supabase
                    .from('sources')
                    .insert([{
                        channel_name: name,
                        rss_link: rssLink,
                        url: url || null,
                        category: category || 'General'
                    }])

                if (error) throw error

                toast.success('Source added successfully!', { id: toastId })
            }

            resetForm()
            fetchSources() // Manually trigger fetch to ensure UI is up to date immediately
        } catch (error: any) {
            console.error('Submit error:', error)
            toast.error('Operation failed: ' + error.message, { id: toastId })
        } finally {
            setSubmitting(false)
        }
    }

    const handleDelete = async (id: string) => {
        if (!confirm('Are you sure you want to delete this source?')) return

        const toastId = toast.loading('Deleting source...')
        try {
            const { error } = await supabase.from('sources').delete().eq('id', id)
            if (error) throw error

            toast.success('Source deleted', { id: toastId })
            setSources(sources.filter(s => s.id !== id)) // Optimistic update
            fetchSources() // Sync
        } catch (error: any) {
            toast.error('Deletion failed: ' + error.message, { id: toastId })
        }
    }

    return (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
            {/* Control Panel (Left Side - 4 cols) */}
            <div className="lg:col-span-4 space-y-6">
                <div className="glass p-6 rounded-2xl border border-white/10 sticky top-24">
                    <div className="flex items-center justify-between mb-6">
                        <div className="flex items-center gap-3">
                            <div className={`w-10 h-10 rounded-full flex items-center justify-center transition-colors ${editingId ? 'bg-orange-500/10 text-orange-500' : 'bg-blue-500/10 text-blue-500'}`}>
                                {editingId ? <Edit2 size={20} /> : <Plus size={20} />}
                            </div>
                            <div>
                                <h2 className="text-lg font-bold">{editingId ? 'Edit Source' : 'Add Source'}</h2>
                                <p className="text-xs text-gray-500">{editingId ? 'Update details' : 'New provider'}</p>
                            </div>
                        </div>
                        {editingId && (
                            <button onClick={resetForm} className="p-2 hover:bg-gray-100 dark:hover:bg-white/5 rounded-lg text-gray-500 transition-colors" title="Cancel Edit">
                                <X size={20} />
                            </button>
                        )}
                    </div>

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-1">
                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider ml-1">Source Name / Channel Name</label>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="glass-input"
                                placeholder="e.g. The Verge, BBC News"
                                required
                            />
                        </div>

                        <div className="space-y-1">
                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider ml-1">RSS Link (Required)</label>
                            <div className="relative">
                                <RefreshCw size={14} className="absolute left-3 top-3.5 text-gray-400" />
                                <input
                                    type="url"
                                    value={rssLink}
                                    onChange={(e) => setRssLink(e.target.value)}
                                    className="glass-input pl-9"
                                    placeholder="https://site.com/rss.xml"
                                    required
                                />
                            </div>
                        </div>

                        <div className="space-y-1">
                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider ml-1">Website URL</label>
                            <div className="relative">
                                <LinkIcon size={14} className="absolute left-3 top-3.5 text-gray-400" />
                                <input
                                    type="url"
                                    value={url}
                                    onChange={(e) => setUrl(e.target.value)}
                                    className="glass-input pl-9"
                                    placeholder="https://"
                                />
                            </div>
                        </div>

                        <div className="space-y-1">
                            <label className="text-[10px] font-bold text-gray-400 uppercase tracking-wider ml-1">Category</label>
                            <div className="relative">
                                <Globe size={14} className="absolute left-3 top-3.5 text-gray-400" />
                                <input
                                    type="text"
                                    value={category}
                                    onChange={(e) => setCategory(e.target.value)}
                                    className="glass-input pl-9"
                                    placeholder="e.g. World, Tech, Politics"
                                />
                            </div>
                        </div>

                        <button
                            type="submit"
                            disabled={submitting}
                            className={`w-full py-3 text-white rounded-xl font-bold flex items-center justify-center gap-2 transition-all shadow-lg active:scale-95 ${submitting ? 'opacity-70 cursor-not-allowed' : ''
                                } ${editingId
                                    ? 'bg-orange-600 hover:bg-orange-700 shadow-orange-600/20'
                                    : 'bg-blue-600 hover:bg-blue-700 shadow-blue-600/20'
                                }`}
                        >
                            {submitting ? (
                                <Loader2 size={18} className="animate-spin" />
                            ) : (
                                <>
                                    {editingId ? <Save size={18} /> : <Plus size={18} />}
                                    {editingId ? 'Update Source' : 'Add Source'}
                                </>
                            )}
                        </button>
                    </form>
                </div>
            </div>

            {/* Sources List (Right Side - 8 cols) */}
            <div className="lg:col-span-8">
                <div className="flex items-center justify-between mb-6">
                    <h2 className="text-xl font-bold flex items-center gap-2">
                        <Globe className="text-blue-500" size={24} />
                        Managed Sources
                    </h2>
                    <div className="flex items-center gap-3">
                        <button
                            onClick={() => { setLoading(true); fetchSources(); }}
                            className="p-2 hover:bg-white/10 rounded-lg text-gray-500 transition-colors"
                            title="Refresh List"
                        >
                            <RefreshCw size={18} className={loading ? 'animate-spin' : ''} />
                        </button>
                        <span className="text-sm font-mono bg-white/10 px-3 py-1 rounded-full border border-white/5">
                            {sources.length} sources
                        </span>
                    </div>
                </div>

                <div className="space-y-3">
                    <AnimatePresence mode='popLayout'>
                        {loading && sources.length === 0 ? (
                            <div className="flex flex-col items-center justify-center py-20 text-gray-400 animate-pulse">
                                <Loader2 className="animate-spin mb-2" size={32} />
                                <p>Syncing with database...</p>
                            </div>
                        ) : sources.length === 0 ? (
                            <motion.div
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                className="text-center py-20 text-gray-400 glass rounded-xl border-dashed"
                            >
                                <Globe size={48} className="mx-auto mb-4 opacity-20" />
                                <p>No sources found in the database.</p>
                                <p className="text-xs opacity-60 mt-1">Add your first source on the left.</p>
                            </motion.div>
                        ) : (
                            sources.map((source) => (
                                <motion.div
                                    key={source.id}
                                    layout
                                    initial={{ opacity: 0, scale: 0.95 }}
                                    animate={{ opacity: 1, scale: 1 }}
                                    exit={{ opacity: 0, scale: 0.95 }}
                                    className="glass p-4 rounded-xl border border-white/5 flex gap-4 group hover:bg-white/40 dark:hover:bg-white/5 transition-all hover:shadow-lg hover:-translate-y-0.5"
                                >
                                    <div className="w-14 h-14 rounded-xl bg-orange-500/10 dark:bg-orange-500/5 border border-orange-500/20 flex items-center justify-center overflow-hidden shrink-0 flex-col gap-1">
                                        <Globe size={18} className="text-orange-500" />
                                        <span className="text-[9px] font-bold text-orange-500 uppercase tracking-wider w-full text-center px-1 truncate">
                                            {source.category || 'News'}
                                        </span>
                                    </div>

                                    <div className="flex-1 min-w-0 flex flex-col justify-center">
                                        <div className="flex items-center justify-between">
                                            <h3 className="font-bold text-gray-900 dark:text-white truncate text-lg">
                                                {source.channel_name}
                                            </h3>
                                            <span className="text-[10px] font-mono text-gray-400 opacity-60">
                                                {new Date(source.created_at).toLocaleDateString()}
                                            </span>
                                        </div>

                                        {source.url && (
                                            <a
                                                href={source.url}
                                                target="_blank"
                                                rel="noreferrer"
                                                className="text-xs text-blue-500 hover:text-blue-400 hover:underline flex items-center gap-1 truncate mt-0.5 w-fit font-medium"
                                            >
                                                <LinkIcon size={12} />
                                                {source.url.replace(/^https?:\/\//, '').replace(/\/$/, '')}
                                            </a>
                                        )}
                                    </div>

                                    <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-all translate-x-4 group-hover:translate-x-0">
                                        <button
                                            onClick={() => handleEdit(source)}
                                            className="p-2 text-gray-400 hover:text-orange-500 hover:bg-orange-500/10 rounded-lg transition-all"
                                            title="Edit source"
                                        >
                                            <Edit2 size={18} />
                                        </button>
                                        <button
                                            onClick={() => handleDelete(source.id)}
                                            className="p-2 text-gray-400 hover:text-red-500 hover:bg-red-500/10 rounded-lg transition-all"
                                            title="Delete source"
                                        >
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </motion.div>
                            ))
                        )}
                    </AnimatePresence>
                </div>
            </div>
        </div>
    )
}
