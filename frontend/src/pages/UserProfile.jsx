import { useState, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import VideoCard from '../components/VideoCard'
import api from '../services/api'

// UserProfile page — channel page showing profile info and uploaded videos
export default function UserProfile() {
  const { id } = useParams()
  const [profile, setProfile] = useState(null)
  const [videos, setVideos] = useState([])
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)
  const [subscribed, setSubscribed] = useState(false)

  useEffect(() => {
    setLoading(true)
    Promise.all([
      api.get(`/users/${id}`).catch(() => ({ data: null })),
      api.get(`/users/${id}/videos`).catch(() => ({ data: [] })),
      api.get(`/users/${id}/stats`).catch(() => ({ data: null })),
    ]).then(([profileRes, videosRes, statsRes]) => {
      setProfile(profileRes.data)
      setVideos(videosRes.data?.videos || videosRes.data || [])
      setStats(statsRes.data)
      setSubscribed(profileRes.data?.subscribed || false)
    }).finally(() => setLoading(false))
  }, [id])

  const handleSubscribe = async () => {
    try {
      await api.post(`/users/${id}/subscribe`)
      setSubscribed(!subscribed)
    } catch {}
  }

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="flex items-center gap-4 mb-6">
          <div className="w-16 h-16 bg-gray-800 rounded-full" />
          <div>
            <div className="h-6 bg-gray-800 rounded w-40 mb-2" />
            <div className="h-4 bg-gray-800 rounded w-24" />
          </div>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-gray-800 rounded-lg">
              <div className="aspect-video bg-gray-700 rounded-t-lg" />
              <div className="p-3"><div className="h-4 bg-gray-700 rounded w-3/4" /></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (!profile) {
    return <div className="text-gray-400 text-center py-10">Channel not found</div>
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <div className="w-16 h-16 bg-gray-600 rounded-full overflow-hidden">
          {profile.avatar && (
            <img src={profile.avatar} alt="" className="w-full h-full object-cover" />
          )}
        </div>
        <div>
          <h1 className="text-2xl font-bold">{profile.username}</h1>
          <p className="text-gray-400 text-sm">
            {stats?.subscriber_count?.toLocaleString() || 0} subscribers
          </p>
        </div>
        <button
          onClick={handleSubscribe}
          className={`ml-4 px-4 py-1.5 rounded-full text-sm font-medium ${
            subscribed ? 'bg-gray-700 text-white' : 'bg-white text-black'
          }`}
        >
          {subscribed ? 'Subscribed' : 'Subscribe'}
        </button>
      </div>

      {videos.length === 0 ? (
        <p className="text-gray-400 text-center py-10">No videos yet</p>
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
