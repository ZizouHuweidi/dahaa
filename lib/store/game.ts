import { create } from 'zustand';
import { apiClient } from '~/lib/api/client';
import { useAuthStore } from './auth';

interface GameSettings {
  rounds: number;
  time_limits: {
    category_selection: number;
    answer_writing: number;
    voting: number;
  };
  selected_categories: string[];
  max_players: number;
}

export interface Player {
  id: string;
  name: string;
  score: number;
  is_connected: boolean;
  last_seen: string;
  is_active: boolean;
}

export interface Game {
  id: string;
  code: string;
  status: 'waiting' | 'playing' | 'ended';
  players: Player[];
  rounds: Round[];
  currentRound?: number;
  settings: GameSettings;
  created_at: string;
  updated_at: string;
  last_activity: string;
  hostId?: string;
}

interface Round {
  number: number;
  category: string;
  question: string;
  question_id: string;
  status: 'waiting' | 'active' | 'voting' | 'completed';
  start_time: string;
  end_time: string;
  current_turn?: Turn;
  answer_pool: AnswerPool;
  timer?: Timer;
}

interface Turn {
  player_id: string;
  start_time: string;
  end_time: string;
  status: 'waiting' | 'active' | 'ended';
  category?: string;
  timer?: Timer;
}

interface AnswerPool {
  correct_answer: string;
  fake_answers: Answer[];
  filler_answers: Answer[];
}

interface Answer {
  id: string;
  player_id: string;
  text: string;
  votes: string[];
  created_at: string;
}

interface Timer {
  type: 'category_selection' | 'answer_writing' | 'voting';
  start_time: string;
  duration: number;
  end_time: string;
}

interface GameState {
  currentGame: Game | null;
  isLoading: boolean;
  error: string | null;
  createGame: (playerName: string) => Promise<void>;
  joinGame: (gameCode: string, playerName: string) => Promise<void>;
  startGame: () => Promise<void>;
  leaveGame: () => Promise<void>;
  submitAnswer: (answer: string) => Promise<void>;
  clearError: () => void;
}

export const useGameStore = create<GameState>((set) => ({
  currentGame: null,
  isLoading: false,
  error: null,

  createGame: async (playerName: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.createGame(playerName);
      set({ currentGame: response, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to create game', isLoading: false });
    }
  },

  joinGame: async (gameCode: string, playerName: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.joinGame(gameCode, playerName);
      set({ currentGame: response, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to join game', isLoading: false });
    }
  },

  startGame: async () => {
    set({ isLoading: true, error: null });
    try {
      const gameCode = useGameStore.getState().currentGame?.code;
      if (!gameCode) throw new Error('No active game');
      const response = await apiClient.startGame(gameCode);
      set({ currentGame: response, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to start game', isLoading: false });
    }
  },

  leaveGame: async () => {
    set({ isLoading: true, error: null });
    try {
      const gameCode = useGameStore.getState().currentGame?.code;
      if (!gameCode) throw new Error('No active game');
      await apiClient.leaveGame(gameCode);
      set({ currentGame: null, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to leave game', isLoading: false });
    }
  },

  submitAnswer: async (answer: string) => {
    set({ isLoading: true, error: null });
    try {
      const gameCode = useGameStore.getState().currentGame?.code;
      if (!gameCode) throw new Error('No active game');
      const response = await apiClient.submitAnswer(gameCode, answer);
      set({ currentGame: response, isLoading: false });
    } catch (error) {
      set({ error: 'Failed to submit answer', isLoading: false });
    }
  },

  clearError: () => set({ error: null }),
})); 