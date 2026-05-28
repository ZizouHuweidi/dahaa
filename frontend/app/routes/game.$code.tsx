import type { Route } from "./+types/game.$code";
import { useEffect, useState } from "react";
import { endRound, gameWebSocketURL, getGame, joinGame, loadGuestPlayer, makeGuestPlayer, saveGuestPlayer, selectCategory, startGame, startTurn, submitAnswer, submitVote, type Game, type Player } from "~/lib/api";

export function meta({ params }: Route.MetaArgs) {
	return [{ title: `Join ${params.code} - Dahaa` }];
}

export default function GameInvite({ params }: Route.ComponentProps) {
	const code = params.code.toUpperCase();
	const [name, setName] = useState("");
	const [player, setPlayer] = useState<Player | null>(null);
	const [game, setGame] = useState<Game | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const [answer, setAnswer] = useState("");
	const [isLive, setIsLive] = useState(false);

	useEffect(() => {
		let cancelled = false;
		const saved = loadGuestPlayer(code);
		if (saved) setPlayer(saved);
		getGame(code)
			.then((loaded) => {
				if (!cancelled) setGame(loaded);
			})
			.catch((err: Error) => {
				if (!cancelled) setError(err.message);
			})
			.finally(() => {
				if (!cancelled) setIsLoading(false);
			});
		return () => {
			cancelled = true;
		};
	}, [code]);

	useEffect(() => {
		if (!player) return;
		const interval = window.setInterval(() => {
			getGame(code).then(setGame).catch(() => {});
		}, isLive ? 10000 : 2500);
		return () => window.clearInterval(interval);
	}, [code, player, isLive]);

	useEffect(() => {
		if (!game?.id || !player) return;
		let closedByEffect = false;
		const socket = new WebSocket(gameWebSocketURL(game.id));
		socket.addEventListener("open", () => setIsLive(true));
		socket.addEventListener("close", () => {
			if (!closedByEffect) setIsLive(false);
		});
		socket.addEventListener("error", () => setIsLive(false));
		socket.addEventListener("message", (event) => {
			try {
				const message = JSON.parse(event.data) as { type: string; payload?: Game | Player };
				if (["game_updated", "game_started", "round_ended", "game_ended"].includes(message.type) && message.payload) {
					setGame(message.payload as Game);
				}
			} catch {}
		});
		return () => {
			closedByEffect = true;
			setIsLive(false);
			socket.close();
		};
	}, [game?.id, player]);

	async function refresh() {
		setGame(await getGame(code));
	}

	async function handleJoin() {
		setError(null);
		if (!name.trim()) {
			setError("Enter your display name.");
			return;
		}
		const guest = makeGuestPlayer(name.trim());
		try {
			await joinGame(code, guest);
			saveGuestPlayer(code, guest);
			setPlayer(guest);
			setGame(await getGame(code));
		} catch (err) {
			setError(err instanceof Error ? err.message : "Failed to join game.");
		}
	}

	async function handleStart() {
		setError(null);
		if (game && game.players.length < 2) {
			setError("Need at least 2 players to start.");
			return;
		}
		try {
			await startGame(code);
			setGame(await getGame(code));
		} catch (err) {
			setError(err instanceof Error ? err.message : "Failed to start game.");
		}
	}

	async function runAction(action: () => Promise<unknown>) {
		setError(null);
		try {
			await action();
			await refresh();
		} catch (err) {
			setError(err instanceof Error ? err.message : "Action failed.");
		}
	}

	const isHost = !!player && !!game && game.host_id === player.id;
	const currentRound = game?.rounds?.[game.rounds.length - 1];
	const currentTurnPlayer = currentRound?.current_turn ? game?.players.find((item) => item.id === currentRound.current_turn?.player_id) : null;
	const hasSubmitted = !!currentRound && !!player && currentRound.answer_pool.fake_answers.some((item) => item.player_id === player.id);
	const hasVoted = !!currentRound && !!player && currentRound.answer_pool.options.some((item) => item.votes.includes(player.id));

	return (
		<main className="min-h-svh bg-orange-50 px-5 py-[calc(1rem+env(safe-area-inset-top))] text-slate-950 sm:px-8">
			<section className="mx-auto flex min-h-[calc(100svh-2rem)] max-w-xl flex-col justify-center pb-[calc(1rem+env(safe-area-inset-bottom))]">
				<div className="rounded-[2rem] bg-white p-5 shadow-2xl shadow-orange-950/10 ring-1 ring-orange-200 sm:p-8">
					<div className="flex items-center justify-between gap-3">
						<p className="text-sm font-black uppercase tracking-[0.2em] text-orange-600">Invite code</p>
						{player ? <span className={`rounded-full px-3 py-1 text-xs font-black ${isLive ? "bg-green-100 text-green-700" : "bg-slate-100 text-slate-600"}`}>{isLive ? "Live" : "Polling"}</span> : null}
					</div>
					<h1 className="mt-3 text-5xl font-black tracking-tight">{code}</h1>
					{game ? <p className="mt-2 text-sm font-bold text-slate-500">Share this page link to invite friends.</p> : null}
					{isLoading ? <p className="mt-4 text-slate-600">Loading room...</p> : null}
					{!isLoading && !player ? (
						<>
							<p className="mt-4 text-slate-600">Enter a name to join this room. Guest sessions are stored locally so reconnects work without accounts.</p>
							<form className="mt-8 grid gap-4">
								<label className="grid gap-2 text-sm font-bold">
									Display name
									<input className="min-h-14 rounded-2xl border border-slate-200 px-4 text-lg" placeholder="Your name" value={name} onChange={(event) => setName(event.target.value)} />
								</label>
								<button type="button" onClick={handleJoin} className="min-h-14 rounded-2xl bg-slate-950 px-5 text-lg font-black text-white">Join game</button>
							</form>
						</>
					) : null}
					{player && game && game.status === "waiting" ? (
						<div className="mt-8 grid gap-5">
							<div className="rounded-3xl bg-orange-50 p-4">
								<p className="text-sm font-bold text-orange-700">Room status</p>
								<p className="mt-1 text-2xl font-black capitalize">{game.status}</p>
							</div>
							<div>
								<p className="mb-3 text-sm font-bold">Players</p>
								<div className="grid gap-2">
									{game.players.map((roomPlayer) => (
										<div key={roomPlayer.id} className="flex min-h-12 items-center justify-between rounded-2xl bg-slate-100 px-4">
											<span className="font-black">{roomPlayer.name}</span>
											<span className="text-sm font-bold text-slate-500">{roomPlayer.id === game.host_id ? "Host" : `${roomPlayer.score} pts`}</span>
										</div>
									))}
								</div>
							</div>
							{isHost ? <button type="button" onClick={handleStart} className="min-h-14 rounded-2xl bg-slate-950 px-5 text-lg font-black text-white">Start game</button> : null}
							{!isHost ? <p className="rounded-2xl bg-slate-100 p-4 text-sm font-bold text-slate-600">Waiting for the host to start.</p> : null}
						</div>
					) : null}
					{player && game && game.status === "playing" ? (
						<div className="mt-8 grid gap-5">
							<div className="grid grid-cols-3 gap-2 text-center">
								{game.players.map((roomPlayer) => (
									<div key={roomPlayer.id} className={`rounded-2xl p-3 ${roomPlayer.id === player.id ? "bg-orange-100" : "bg-slate-100"}`}>
										<p className="truncate text-xs font-bold text-slate-500">{roomPlayer.name}</p>
										<p className="text-xl font-black">{roomPlayer.score}</p>
									</div>
								))}
							</div>
							<div className="rounded-3xl bg-slate-950 p-5 text-white">
								<p className="text-sm font-bold text-orange-200">Round {currentRound?.number ?? 1}</p>
								<h2 className="mt-2 text-2xl font-black">{currentRound?.question || "Waiting for category selection"}</h2>
								{currentRound?.category ? <p className="mt-3 text-sm font-bold text-orange-100">Category: {currentRound.category}</p> : null}
								{currentTurnPlayer ? <p className="mt-2 text-sm font-bold text-orange-100">Category chooser: {currentTurnPlayer.name}</p> : null}
							</div>

							{isHost && currentRound && !currentRound.current_turn ? (
								<button type="button" onClick={() => runAction(() => startTurn(code, player.id))} className="min-h-14 rounded-2xl bg-orange-500 px-5 text-lg font-black text-white">Start category selection</button>
							) : null}
							{currentRound && !currentRound.current_turn && !isHost ? <p className="rounded-2xl bg-slate-100 p-4 text-sm font-bold text-slate-600">Waiting for the host to start this round.</p> : null}

							{isHost && currentRound?.current_turn && !currentRound.question ? (
								<div className="grid gap-3">
									<p className="text-sm font-bold">Choose a category</p>
									<div className="flex flex-wrap gap-2">
										{game.settings.selected_categories.map((category) => (
											<button key={category} type="button" onClick={() => runAction(() => selectCategory(code, category))} className="min-h-12 rounded-full bg-orange-500 px-5 font-black text-white">{category}</button>
										))}
									</div>
								</div>
							) : null}

							{currentRound?.question && currentRound.status === "waiting" && currentRound.current_turn?.player_id !== player.id ? (
								<div className="grid gap-3">
									<label className="grid gap-2 text-sm font-bold">
										Write a convincing fake answer
										<input disabled={hasSubmitted} value={answer} onChange={(event) => setAnswer(event.target.value)} className="min-h-14 rounded-2xl border border-slate-200 px-4 text-lg" placeholder="Your fake answer" />
									</label>
									<button type="button" disabled={hasSubmitted} onClick={() => runAction(async () => { await submitAnswer(code, player.id, answer); setAnswer(""); })} className="min-h-14 rounded-2xl bg-slate-950 px-5 text-lg font-black text-white disabled:opacity-60">{hasSubmitted ? "Answer submitted" : "Submit answer"}</button>
								</div>
							) : null}
							{currentRound?.question && currentRound.status === "waiting" && currentRound.current_turn?.player_id === player.id ? <p className="rounded-2xl bg-slate-100 p-4 text-sm font-bold text-slate-600">You picked the category this round. Waiting for everyone else to submit fake answers.</p> : null}
							{currentRound?.question && currentRound.status === "waiting" && !hasSubmitted && currentRound.current_turn?.player_id !== player.id ? <p className="text-center text-sm font-bold text-slate-500">Submitted answers: {currentRound.answer_pool.fake_answers.length} / {Math.max(game.players.length - 1, 0)}</p> : null}

							{currentRound?.status === "voting" ? (
								<div className="grid gap-3">
									<p className="text-sm font-bold">Pick the answer you think is correct</p>
									{currentRound.answer_pool.options.map((option) => (
										<button key={option.id} type="button" disabled={hasVoted || option.player_id === player.id} onClick={() => runAction(() => submitVote(code, player.id, option.id))} className="min-h-14 rounded-2xl bg-orange-100 px-4 text-left text-lg font-black text-slate-950 disabled:opacity-50">{option.text}</button>
									))}
								</div>
							) : null}

							{currentRound?.status === "completed" ? (
								<div className="grid gap-3 rounded-3xl bg-orange-50 p-4">
									<p className="text-sm font-bold text-orange-700">Correct answer</p>
									<p className="text-2xl font-black">{currentRound.answer_pool.correct_answer}</p>
									<div className="grid gap-2">
										{currentRound.answer_pool.options.map((option) => (
											<div key={option.id} className="rounded-2xl bg-white p-3 text-sm">
												<p className="font-black">{option.text}</p>
												<p className="mt-1 text-slate-600">{option.kind === "correct" ? "Correct answer" : option.kind === "fake" ? "Fake answer" : "Filler answer"} · {option.votes.length} votes</p>
											</div>
										))}
									</div>
									{isHost ? <button type="button" onClick={() => runAction(() => endRound(code))} className="min-h-12 rounded-2xl bg-slate-950 px-5 font-black text-white">Close round</button> : null}
								</div>
							) : null}
						</div>
					) : null}
					{player && game && game.status === "ended" ? (
						<div className="mt-8 grid gap-3">
							<p className="text-sm font-bold text-orange-700">Final scores</p>
							{[...game.players].sort((a, b) => b.score - a.score).map((roomPlayer, index) => (
								<div key={roomPlayer.id} className="flex min-h-14 items-center justify-between rounded-2xl bg-slate-100 px-4">
									<span className="font-black">{index + 1}. {roomPlayer.name}</span>
									<span className="font-black">{roomPlayer.score}</span>
								</div>
							))}
							<a href="/create" className="mt-3 inline-flex min-h-12 items-center justify-center rounded-2xl bg-slate-950 px-5 font-black text-white">Create another game</a>
						</div>
					) : null}
					{error ? <p className="mt-5 rounded-2xl bg-red-50 p-4 text-sm font-bold text-red-700">{error}</p> : null}
				</div>
			</section>
		</main>
	);
}
