import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Image from '@tiptap/extension-image'
import { Bold, Italic, Strikethrough, Heading1, Heading2, List, ListOrdered, Quote, Undo, Redo, Image as ImageIcon } from 'lucide-react'

interface TiptapEditorProps {
    content: string
    onChange: (html: string) => void
}

const MenuBar = ({ editor }: { editor: any }) => {
    if (!editor) {
        return null
    }

    const buttons = [
        {
            icon: Bold,
            onClick: () => editor.chain().focus().toggleBold().run(),
            isActive: editor.isActive('bold'),
            title: 'Bold'
        },
        {
            icon: Italic,
            onClick: () => editor.chain().focus().toggleItalic().run(),
            isActive: editor.isActive('italic'),
            title: 'Italic'
        },
        {
            icon: Strikethrough,
            onClick: () => editor.chain().focus().toggleStrike().run(),
            isActive: editor.isActive('strike'),
            title: 'Strike'
        },
        {
            icon: Heading1,
            onClick: () => editor.chain().focus().toggleHeading({ level: 1 }).run(),
            isActive: editor.isActive('heading', { level: 1 }),
            title: 'H1'
        },
        {
            icon: Heading2,
            onClick: () => editor.chain().focus().toggleHeading({ level: 2 }).run(),
            isActive: editor.isActive('heading', { level: 2 }),
            title: 'H2'
        },
        {
            icon: List,
            onClick: () => editor.chain().focus().toggleBulletList().run(),
            isActive: editor.isActive('bulletList'),
            title: 'Bullet List'
        },
        {
            icon: ListOrdered,
            onClick: () => editor.chain().focus().toggleOrderedList().run(),
            isActive: editor.isActive('orderedList'),
            title: 'Ordered List'
        },
        {
            icon: Quote,
            onClick: () => editor.chain().focus().toggleBlockquote().run(),
            isActive: editor.isActive('blockquote'),
            title: 'Blockquote'
        },
    ]

    const addImage = () => {
        const url = window.prompt('URL')
        if (url) {
            editor.chain().focus().setImage({ src: url }).run()
        }
    }

    return (
        <div className="flex flex-wrap gap-2 p-3 border-b border-gray-200 dark:border-white/10 bg-gray-50/50 dark:bg-black/20">
            {buttons.map((btn, index) => {
                const Icon = btn.icon
                return (
                    <button
                        key={index}
                        onClick={btn.onClick}
                        className={`p-2 rounded-lg transition-colors ${btn.isActive
                            ? 'bg-blue-100 text-blue-600 dark:bg-blue-500/20 dark:text-blue-400'
                            : 'text-gray-500 hover:bg-gray-200 dark:hover:bg-white/10'
                            }`}
                        title={btn.title}
                        type="button"
                    >
                        <Icon size={18} />
                    </button>
                )
            })}
            <div className="w-[1px] h-8 bg-gray-300 dark:bg-white/10 mx-1 self-center" />
            <button
                onClick={addImage}
                className="p-2 rounded-lg text-gray-500 hover:bg-gray-200 dark:hover:bg-white/10 transition-colors"
                title="Add Image"
                type="button"
            >
                <ImageIcon size={18} />
            </button>
            <div className="ml-auto flex gap-2">
                <button
                    onClick={() => editor.chain().focus().undo().run()}
                    className="p-2 rounded-lg text-gray-500 hover:bg-gray-200 dark:hover:bg-white/10"
                    type="button"
                >
                    <Undo size={18} />
                </button>
                <button
                    onClick={() => editor.chain().focus().redo().run()}
                    className="p-2 rounded-lg text-gray-500 hover:bg-gray-200 dark:hover:bg-white/10"
                    type="button"
                >
                    <Redo size={18} />
                </button>
            </div>
        </div>
    )
}

export default function TiptapEditor({ content, onChange }: TiptapEditorProps) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Image,
        ],
        content,
        editorProps: {
            attributes: {
                class: 'prose dark:prose-invert max-w-none focus:outline-none min-h-[300px] p-4 text-gray-700 dark:text-gray-200',
            },
        },
        onUpdate: ({ editor }) => {
            onChange(editor.getHTML())
        },
        immediatelyRender: false,
    })

    return (
        <div className="glass rounded-xl border border-white/10 overflow-hidden shadow-sm">
            <MenuBar editor={editor} />
            <EditorContent editor={editor} />
        </div>
    )
}
