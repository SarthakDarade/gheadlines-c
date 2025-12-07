import DashboardLayout from '@/layouts/DashboardLayout'
import LiveUpdatesModule from '@/modules/live-updates'

export default function LivePage() {
    return (
        <DashboardLayout title="Live Updates">
            <div className="mb-8">
                <h1 className="text-3xl font-bold">Live Control Center</h1>
                <p className="text-gray-500 mt-1">Manage real-time news ticker and flash updates.</p>
            </div>
            <LiveUpdatesModule />
        </DashboardLayout>
    )
}
