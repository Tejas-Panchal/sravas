import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import VideoPlayer from '../components/VideoPlayer'
import CommentList from '../components/CommentList'
import api from '../services/api'

// Watch page — video player with metadata, interactions, and comments
export default function Watch() {
  const { id } = useParams()
  const [video, setVideo] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    (async () => {
      try {
        const res = await api.get(`/videos/${id}`)
        setVideo(res.data)
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load video')
      } finally {
        setLoading(false)
      }
    })()
  }, [id])

  if (loading) {
    return (
      <div className="max-w-4xl animate-pulse">
        <div className="aspect-video bg-gray-800 rounded-lg mb-4" />
        <div className="h-6 bg-gray-800 rounded w-3/4 mb-2" />
        <div className="h-4 bg-gray-800 rounded w-1/2" />
      </div>
    )
  }

  if (error) {
    return <div className="text-red-400 text-center py-10">{error}</div>
  }

  const src = video?.hls_url || `/api/v1/videos/${id}/manifest.m3u8`

  return (
    <div className="max-w-4xl">
      <VideoPlayer src={src} poster={video?.thumbnail} />

      <div className="mt-4">
        <h1 className="text-xl font-bold">{video?.title || 'Untitled'}</h1>

        <div className="flex items-center justify-between mt-2">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gray-600 rounded-full" />
            <div>
              <p className="text-sm font-medium">{video?.channel || 'Channel'}</p>
              <p className="text-xs text-gray-400">
                {video?.view_count?.toLocaleString() || 0} views
              </p>
            </div>
            <button className="ml-4 bg-white text-black px-4 py-1.5 rounded-full text-sm font-medium">
              Subscribe
            </button>
          </div>

          <div className="flex gap-2">
            <button className="bg-gray-800 px-4 py-1.5 rounded-full text-sm flex items-center gap-1">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M1 21h4V9H1v12zm22-11c0-1.1-.9-2-2-2h-6.31l.95-4.57.03-.32c0-.41-.17-.79-.44-1.06L14.17 1 7.59 7.59C7.22 7.95 7 8.45 7 9v10c0 1.1.9 2 2 2h9c.83 0 1.54-.5 1.84-1.22l3.02-7.05c.09-.23.14-.47.14-.73v-2z" />
              </svg>
              {video?.like_count || 0}
            </button>
          </div>
        </div>

        <div className="mt-4 bg-gray-800 rounded-lg p-4 text-sm">
          <p className="text-gray-300 whitespace-pre-wrap">
            {video?.description || 'No description'}
          </p>
        </div>

        <CommentList videoId={id} />
      </div>
    </div>
  )
}
