import { create } from 'zustand';
import { useAuthStore } from './auth';
import { websocketService } from '../services/websocket';

interface Player {
  id: string;
  username: string;
  displayName: string;
  score: number;
  isHost: boolean;
}

interface GameSettings {
  maxPlayers: number;
  rounds: number;
  timePerRound: number;
  categories: string[];
}

interface Game {
  id: string;
  code: string;
  status: 'waiting' | 'playing' | 'ended';
  players: Player[];
  settings: GameSettings;
  currentRound: number;
  currentCategory: string | null;
  currentPhase: 'category' | 'answer' | 'vote' | 'end';
  timeRemaining: number;
}

interface GameState {
  currentGame: Game | null;
  isLoading: boolean;
  error: string | null;
  createGame: (settings: GameSettings) => Promise<void>;
  joinGame: (code: string) => Promise<void>;
  startGame: () => Promise<void>;
  submitAnswer: (answer: string) => Promise<void>;
  submitVote: (answerId: string) => Promise<void>;
  endRound: () => Promise<void>;
  endGame: () => Promise<void>;
  clearError: () => void;
}

const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

export const useGameStore = create<GameState>()((set, get) => ({
  currentGame: null,
  isLoading: false,
  error: null,

  createGame: async (settings: GameSettings) => {
    set({ isLoading: true, error: null });
    try {
      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(settings),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to create game');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to create game', isLoading: false });
    }
  },

  joinGame: async (code: string) => {
    set({ isLoading: true, error: null });
    try {
      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${code}/join`, {
        method: 'POST',
        headers: getAuthHeaders(),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to join game');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to join game', isLoading: false });
    }
  },

  startGame: async () => {
    set({ isLoading: true, error: null });
    try {
      const { currentGame } = get();
      if (!currentGame) throw new Error('No active game');

      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${currentGame.id}/start`, {
        method: 'POST',
        headers: getAuthHeaders(),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to start game');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to start game', isLoading: false });
    }
  },

  submitAnswer: async (answer: string) => {
    set({ isLoading: true, error: null });
    try {
      const { currentGame } = get();
      if (!currentGame) throw new Error('No active game');

      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${currentGame.id}/answer`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ answer }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to submit answer');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to submit answer', isLoading: false });
    }
  },

  submitVote: async (answerId: string) => {
    set({ isLoading: true, error: null });
    try {
      const { currentGame } = get();
      if (!currentGame) throw new Error('No active game');

      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${currentGame.id}/vote`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ answerId }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to submit vote');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to submit vote', isLoading: false });
    }
  },

  endRound: async () => {
    set({ isLoading: true, error: null });
    try {
      const { currentGame } = get();
      if (!currentGame) throw new Error('No active game');

      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${currentGame.id}/round/end`, {
        method: 'POST',
        headers: getAuthHeaders(),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to end round');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to end round', isLoading: false });
    }
  },

  endGame: async () => {
    set({ isLoading: true, error: null });
    try {
      const { currentGame } = get();
      if (!currentGame) throw new Error('No active game');

      const { getAuthHeaders } = useAuthStore.getState();
      const response = await fetch(`${API_URL}/games/${currentGame.id}/end`, {
        method: 'POST',
        headers: getAuthHeaders(),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to end game');
      }

      set({ currentGame: data, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to end game', isLoading: false });
    }
  },

  clearError: () => {
    set({ error: null });
  },
}));

// Subscribe to WebSocket events
websocketService.subscribe('game_update', (data) => {
  useGameStore.setState({ currentGame: data });
});

websocketService.subscribe('player_update', (data) => {
  const { currentGame } = useGameStore.getState();
  if (currentGame) {
    useGameStore.setState({
      currentGame: {
        ...currentGame,
        players: data,
      },
    });
  }
});

websocketService.subscribe('round_update', (data) => {
  const { currentGame } = useGameStore.getState();
  if (currentGame) {
    useGameStore.setState({
      currentGame: {
        ...currentGame,
        currentRound: data.round,
        currentCategory: data.category,
        currentPhase: data.phase,
      },
    });
  }
});

websocketService.subscribe('timer_update', (data) => {
  const { currentGame } = useGameStore.getState();
  if (currentGame) {
    useGameStore.setState({
      currentGame: {
        ...currentGame,
        timeRemaining: data.timeRemaining,
      },
    });
  }
}); 