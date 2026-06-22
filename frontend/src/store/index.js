import { configureStore } from '@reduxjs/toolkit'
import authReducer from './slices/authSlice'
import videoReducer from './slices/videoSlice'
import commentReducer from './slices/commentSlice'
import notificationReducer from './slices/notificationSlice'

export const store = configureStore({
  reducer: {
    auth: authReducer,
    videos: videoReducer,
    comments: commentReducer,
    notifications: notificationReducer,
  },
})
