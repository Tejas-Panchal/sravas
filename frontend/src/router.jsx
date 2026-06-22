import { createBrowserRouter } from 'react-router-dom'

import Layout from './components/Layout'
import Home from './pages/Home'
import Upload from './pages/Upload'
import Watch from './pages/Watch'
import SearchResults from './pages/SearchResults'
import UserProfile from './pages/UserProfile'
import Trending from './pages/Trending'
import Subscriptions from './pages/Subscriptions'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <Home /> },
      { path: 'upload', element: <Upload /> },
      { path: 'watch/:id', element: <Watch /> },
      { path: 'search', element: <SearchResults /> },
      { path: 'channel/:id', element: <UserProfile /> },
      { path: 'trending', element: <Trending /> },
      { path: 'feed/subscriptions', element: <Subscriptions /> },
    ],
  },
])
