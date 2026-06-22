import { useSelector } from 'react-redux'
import { Link } from 'react-router-dom'
import UploadForm from '../components/UploadForm'

// Upload page — auth-gated upload form
export default function Upload() {
  const { user } = useSelector((s) => s.auth)

  if (!user) {
    return (
      <div className="max-w-2xl mx-auto text-center py-16">
        <h1 className="text-2xl font-bold mb-4">Sign in to upload</h1>
        <p className="text-gray-400 mb-6">You need an account to upload videos</p>
        <Link to="/" className="bg-white text-black px-6 py-2 rounded-full text-sm font-medium">
          Go home
        </Link>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Upload video</h1>
      <UploadForm />
    </div>
  )
}
