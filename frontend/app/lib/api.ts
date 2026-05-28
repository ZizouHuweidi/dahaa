export type Player = {
  id: string;
  name: string;
  score: number;
  is_connected: boolean;
  is_active: boolean;
  last_seen: string;
};

export type Game = {
  id: string;
  code: string;
  status: "waiting" | "playing" | "ended";
  players: Player[];
  rounds: Round[];
  settings: GameSettings;
  host_id: string;
};

export type Answer = {
  id: string;
  player_id: string;
  text: string;
  kind: "correct" | "fake" | "filler";
  votes: string[];
  created_at: string;
};

export type Round = {
  number: number;
  category: string;
  question: string;
  question_id: string;
  status: "waiting" | "active" | "voting" | "completed";
  current_turn?: {
    player_id: string;
    status: "waiting" | "active" | "ended";
    category?: string;
  };
  answer_pool: {
    correct_answer: string;
    fake_answers: Answer[];
    filler_answers: Answer[];
    options: Answer[];
  };
};

export type GameSettings = {
  rounds: number;
  max_players: number;
  selected_categories: string[];
  time_limits: {
    category_selection: number;
    answer_writing: number;
    voting: number;
  };
};

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

export function gameWebSocketURL(gameID: string) {
  const configured = import.meta.env.VITE_WS_URL as string | undefined;
  if (configured) return `${configured.replace(/\/$/, "")}/ws?game_id=${encodeURIComponent(gameID)}`;

  const url = new URL(API_URL);
  url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
  url.pathname = "/ws";
  url.search = `game_id=${encodeURIComponent(gameID)}`;
  return url.toString();
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  if (!response.ok) {
    let message = `Request failed with ${response.status}`;
    try {
      const body = await response.json();
      if (body?.error) message = body.error;
    } catch {}
    throw new Error(message);
  }

  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

export async function getCategories() {
  const body = await request<{ categories: string[] }>("/api/questions/categories");
  return body.categories;
}

export function createGame(player: Player, settings: GameSettings) {
  return request<Game>("/api/games", {
    method: "POST",
    body: JSON.stringify({ player, settings }),
  });
}

export function getGame(code: string) {
  return request<Game>(`/api/games/${encodeURIComponent(code)}`);
}

export function joinGame(code: string, player: Player) {
  return request<{ message: string }>(`/api/games/${encodeURIComponent(code)}/join`, {
    method: "POST",
    body: JSON.stringify({ player }),
  });
}

export function startGame(code: string) {
  return request<void>(`/api/games/${encodeURIComponent(code)}/start`, { method: "POST" });
}

export function startTurn(code: string, playerId: string) {
  return request<void>(`/api/games/${encodeURIComponent(code)}/turns`, {
    method: "POST",
    body: JSON.stringify({ player_id: playerId }),
  });
}

export function selectCategory(code: string, category: string) {
  return request<void>(`/api/games/${encodeURIComponent(code)}/turns/category`, {
    method: "POST",
    body: JSON.stringify({ category }),
  });
}

export function submitAnswer(code: string, playerId: string, answer: string) {
  return request<{ message: string }>(`/api/games/${encodeURIComponent(code)}/rounds/current/answers`, {
    method: "POST",
    body: JSON.stringify({ player_id: playerId, answer }),
  });
}

export function submitVote(code: string, playerId: string, answerId: string) {
  return request<{ message: string }>(`/api/games/${encodeURIComponent(code)}/rounds/current/votes`, {
    method: "POST",
    body: JSON.stringify({ player_id: playerId, answer_id: answerId }),
  });
}

export function endRound(code: string) {
  return request<{ message: string }>(`/api/games/${encodeURIComponent(code)}/rounds/current/end`, { method: "POST" });
}

export function makeGuestPlayer(name: string): Player {
  return {
    id: crypto.randomUUID(),
    name,
    score: 0,
    is_connected: true,
    is_active: true,
    last_seen: new Date().toISOString(),
  };
}

export function saveGuestPlayer(code: string, player: Player) {
  localStorage.setItem(`dahaa:player:${code}`, JSON.stringify(player));
}

export function loadGuestPlayer(code: string): Player | null {
  const value = localStorage.getItem(`dahaa:player:${code}`);
  if (!value) return null;
  try {
    return JSON.parse(value) as Player;
  } catch {
    return null;
  }
}
