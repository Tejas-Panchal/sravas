import { createSlice } from '@reduxjs/toolkit'

const notificationSlice = createSlice({
  name: 'notifications',
  initialState: {
    items: [],
    unreadCount: 0,
    loading: false,
    error: null,
  },
  reducers: {
    setItems(state, action) {
      state.items = action.payload
      state.unreadCount = action.payload.filter((n) => !n.read).length
    },
    addItem(state, action) {
      state.items.unshift(action.payload)
      if (!action.payload.read) state.unreadCount++
    },
    setLoading(state, action) { state.loading = action.payload },
    setError(state, action) { state.error = action.payload },
  },
})

export const { setItems, addItem, setLoading, setError } = notificationSlice.actions
export default notificationSlice.reducer
