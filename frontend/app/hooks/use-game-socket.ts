import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { gameQueryKey, gameWebSocketURL } from "~/lib/game-api";
import { type Game, gameSchema } from "~/lib/game-schemas";

export function useGameSocket(
	code: string,
	game: Game | null | undefined,
	enabled: boolean,
) {
	const queryClient = useQueryClient();
	const [isLive, setIsLive] = useState(false);

	useEffect(() => {
		if (!enabled || !game?.id) return;
		let closedByEffect = false;
		const socket = new WebSocket(gameWebSocketURL(game.id));

		socket.addEventListener("open", () => setIsLive(true));
		socket.addEventListener("close", () => {
			if (!closedByEffect) setIsLive(false);
		});
		socket.addEventListener("error", () => setIsLive(false));
		socket.addEventListener("message", (event) => {
			try {
				const message = JSON.parse(event.data) as {
					type: string;
					payload?: unknown;
				};
				if (
					[
						"game_updated",
						"game_started",
						"round_ended",
						"game_ended",
					].includes(message.type) &&
					message.payload
				) {
					queryClient.setQueryData(
						gameQueryKey(code),
						gameSchema.parse(message.payload),
					);
				}
			} catch {}
		});

		return () => {
			closedByEffect = true;
			setIsLive(false);
			socket.close();
		};
	}, [code, enabled, game?.id, queryClient]);

	return isLive;
}
