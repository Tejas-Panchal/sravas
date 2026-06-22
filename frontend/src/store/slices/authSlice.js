import { createSlice } from '@reduxjs/toolkit'

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    token: null,
    loading: false,
    error: null,
  },
  reducers: {
    setUser(state, action) { state.user = action.payload },
    setToken(state, action) { state.token = action.payload },
    setLoading(state, action) { state.loading = action.payload },
    setError(state, action) { state.error = action.payload },
    logout(state) { state.user = null; state.token = null },
  },
})

export const { setUser, setToken, setLoading, setError, logout } = authSlice.actions
export default authSlice.reducer
