import { createSlice } from '@reduxjs/toolkit'

const videoSlice = createSlice({
  name: 'videos',
  initialState: {
    items: [],
    current: null,
    loading: false,
    error: null,
  },
  reducers: {
    setItems(state, action) { state.items = action.payload },
    setCurrent(state, action) { state.current = action.payload },
    setLoading(state, action) { state.loading = action.payload },
    setError(state, action) { state.error = action.payload },
  },
})

export const { setItems, setCurrent, setLoading, setError } = videoSlice.actions
export default videoSlice.reducer
