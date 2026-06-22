import { Link } from 'react-router-dom'

// Sidebar contains navigation links
export default function Sidebar() {
  return (
    <aside className="fixed left-0 top-16 w-64 h-[calc(100vh-4rem)] bg-gray-900 border-r border-gray-700 p-4 overflow-y-auto">
      <nav className="flex flex-col gap-1">
        <Link to="/" className="px-4 py-2 rounded-lg hover:bg-gray-800">Home</Link>
        <Link to="/trending" className="px-4 py-2 rounded-lg hover:bg-gray-800">Trending</Link>
        <Link to="/feed/subscriptions" className="px-4 py-2 rounded-lg hover:bg-gray-800">Subscriptions</Link>
        <hr className="border-gray-700 my-2" />
        <Link to="/channel/me" className="px-4 py-2 rounded-lg hover:bg-gray-800">Your Channel</Link>
      </nav>
    </aside>
  )
}
