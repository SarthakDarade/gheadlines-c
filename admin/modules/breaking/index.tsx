'use client';

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { supabase } from '@/lib/supabaseClient'
import { Zap, Plus, X, Edit2, Trash2, ExternalLink, Save } from 'lucide-react'
import { toast } from 'react-hot-toast'

interface BreakingNewsItem {
    id: string
    headline: string
    url?: string
    category?: string
    created_at: string
}

export default function BreakingNewsModule() {
    const [news, setNews] = useState<BreakingNewsItem[]>([])
    const [loading, setLoading] = useState(true)
    const [isAdding, setIsAdding] = useState(false)
    const [editingId, setEditingId] = useState<string | null>(null)

    // Form States
    const [headline, setHeadline] = useState('')
    const [url, setUrl] = useState('')
    const [category, setCategory] = useState('')

    useEffect(() => {
        fetchBreakingNews()
    }, [])

    const fetchBreakingNews = async () => {
        const { data, error } = await supabase
            .from('breaking_news')
            .select('*')
            .order('created_at', { ascending: false })

        if (error) {
            console.error(error)
            toast.error("Failed to load breaking news")
        } else {
            setNews(data || [])
        }
        setLoading(false)
    }

    const resetForm = () => {
        setHeadline('')
        setUrl('')
        setCategory('')
        setEditingId(null)
        setIsAdding(false)
    }

    const startEdit = (item: BreakingNewsItem) => {
        setHeadline(item.headline)
        setUrl(item.url || '')
        setCategory(item.category || '')
        setEditingId(item.id)
        setIsAdding(true) // Reuse the add form UI for editing
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!headline) return toast.error("Headline is required")

        const payload = { headline, url, category }

        try {
            if (editingId) {
                // Update
                const { error } = await supabase
                    .from('breaking_news')
                    .update(payload)
                    .eq('id', editingId)

                if (error) throw error
                setNews(prev => prev.map(item => item.id === editingId ? { ...item, ...payload } : item))
                toast.success("Breaking news updated")
            } else {
                // Create
                const { data, error } = await supabase
                    .from('breaking_news')
                    .insert([payload])
                    .select()
                    .single()

                if (error) throw error
                setNews(prev => [data, ...prev])
                toast.success("Published breaking news")
            }
            resetForm()
        } catch (err: any) {
            toast.error(err.message || "Operation failed")
        }
    }

    const handleDelete = async (id: string) => {
        if (!confirm("Delete this breaking news?")) return

        const { error } = await supabase.from('breaking_news').delete().eq('id', id)
        if (error) {
            toast.error("Failed to delete")
        } else {
            setNews(prev => prev.filter(i => i.id !== id))
            toast.success("Deleted successfully")
        }
    }

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center bg-red-500/10 p-4 rounded-xl border border-red-500/20">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-red-500 rounded-lg animate-pulse">
                        <Zap className="text-white" size={20} />
                    </div>
                    <div>
                        <h2 className="font-bold text-red-600 dark:text-red-400">Breaking News Control</h2>
                        <p className="text-sm text-red-500/80">Urgent alerts displayed at the top of the site</p>
                    </div>
                </div>
                <button
                    onClick={() => { resetForm(); setIsAdding(!isAdding) }}
                    className={`px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2 ${isAdding ? 'bg-gray-200 text-gray-700 hover:bg-gray-300' : 'bg-red-600 text-white hover:bg-red-700'
                        }`}
                >
                    {isAdding ? <><X size={18} /> Cancel</> : <><Plus size={18} /> Add Alert</>}
                </button>
            </div>

            <AnimatePresence>
                {isAdding && (
                    <motion.div
                        initial={{ height: 0, opacity: 0 }}
                        animate={{ height: 'auto', opacity: 1 }}
                        exit={{ height: 0, opacity: 0 }}
                        className="overflow-hidden"
                    >
                        <form onSubmit={handleSubmit} className="glass-card p-6 border border-red-500/30 ring-1 ring-red-500/20 space-y-4">
                            <h3 className="font-bold mb-2">{editingId ? 'Edit Alert' : 'New Breaking Alert'}</h3>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                <div className="md:col-span-2">
                                    <label className="block text-sm font-medium mb-1">Headline</label>
                                    <input
                                        value={headline}
                                        onChange={e => setHeadline(e.target.value)}
                                        className="glass-input border-red-200 focus:border-red-500 focus:ring-red-500/20"
                                        placeholder="URGENT: Market hits all time high..."
                                        autoFocus
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium mb-1">Category (Optional)</label>
                                    <input
                                        value={category}
                                        onChange={e => setCategory(e.target.value)}
                                        className="glass-input"
                                        placeholder="Global, Markets, Politics..."
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium mb-1">Link URL (Optional)</label>
                                    <input
                                        value={url}
                                        onChange={e => setUrl(e.target.value)}
                                        className="glass-input"
                                        placeholder="https://..."
                                    />
                                </div>
                            </div>
                            <div className="flex justify-end pt-2">
                                <button type="submit" className="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium flex items-center gap-2 shadow-lg shadow-red-600/20">
                                    <Save size={18} />
                                    {editingId ? 'Update Alert' : 'Publish Alert'}
                                </button>
                            </div>
                        </form>
                    </motion.div>
                )}
            </AnimatePresence>

            <div className="space-y-3">
                {news.map((item) => (
                    <motion.div
                        key={item.id}
                        layout
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="glass-card p-4 flex flex-col md:flex-row gap-4 items-start md:items-center justify-between group hover:border-red-500/30 transition-colors"
                    >
                        <div className="flex items-start gap-4 flex-1">
                            <div className="mt-1 min-w-[30px] flex justify-center">
                                <Zap size={20} className="text-red-500" />
                            </div>
                            <div>
                                <h4 className="font-bold text-gray-900 dark:text-gray-100 leading-tight">{item.headline}</h4>
                                <div className="flex gap-3 text-xs text-gray-500 mt-1">
                                    <span>{new Date(item.created_at).toLocaleString()}</span>
                                    {item.category && <span className="text-red-500 font-medium uppercase">{item.category}</span>}
                                    {item.url && (
                                        <a href={item.url} target="_blank" className="flex items-center gap-1 hover:text-blue-500 transition-colors">
                                            <ExternalLink size={10} /> Link
                                        </a>
                                    )}
                                </div>
                            </div>
                        </div>

                        <div className="flex items-center gap-2 md:opacity-0 group-hover:opacity-100 transition-opacity self-end md:self-auto">
                            <button
                                onClick={() => startEdit(item)}
                                className="p-2 text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-500/10 rounded-lg transition-colors"
                                title="Edit"
                            >
                                <Edit2 size={18} />
                            </button>
                            <button
                                onClick={() => handleDelete(item.id)}
                                className="p-2 text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10 rounded-lg transition-colors"
                                title="Delete"
                            >
                                <Trash2 size={18} />
                            </button>
                        </div>
                    </motion.div>
                ))}

                {news.length === 0 && !loading && (
                    <div className="text-center py-10 text-gray-500">
                        No active breaking news.
                    </div>
                )}
            </div>
        </div>
    )
}
