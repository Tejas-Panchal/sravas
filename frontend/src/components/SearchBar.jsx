import { useState, useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../services/api'

// SearchBar provides an input with debounced autocomplete suggestions
export default function SearchBar() {
  const [query, setQuery] = useState('')
  const [suggestions, setSuggestions] = useState([])
  const [showDropdown, setShowDropdown] = useState(false)
  const navigate = useNavigate()
  const wrapperRef = useRef()
  const debounceRef = useRef()

  useEffect(() => {
    if (!query.trim()) {
      setSuggestions([])
      return
    }
    clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(async () => {
      try {
        const res = await api.get('/videos/autocomplete', { params: { q: query } })
        setSuggestions(res.data?.suggestions || [])
      } catch {
        setSuggestions([])
      }
    }, 200)
    return () => clearTimeout(debounceRef.current)
  }, [query])

  useEffect(() => {
    const handleClick = (e) => {
      if (wrapperRef.current && !wrapperRef.current.contains(e.target)) {
        setShowDropdown(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [])

  const handleSubmit = (e) => {
    e.preventDefault()
    if (query.trim()) {
      navigate(`/search?q=${encodeURIComponent(query.trim())}`)
      setShowDropdown(false)
    }
  }

  const pickSuggestion = (s) => {
    setQuery(s)
    navigate(`/search?q=${encodeURIComponent(s)}`)
    setShowDropdown(false)
  }

  return (
    <div ref={wrapperRef} className="relative flex-1 max-w-xl mx-8">
      <form onSubmit={handleSubmit} className="relative">
        <input
          type="text"
          value={query}
          onChange={(e) => { setQuery(e.target.value); setShowDropdown(true) }}
          onFocus={() => query.trim() && setShowDropdown(true)}
          placeholder="Search"
          className="w-full bg-gray-800 text-white px-4 py-2 rounded-full outline-none"
        />
        {showDropdown && suggestions.length > 0 && (
          <ul className="absolute top-full left-0 right-0 mt-1 bg-gray-800 rounded-lg shadow-lg max-h-60 overflow-y-auto z-50">
            {suggestions.map((s, i) => (
              <li
                key={i}
                className="px-4 py-2 text-sm cursor-pointer hover:bg-gray-700"
                onMouseDown={() => pickSuggestion(s)}
              >
                <span className="text-gray-400 mr-2">🔍</span>{s}
              </li>
            ))}
          </ul>
        )}
      </form>
    </div>
  )
}
