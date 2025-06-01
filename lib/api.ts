import { Game, Player } from './store/game';

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async post<T>(path: string, data?: any): Promise<{ data: T }> {
    const response = await fetch(`${this.baseUrl}${path}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: data ? JSON.stringify(data) : undefined,
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
  }

  async createGame(playerName: string): Promise<Game> {
    const player: Player = {
      id: '', // Backend will assign ID
      name: playerName,
      score: 0,
      is_connected: true,
      last_seen: new Date().toISOString(),
      is_active: true,
    };
    const response = await this.post<Game>('/api/games', { player });
    return response.data;
  }

  async joinGame(gameCode: string, playerName: string): Promise<Game> {
    const player: Player = {
      id: '', // Backend will assign ID
      name: playerName,
      score: 0,
      is_connected: true,
      last_seen: new Date().toISOString(),
      is_active: true,
    };
    const response = await this.post<Game>('/api/games/join', { gameCode, player });
    return response.data;
  }

  async startGame(gameCode: string): Promise<Game> {
    const response = await this.post<Game>(`/api/games/${gameCode}/start`);
    return response.data;
  }

  async leaveGame(gameCode: string): Promise<Game> {
    const response = await this.post<Game>(`/api/games/${gameCode}/leave`);
    return response.data;
  }

  async submitAnswer(gameCode: string, answer: string): Promise<Game> {
    const response = await this.post<Game>(`/api/games/${gameCode}/answer`, { answer });
    return response.data;
  }
} 