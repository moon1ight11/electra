import { createContext, useContext, useState, useEffect } from 'react';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const restoreSession = async () => {
      try {
        const meRes = await fetch('/api/v1/worker/me', { credentials: 'include' });
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
    await fetch('/api/v1/worker/logout', { method: 'POST', credentials: 'include' });
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