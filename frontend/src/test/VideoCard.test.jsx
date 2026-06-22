import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, it, expect } from 'vitest'
import VideoCard from '../components/VideoCard'

const mockVideo = {
  id: 'abc-123',
  title: 'Test Video Title',
  channel: 'Test Channel',
  view_count: 1500,
  duration: 125,
  uploaded_at: '2024-01-15T10:00:00Z',
}

describe('VideoCard', () => {
  it('renders video title and channel', () => {
    render(
      <MemoryRouter>
        <VideoCard video={mockVideo} />
      </MemoryRouter>
    )
    expect(screen.getByText('Test Video Title')).toBeInTheDocument()
    expect(screen.getByText('Test Channel')).toBeInTheDocument()
  })

  it('renders view count', () => {
    render(
      <MemoryRouter>
        <VideoCard video={mockVideo} />
      </MemoryRouter>
    )
    expect(screen.getByText(/1,500 views/)).toBeInTheDocument()
  })

  it('formats duration correctly', () => {
    render(
      <MemoryRouter>
        <VideoCard video={mockVideo} />
      </MemoryRouter>
    )
    expect(screen.getByText('2:05')).toBeInTheDocument()
  })

  it('returns null for no video', () => {
    const { container } = render(
      <MemoryRouter>
        <VideoCard video={null} />
      </MemoryRouter>
    )
    expect(container.firstChild).toBeNull()
  })
})
