import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Player } from "~/lib/game-schemas";

type GuestSessionState = {
	playersByCode: Record<string, Player>;
	getPlayer: (code: string) => Player | null;
	savePlayer: (code: string, player: Player) => void;
};

export const useGuestSessionStore = create<GuestSessionState>()(
	persist(
		(set, get) => ({
			playersByCode: {},
			getPlayer: (code) => get().playersByCode[code.toUpperCase()] ?? null,
			savePlayer: (code, player) =>
				set((state) => ({
					playersByCode: {
						...state.playersByCode,
						[code.toUpperCase()]: player,
					},
				})),
		}),
		{ name: "dahaa-guest-session" },
	),
);
