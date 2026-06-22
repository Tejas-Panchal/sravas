import { useState, useEffect } from 'react'
import VideoCard from '../components/VideoCard'
import api from '../services/api'

// Trending page — fetches and displays trending videos
export default function Trending() {
  const [videos, setVideos] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/analytics/trending')
      .then((res) => setVideos(res.data?.videos || []))
      .catch(() => setVideos([]))
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div>
        <h1 className="text-2xl font-bold mb-6">Trending</h1>
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
      <h1 className="text-2xl font-bold mb-6">Trending</h1>
      {videos.length === 0 ? (
        <p className="text-gray-400 text-center py-10">No trending videos</p>
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
