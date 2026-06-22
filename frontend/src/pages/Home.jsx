import { useState, useEffect, useCallback, useRef } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import VideoCard from '../components/VideoCard'
import LoadingSkeleton from '../components/LoadingSkeleton'
import api from '../services/api'

const PAGE_SIZE = 20

// Home page — video feed with infinite scroll
export default function Home() {
  const dispatch = useDispatch()
  const { items: videos, loading, error } = useSelector((s) => s.videos)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(true)
  const observerRef = useRef()
  const sentinelRef = useRef()

  const fetchVideos = useCallback(async () => {
    dispatch({ type: 'videos/setLoading', payload: true })
    try {
      const res = await api.get('/videos', { params: { page, page_size: PAGE_SIZE } })
      const data = res.data
      if (page === 1) {
        dispatch({ type: 'videos/setItems', payload: data.results || [] })
      } else {
        dispatch({ type: 'videos/setItems', payload: [...videos, ...(data.results || [])] })
      }
      setHasMore(data.results?.length === PAGE_SIZE)
    } catch (err) {
      dispatch({ type: 'videos/setError', payload: err.message })
    } finally {
      dispatch({ type: 'videos/setLoading', payload: false })
    }
  }, [page])

  useEffect(() => {
    fetchVideos()
  }, [page])

  useEffect(() => {
    if (!sentinelRef.current) return
    observerRef.current = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading) {
          setPage((p) => p + 1)
        }
      },
      { threshold: 0.1 }
    )
    observerRef.current.observe(sentinelRef.current)
    return () => observerRef.current?.disconnect()
  }, [hasMore, loading])

  if (error) {
    return <div className="text-red-400 text-center py-10">{error}</div>
  }

  return (
    <div>
      {videos.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {videos.map((video) => (
            <VideoCard key={video.id} video={video} />
          ))}
        </div>
      )}

      {loading && <LoadingSkeleton count={PAGE_SIZE} />}

      {!loading && videos.length === 0 && (
        <p className="text-gray-400 text-center py-10">No videos yet</p>
      )}

      {hasMore && <div ref={sentinelRef} className="h-4" />}
    </div>
  )
}
