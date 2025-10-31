Zustand user store

Install:

npm install zustand
# or
pnpm add zustand

types (optional):

npm i -D @types/zustand

Usage example (React component):

import { useUserStore } from '../lib/store/userStore'

function MyComponent() {
  const { user, setUser, clearUser } = useUserStore()

  // set after login
  setUser({ id: 123, email: 'a@b.com', token: 'jwt' })

  return <div>{user?.email}</div>
}
