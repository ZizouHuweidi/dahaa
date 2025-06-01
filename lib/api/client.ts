import { useAuthStore } from '~/lib/store/auth';

const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

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

interface Player {
  id: string;
  name: string;
  score: number;
  is_connected: boolean;
  last_seen: string;
  is_active: boolean;
}

interface Game {
  id: string;
  code: string;
  status: 'waiting' | 'playing' | 'ended';
  players: Player[];
  rounds: any[]; // We'll type this properly when needed
  settings: GameSettings;
  created_at: string;
  updated_at: string;
  last_activity: string;
}

class ApiClient {
  private getHeaders(token: string | null) {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
  }

  async createGame(player: Player, settings?: Partial<GameSettings>): Promise<Game> {
    const { token } = useAuthStore.getState();
    const response = await fetch(`${API_URL}/api/games`, {
      method: 'POST',
      headers: this.getHeaders(token),
      body: JSON.stringify({
        player,
        settings,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create game');
    }

    return response.json();
  }

  async joinGame(code: string, player: Player): Promise<void> {
    const { token } = useAuthStore.getState();
    const response = await fetch(`${API_URL}/api/games/join`, {
      method: 'POST',
      headers: this.getHeaders(token),
      body: JSON.stringify({
        code,
        player,
      }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to join game');
    }
  }

  async getGame(code: string): Promise<Game> {
    const { token } = useAuthStore.getState();
    const response = await fetch(`${API_URL}/api/games/${code}`, {
      headers: this.getHeaders(token),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get game');
    }

    return response.json();
  }

  async startGame(code: string): Promise<void> {
    const { token } = useAuthStore.getState();
    const response = await fetch(`${API_URL}/api/games/${code}/start`, {
      method: 'POST',
      headers: this.getHeaders(token),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to start game');
    }
  }
}

export const apiClient = new ApiClient(); 