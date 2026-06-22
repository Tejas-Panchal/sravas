import { Link } from 'react-router-dom'

// VideoCard displays a video thumbnail with metadata
export default function VideoCard({ video }) {
  if (!video) return null

  const views = video.view_count?.toLocaleString() || '0'
  const time = video.uploaded_at
    ? new Date(video.uploaded_at).toLocaleDateString()
    : ''

  return (
    <Link to={`/watch/${video.id}`} className="group block">
      <div className="relative aspect-video bg-gray-800 rounded-lg overflow-hidden mb-2">
        {video.thumbnail ? (
          <img src={video.thumbnail} alt={video.title} className="w-full h-full object-cover" />
        ) : (
          <div className="w-full h-full bg-gray-700" />
        )}
        {video.duration && (
          <span className="absolute bottom-1 right-1 bg-black/80 text-xs px-1 rounded">
            {formatDuration(video.duration)}
          </span>
        )}
      </div>
      <div className="flex gap-3">
        <div className="w-9 h-9 bg-gray-600 rounded-full shrink-0 mt-1" />
        <div className="min-w-0">
          <h3 className="text-sm font-medium line-clamp-2 leading-tight">{video.title}</h3>
          <p className="text-xs text-gray-400 mt-1">{video.channel || 'Channel'}</p>
          <p className="text-xs text-gray-400">{views} views · {time}</p>
        </div>
      </div>
    </Link>
  )
}

function formatDuration(seconds) {
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
