import { useState, useEffect } from 'react'
import { useSelector } from 'react-redux'
import { Link } from 'react-router-dom'
import VideoCard from '../components/VideoCard'
import api from '../services/api'

// Subscriptions page — feed of videos from subscribed channels
export default function Subscriptions() {
  const { user } = useSelector((s) => s.auth)
  const [videos, setVideos] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      setLoading(false)
      return
    }
    api.get(`/users/${user.id}/subscriptions`)
      .then((res) => setVideos(res.data?.videos || []))
      .catch(() => setVideos([]))
      .finally(() => setLoading(false))
  }, [user])

  if (!user) {
    return (
      <div className="text-center py-16">
        <h1 className="text-2xl font-bold mb-4">Subscriptions</h1>
        <p className="text-gray-400 mb-6">Sign in to see your subscriptions</p>
        <Link to="/" className="bg-white text-black px-6 py-2 rounded-full text-sm font-medium">
          Go home
        </Link>
      </div>
    )
  }

  if (loading) {
    return (
      <div>
        <h1 className="text-2xl font-bold mb-6">Subscriptions</h1>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 animate-pulse">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="bg-gray-800 rounded-lg">
              <div className="aspect-video bg-gray-700 rounded-t-lg" />
              <div className="p-3"><div className="h-4 bg-gray-700 rounded w-3/4" /></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Subscriptions</h1>
      {videos.length === 0 ? (
        <p className="text-gray-400 text-center py-10">
          No videos from your subscriptions
        </p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {videos.map((v) => (
            <VideoCard key={v.id} video={v} />
          ))}
        </div>
      )}
    </div>
  )
}
