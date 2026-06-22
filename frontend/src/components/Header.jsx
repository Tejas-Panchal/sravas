import { Link } from 'react-router-dom'
import SearchBar from './SearchBar'

// Header contains logo, search bar, and user menu
export default function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 h-16 bg-gray-900 border-b border-gray-700 flex items-center justify-between px-4 z-50">
      <Link to="/" className="text-xl font-bold">
        Sravas
      </Link>
      <SearchBar />
      <div className="flex items-center gap-4">
        <Link to="/upload" className="text-sm bg-gray-700 px-4 py-2 rounded-full">
          Upload
        </Link>
        <div className="w-8 h-8 bg-gray-600 rounded-full" />
      </div>
    </header>
  )
}
