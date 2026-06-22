import { Outlet } from 'react-router-dom'
import Header from './Header'
import Sidebar from './Sidebar'

// Layout wraps all pages with a shared Header and Sidebar
export default function Layout() {
  return (
    <div className="min-h-screen bg-black text-white">
      <Header />
      <div className="flex pt-16">
        <Sidebar />
        <main className="flex-1 ml-64 p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
