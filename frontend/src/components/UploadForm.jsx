import { useState, useRef, useCallback } from 'react'
import api from '../services/api'

const ACCEPTED_TYPES = ['video/mp4', 'video/webm', 'video/avi']
const MAX_SIZE = 2 * 1024 * 1024 * 1024 // 2GB

// UploadForm handles drag-and-drop file selection, validation, metadata, and upload with progress
export default function UploadForm() {
  const [state, setState] = useState('idle') // idle | selected | uploading | success | error
  const [dragging, setDragging] = useState(false)
  const [file, setFile] = useState(null)
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [progress, setProgress] = useState(0)
  const [videoId, setVideoId] = useState(null)
  const [errorMsg, setErrorMsg] = useState('')
  const inputRef = useRef()
  const abortRef = useRef(null)

  const validate = useCallback((f) => {
    if (!ACCEPTED_TYPES.includes(f.type)) {
      return 'Only MP4, WebM, and AVI files are accepted'
    }
    if (f.size > MAX_SIZE) {
      return 'File size must be under 2GB'
    }
    return null
  }, [])

  const handleFile = useCallback((f) => {
    const err = validate(f)
    if (err) {
      setErrorMsg(err)
      setState('error')
      return
    }
    setFile(f)
    setTitle(f.name.replace(/\.[^/.]+$/, ''))
    setState('selected')
    setErrorMsg('')
  }, [validate])

  const onDrop = useCallback((e) => {
    e.preventDefault()
    setDragging(false)
    const f = e.dataTransfer.files[0]
    if (f) handleFile(f)
  }, [handleFile])

  const onDragOver = useCallback((e) => {
    e.preventDefault()
    setDragging(true)
  }, [])

  const onDragLeave = useCallback(() => {
    setDragging(false)
  }, [])

  const handleUpload = useCallback(async () => {
    if (!file || !title.trim()) return

    const formData = new FormData()
    formData.append('video', file)
    formData.append('title', title.trim())
    formData.append('description', description.trim())

    const controller = new AbortController()
    abortRef.current = controller

    setState('uploading')
    setProgress(0)

    try {
      const res = await api.post('/videos/upload', formData, {
        signal: controller.signal,
        headers: { 'Content-Type': 'multipart/form-data' },
        onUploadProgress: (e) => {
          if (e.total) setProgress(Math.round((e.loaded / e.total) * 100))
        },
      })
      setVideoId(res.data.id)
      setState('success')
    } catch (err) {
      if (err.name !== 'CanceledError') {
        setErrorMsg(err.response?.data?.error || 'Upload failed')
        setState('error')
      }
    } finally {
      abortRef.current = null
    }
  }, [file, title, description])

  const handleReset = useCallback(() => {
    setFile(null)
    setTitle('')
    setDescription('')
    setProgress(0)
    setVideoId(null)
    setErrorMsg('')
    setState('idle')
  }, [])

  if (state === 'success') {
    return (
      <div className="text-center py-10">
        <div className="w-16 h-16 bg-green-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" /></svg>
        </div>
        <h2 className="text-xl font-bold mb-2">Upload complete</h2>
        <p className="text-gray-400 mb-4">{title}</p>
        <div className="flex gap-3 justify-center">
          <a href={`/watch/${videoId}`} className="bg-white text-black px-4 py-2 rounded-full text-sm font-medium">
            Watch now
          </a>
          <button onClick={handleReset} className="bg-gray-800 px-4 py-2 rounded-full text-sm">
            Upload another
          </button>
        </div>
      </div>
    )
  }

  if (state === 'uploading') {
    return (
      <div className="text-center py-10">
        <p className="text-gray-400 mb-2">Uploading {file?.name}</p>
        <div className="w-full max-w-md mx-auto h-2 bg-gray-700 rounded-full overflow-hidden mb-2">
          <div className="h-full bg-white transition-all" style={{ width: `${progress}%` }} />
        </div>
        <p className="text-sm text-gray-500">{progress}%</p>
        <button onClick={() => abortRef.current?.abort()} className="mt-4 text-sm text-gray-400 hover:text-white">
          Cancel
        </button>
      </div>
    )
  }

  if (state === 'error') {
    return (
      <div className="text-center py-10">
        <div className="w-16 h-16 bg-red-600 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg className="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" /></svg>
        </div>
        <h2 className="text-xl font-bold mb-2">Upload failed</h2>
        <p className="text-gray-400 mb-4">{errorMsg}</p>
        <button onClick={handleReset} className="bg-gray-800 px-4 py-2 rounded-full text-sm">
          Try again
        </button>
      </div>
    )
  }

  // idle or selected state
  return (
    <div>
      {state === 'idle' && (
        <div
          className={`border-2 border-dashed rounded-lg p-12 text-center cursor-pointer transition-colors ${
            dragging ? 'border-white bg-gray-800' : 'border-gray-600 hover:border-gray-400'
          }`}
          onDrop={onDrop}
          onDragOver={onDragOver}
          onDragLeave={onDragLeave}
          onClick={() => inputRef.current?.click()}
        >
          <svg className="w-12 h-12 mx-auto text-gray-400 mb-3" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19 7v2.99s-1.99.01-2 0V7h-3s.01-1.99 0-2h3V2h2v3h3v2h-3zm-3 4V8h-3V5H5c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2v-8h-3zM5 19l3-4 2 3 3-4 4 5H5z" />
          </svg>
          <p className="text-gray-400">Drag & drop a video file here</p>
          <p className="text-gray-500 text-sm mt-2">MP4, WebM, or AVI (max 2GB)</p>
          <input
            ref={inputRef}
            type="file"
            accept="video/mp4,video/webm,video/avi"
            className="hidden"
            onChange={(e) => e.target.files[0] && handleFile(e.target.files[0])}
          />
        </div>
      )}

      {state === 'selected' && (
        <div>
          <div className="bg-gray-800 rounded-lg p-4 mb-4 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <svg className="w-8 h-8 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M18 4l2 4h-3l-2-4h-2l2 4h-3l-2-4H8l2 4H7L5 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V4h-4z" />
              </svg>
              <div>
                <p className="text-sm font-medium">{file?.name}</p>
                <p className="text-xs text-gray-400">{(file?.size / (1024 * 1024)).toFixed(1)} MB</p>
              </div>
            </div>
            <button onClick={handleReset} className="text-gray-400 hover:text-white text-sm">
              Change
            </button>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm text-gray-400 mb-1">Title</label>
              <input
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full bg-gray-800 rounded px-3 py-2 outline-none"
                placeholder="Video title"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-400 mb-1">Description</label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full bg-gray-800 rounded px-3 py-2 outline-none resize-none"
                rows={3}
                placeholder="Optional description"
              />
            </div>
            <button
              onClick={handleUpload}
              disabled={!title.trim()}
              className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white px-6 py-2 rounded-full text-sm font-medium"
            >
              Upload
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
