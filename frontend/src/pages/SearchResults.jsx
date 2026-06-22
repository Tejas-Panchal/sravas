import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import VideoCard from '../components/VideoCard'
import api from '../services/api'

// SearchResults page — fetches and displays search results from the query string
export default function SearchResults() {
  const [params] = useSearchParams()
  const q = params.get('q') || ''
  const [results, setResults] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!q) {
      setResults([])
      setLoading(false)
      return
    }
    setLoading(true)
    api.get('/videos/search', { params: { q } })
      .then((res) => setResults(res.data?.results || []))
      .catch(() => setResults([]))
      .finally(() => setLoading(false))
  }, [q])

  if (!q) {
    return <div className="text-gray-400 text-center py-10">Enter a search term</div>
  }

  return (
    <div>
      <h1 className="text-lg text-gray-400 mb-6">
        Results for "<span className="text-white">{q}</span>"
      </h1>

      {loading ? (
        <div className="space-y-4 animate-pulse">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="flex gap-4">
              <div className="w-40 aspect-video bg-gray-800 rounded-lg shrink-0" />
              <div className="flex-1">
                <div className="h-4 bg-gray-800 rounded w-3/4 mb-2" />
                <div className="h-3 bg-gray-800 rounded w-1/2" />
              </div>
            </div>
          ))}
        </div>
      ) : results.length === 0 ? (
        <div className="text-gray-400 text-center py-10">No results found</div>
      ) : (
        <div className="space-y-4">
          {results.map((v) => (
            <VideoCard key={v.id} video={v} />
          ))}
        </div>
      )}
    </div>
  )
}
