import { createContext, useContext, useState, useEffect } from 'react';

const BASE = process.env.REACT_APP_API_URL || '/api/v1';
const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const restoreSession = async () => {
      try {
        const meRes = await fetch(`${BASE}/worker/me`, { credentials: 'include' });
        if (meRes.ok) {
          const me = await meRes.json();
          setUser({ id: me.id, name: me.name, role: me.role });
        }
      } catch { }
      setLoading(false);
    };

    restoreSession();
  }, []);

  const login = (id, name, role) => setUser({ id, name, role });

  const logout = async () => {
    await fetch(`${BASE}/worker/logout`, { method: 'POST', credentials: 'include' });
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}