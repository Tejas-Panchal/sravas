import { useState, useEffect, useCallback } from 'react'
import { useSelector } from 'react-redux'
import api from '../services/api'

// CommentForm — textarea with submit button for posting a comment
function CommentForm({ videoId, onAdded }) {
  const [text, setText] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    if (!text.trim() || submitting) return
    setSubmitting(true)
    try {
      await api.post(`/videos/${videoId}/comments`, { text: text.trim() })
      setText('')
      onAdded?.()
    } catch {} finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex gap-3 mb-6">
      <div className="w-9 h-9 bg-gray-600 rounded-full shrink-0" />
      <div className="flex-1">
        <textarea
          value={text}
          onChange={(e) => setText(e.target.value)}
          className="w-full bg-transparent border-b border-gray-600 focus:border-white outline-none resize-none text-sm pb-1"
          rows={1}
          placeholder="Add a comment..."
        />
        {text.trim() && (
          <div className="flex justify-end gap-2 mt-2">
            <button onClick={() => setText('')} className="text-sm text-gray-400 px-3 py-1 rounded-full hover:bg-gray-800">
              Cancel
            </button>
            <button
              onClick={handleSubmit}
              disabled={submitting}
              className="text-sm bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white px-4 py-1 rounded-full"
            >
              Comment
            </button>
          </div>
        )}
      </div>
    </div>
  )
}

// CommentItem — single comment with like toggle and reply form
function CommentItem({ comment, videoId, onAdded }) {
  const [showReply, setShowReply] = useState(false)
  const [replyText, setReplyText] = useState('')
  const [liked, setLiked] = useState(comment.liked || false)
  const [likeCount, setLikeCount] = useState(comment.like_count || 0)

  const handleLike = async () => {
    try {
      await api.post(`/comments/${comment.id}/like`)
      setLiked(!liked)
      setLikeCount((c) => liked ? c - 1 : c + 1)
    } catch {}
  }

  const handleReply = async () => {
    if (!replyText.trim()) return
    try {
      await api.post(`/comments/${comment.id}/replies`, { text: replyText.trim() })
      setReplyText('')
      setShowReply(false)
      onAdded?.()
    } catch {}
  }

  const time = comment.created_at
    ? timeAgo(new Date(comment.created_at))
    : ''

  return (
    <div className="mb-4">
      <div className="flex gap-3">
        <div className="w-9 h-9 bg-gray-600 rounded-full shrink-0" />
        <div className="flex-1">
          <div className="flex items-center gap-2 text-xs text-gray-400">
            <span className="text-white font-medium">{comment.username || 'User'}</span>
            <span>{time}</span>
          </div>
          <p className="text-sm mt-1">{comment.text}</p>
          <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
            <button onClick={handleLike} className={`flex items-center gap-1 hover:text-white ${liked ? 'text-white' : ''}`}>
              <svg className="w-4 h-4" fill={liked ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
              </svg>
              {likeCount > 0 && likeCount}
            </button>
            <button onClick={() => setShowReply(!showReply)} className="hover:text-white">
              Reply
            </button>
          </div>

          {showReply && (
            <div className="mt-3 flex gap-2">
              <input
                type="text"
                value={replyText}
                onChange={(e) => setReplyText(e.target.value)}
                className="flex-1 bg-gray-800 rounded px-3 py-1.5 text-sm outline-none"
                placeholder="Write a reply..."
                onKeyDown={(e) => e.key === 'Enter' && handleReply()}
              />
              <button
                onClick={handleReply}
                disabled={!replyText.trim()}
                className="text-sm bg-blue-600 disabled:bg-gray-600 text-white px-3 py-1 rounded-full"
              >
                Reply
              </button>
            </div>
          )}

          {comment.replies?.length > 0 && (
            <div className="mt-3 pl-2 border-l border-gray-700">
              {comment.replies.map((r) => (
                <CommentItem key={r.id} comment={r} videoId={videoId} onAdded={onAdded} />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

// CommentList — fetches and displays comments for a video
export default function CommentList({ videoId }) {
  const { user } = useSelector((s) => s.auth)
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(true)

  const fetchComments = useCallback(() => {
    if (!videoId) return
    setLoading(true)
    api.get(`/videos/${videoId}/comments`)
      .then((res) => setComments(res.data?.comments || []))
      .catch(() => setComments([]))
      .finally(() => setLoading(false))
  }, [videoId])

  useEffect(() => { fetchComments() }, [fetchComments])

  return (
    <div className="mt-6">
      <h2 className="text-lg font-bold mb-4">
        {comments.length} comment{comments.length !== 1 && 's'}
      </h2>

      {user && <CommentForm videoId={videoId} onAdded={fetchComments} />}

      {loading ? (
        <div className="space-y-4 animate-pulse">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="flex gap-3">
              <div className="w-9 h-9 bg-gray-800 rounded-full" />
              <div className="flex-1">
                <div className="h-3 bg-gray-800 rounded w-1/4 mb-2" />
                <div className="h-3 bg-gray-800 rounded w-3/4" />
              </div>
            </div>
          ))}
        </div>
      ) : comments.length === 0 ? (
        <p className="text-gray-400 text-sm">No comments yet</p>
      ) : (
        comments.map((c) => (
          <CommentItem key={c.id} comment={c} videoId={videoId} onAdded={fetchComments} />
        ))
      )}
    </div>
  )
}

function timeAgo(date) {
  const secs = Math.floor((Date.now() - date.getTime()) / 1000)
  if (secs < 60) return 'just now'
  const mins = Math.floor(secs / 60)
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  const days = Math.floor(hrs / 24)
  return `${days}d ago`
}
