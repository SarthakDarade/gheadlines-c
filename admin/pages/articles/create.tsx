import DashboardLayout from '@/layouts/DashboardLayout'
import { useState } from 'react'
import { useRouter } from 'next/router'
import { motion } from 'framer-motion'
import { Save, ArrowLeft, UploadCloud, Link as LinkIcon, Trash2 } from 'lucide-react'
import TiptapEditor from '@/components/TiptapEditor'
import { supabase } from '@/lib/supabaseClient'
import Link from 'next/link'
import slugify from 'slugify'
import { Toaster, toast } from 'react-hot-toast'
import { useAuth } from '@/hooks/useAuth'

export default function CreateArticle() {
    const router = useRouter()
    const { user } = useAuth()
    const [loading, setLoading] = useState(false)
    const [formData, setFormData] = useState({
        title: '',
        slug: '',
        description: '',
        content: '<p>Start writing...</p>',
        category_id: '',
        tags: '',
        published: false,
        meta_title: '',
        meta_description: '',
        image_url: ''
    })

    // Categories
    const categories = [
        { id: 'tech', name: 'Technology' },
        { id: 'politics', name: 'Politics' },
        { id: 'business', name: 'Business' },
        { id: 'sports', name: 'Sports' },
        { id: 'entertainment', name: 'Entertainment' },
    ]

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        const { name, value } = e.target
        setFormData(prev => {
            const updates = { ...prev, [name]: value }
            if (name === 'title') {
                updates.slug = slugify(value, { lower: true, strict: true })
            }
            return updates
        })
    }

    const handleEditorChange = (html: string) => {
        setFormData(prev => ({ ...prev, content: html }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!user) {
            toast.error("You must be logged in to publish.")
            return
        }
        setLoading(true)

        try {
            const { error } = await supabase.from('articles').insert([{
                ...formData,
                author_id: user.id, // Important: Link to author
                tags: formData.tags.split(',').map(t => t.trim()).filter(Boolean),
                created_at: new Date().toISOString()
            }])

            if (error) throw error

            toast.success("Article published successfully!")
            router.push('/articles')
        } catch (err: any) {
            console.error('Error creating article:', err)
            toast.error('Error creating article: ' + err.message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <DashboardLayout title="Create Article">
            <form onSubmit={handleSubmit} className="max-w-5xl mx-auto pb-20">
                <div className="flex items-center justify-between mb-8">
                    <div className="flex items-center gap-4">
                        <Link href="/articles">
                            <button type="button" className="p-2 hover:bg-gray-100 dark:hover:bg-white/10 rounded-lg transition-colors">
                                <ArrowLeft size={24} />
                            </button>
                        </Link>
                        <div>
                            <h1 className="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">New Article</h1>
                            <p className="text-gray-500 text-sm">Create and publish a new story</p>
                        </div>
                    </div>
                    <div className="flex gap-3">
                        <button
                            type="submit"
                            disabled={loading}
                            className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg flex items-center gap-2 font-medium shadow-lg shadow-blue-600/20"
                        >
                            <Save size={18} />
                            {loading ? 'Publishing...' : 'Publish Article'}
                        </button>
                    </div>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                    <div className="lg:col-span-2 space-y-6">
                        <div className="glass-card p-6 border border-white/10 space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Title</label>
                                <input
                                    type="text"
                                    name="title"
                                    value={formData.title}
                                    onChange={handleChange}
                                    className="glass-input"
                                    placeholder="Enter article headline..."
                                    required
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Slug</label>
                                <input
                                    type="text"
                                    name="slug"
                                    value={formData.slug}
                                    onChange={handleChange}
                                    className="glass-input bg-gray-50 dark:bg-white/5 text-gray-500"
                                    placeholder="article-slug"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Summary</label>
                                <textarea
                                    name="description"
                                    value={formData.description}
                                    onChange={handleChange}
                                    rows={3}
                                    className="glass-input"
                                    placeholder="Brief summary for cards and SEO..."
                                />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Content</label>
                            <TiptapEditor content={formData.content} onChange={handleEditorChange} />
                        </div>

                        <div className="glass-card p-6 border border-white/10 space-y-4">
                            <h3 className="font-medium flex items-center gap-2 text-gray-900 dark:text-gray-100">
                                <UploadCloud size={18} />
                                SEO Settings
                            </h3>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meta Title</label>
                                <input
                                    type="text"
                                    name="meta_title"
                                    value={formData.meta_title}
                                    onChange={handleChange}
                                    className="glass-input"
                                    placeholder="SEO Title"
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Meta Keywords</label>
                                <input
                                    type="text"
                                    name="tags"
                                    value={formData.tags}
                                    onChange={handleChange}
                                    className="glass-input"
                                    placeholder="news, world, politics (comma separated)"
                                />
                            </div>
                        </div>
                    </div>

                    <div className="space-y-6">
                        <div className="glass-card p-6 border border-white/10 space-y-4">
                            <h3 className="font-bold text-gray-900 dark:text-gray-100">Publishing</h3>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label>
                                <select
                                    name="category_id"
                                    value={formData.category_id}
                                    onChange={handleChange}
                                    className="glass-input"
                                >
                                    <option value="">Select Category</option>
                                    {categories.map(c => (
                                        <option key={c.id} value={c.id}>{c.name}</option>
                                    ))}
                                </select>
                            </div>
                            <div>
                                <label className="flex items-center gap-2 cursor-pointer">
                                    <input
                                        type="checkbox"
                                        checked={formData.published}
                                        onChange={e => setFormData(prev => ({ ...prev, published: e.target.checked }))}
                                        className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                                    />
                                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Publish Immediately</span>
                                </label>
                            </div>
                        </div>

                        <div className="glass-card p-6 border border-white/10">
                            <h3 className="font-bold mb-4 text-gray-900 dark:text-gray-100">Featured Image</h3>
                            {formData.image_url ? (
                                <div className="relative group">
                                    <img src={formData.image_url} alt="Featured" className="w-full h-40 object-cover rounded-lg" />
                                    <button
                                        type="button"
                                        onClick={() => setFormData(p => ({ ...p, image_url: '' }))}
                                        className="absolute top-2 right-2 bg-red-500 text-white p-1 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                                    >
                                        <Trash2 size={14} />
                                    </button>
                                </div>
                            ) : (
                                <div className="border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-xl p-8 text-center hover:bg-gray-50 dark:hover:bg-white/5 transition-colors cursor-pointer"
                                    onClick={() => {
                                        const url = prompt("Enter Image URL");
                                        if (url) setFormData(p => ({ ...p, image_url: url }))
                                    }}
                                >
                                    <UploadCloud className="mx-auto text-gray-400 mb-3" size={32} />
                                    <p className="text-sm text-gray-500">Click to add URL</p>
                                </div>
                            )}
                            <div className="mt-4 flex items-center gap-2">
                                <LinkIcon size={14} className="text-gray-400" />
                                <input
                                    type="url"
                                    className="glass-input text-xs py-1"
                                    placeholder="Or paste URL here..."
                                    value={formData.image_url}
                                    onChange={e => setFormData({ ...formData, image_url: e.target.value })}
                                />
                            </div>
                        </div>
                    </div>
                </div>
            </form>
        </DashboardLayout>
    )
}
