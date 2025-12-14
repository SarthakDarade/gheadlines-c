import DashboardLayout from '@/layouts/DashboardLayout'
import SourcesModule from '@/modules/sources'

export default function SourcesPage() {
    return (
        <DashboardLayout title="News Sources">
            <div className="mb-8">
                <h1 className="text-3xl font-bold">Manage Sources</h1>
                <p className="text-gray-500 mt-1">Add, update, and remove news sources from the database.</p>
            </div>
            <SourcesModule />
        </DashboardLayout>
    )
}
