import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface User {
  id: string;
  username: string;
  email: string;
  displayName: string;
}

interface AuthState {
  token: string | null;
  user: User | null;
  isLoading: boolean;
  error: string | null;
  register: (username: string, email: string, password: string, displayName: string) => Promise<void>;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  clearError: () => void;
  getAuthHeaders: () => Record<string, string>;
}

const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      isLoading: false,
      error: null,

      getAuthHeaders: () => {
        const { token } = get();
        return {
          'Content-Type': 'application/json',
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
        };
      },

      register: async (username: string, email: string, password: string, displayName: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(`${API_URL}/users/register`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              username,
              email,
              password,
              display_name: displayName,
            }),
          });

          const data = await response.json();

          if (!response.ok) {
            throw new Error(data.error || 'Failed to register');
          }

          set({ user: data, isLoading: false });
        } catch (error) {
          set({ error: error instanceof Error ? error.message : 'Failed to register', isLoading: false });
        }
      },

      login: async (username: string, password: string) => {
        set({ isLoading: true, error: null });
        try {
          const response = await fetch(`${API_URL}/users/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              username,
              password,
            }),
          });

          const data = await response.json();

          if (!response.ok) {
            throw new Error(data.error || 'Failed to login');
          }

          // Fetch user data after successful login
          const userResponse = await fetch(`${API_URL}/users/me`, {
            headers: {
              Authorization: `Bearer ${data.token}`,
            },
          });

          if (!userResponse.ok) {
            throw new Error('Failed to fetch user data');
          }

          const userData = await userResponse.json();

          set({ 
            token: data.token, 
            user: userData,
            isLoading: false 
          });
        } catch (error) {
          set({ error: error instanceof Error ? error.message : 'Failed to login', isLoading: false });
        }
      },

      logout: () => {
        set({ token: null, user: null });
      },

      clearError: () => {
        set({ error: null });
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => AsyncStorage),
    }
  )
); 