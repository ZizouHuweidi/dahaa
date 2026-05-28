import type { Route } from "./+types/create";
import { useEffect, useState, startTransition } from "react";
import { useNavigate } from "react-router";
import { createGame, getCategories, makeGuestPlayer, saveGuestPlayer } from "~/lib/api";

export function meta({}: Route.MetaArgs) {
	return [{ title: "Create Game - Dahaa" }];
}

export default function CreateGame() {
	const navigate = useNavigate();
	const [name, setName] = useState("");
	const [rounds, setRounds] = useState(8);
	const [categories, setCategories] = useState<string[]>([]);
	const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
	const [error, setError] = useState<string | null>(null);
	const [isLoading, setIsLoading] = useState(false);

	useEffect(() => {
		let cancelled = false;
		getCategories()
			.then((loaded) => {
				if (cancelled) return;
				startTransition(() => {
					setCategories(loaded);
					setSelectedCategories(loaded.slice(0, Math.min(3, loaded.length)));
				});
			})
			.catch((err: Error) => {
				if (!cancelled) setError(err.message);
			});
		return () => {
			cancelled = true;
		};
	}, []);

	async function handleCreate() {
		setError(null);
		if (!name.trim()) {
			setError("Enter your display name.");
			return;
		}
		if (selectedCategories.length === 0) {
			setError("Add questions first, then select at least one category.");
			return;
		}
		setIsLoading(true);
		try {
			const player = makeGuestPlayer(name.trim());
			const game = await createGame(player, {
				rounds,
				max_players: 8,
				selected_categories: selectedCategories,
				time_limits: {
					category_selection: 20,
					answer_writing: 45,
					voting: 30,
				},
			});
			saveGuestPlayer(game.code, player);
			navigate(`/g/${game.code}`);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Failed to create game.");
		} finally {
			setIsLoading(false);
		}
	}

	function toggleCategory(category: string) {
		setSelectedCategories((current) => current.includes(category) ? current.filter((item) => item !== category) : [...current, category]);
	}

	return (
		<main className="min-h-svh bg-slate-950 px-5 py-[calc(1.25rem+env(safe-area-inset-top))] text-white sm:px-8">
			<section className="mx-auto flex min-h-[calc(100svh-2.5rem)] max-w-3xl flex-col justify-center pb-[calc(1rem+env(safe-area-inset-bottom))]">
				<a href="/" className="mb-8 text-sm font-bold text-orange-200">Back to Dahaa</a>
				<div className="rounded-[2rem] bg-white p-5 text-slate-950 shadow-2xl sm:p-8">
					<p className="text-sm font-black uppercase tracking-[0.2em] text-orange-600">MVP flow</p>
					<h1 className="mt-3 text-4xl font-black tracking-tight sm:text-5xl">Create a room</h1>
					<p className="mt-4 text-slate-600">Create a guest-hosted room and share the invite link. Accounts can come later for profiles and progression.</p>
					<form className="mt-8 grid gap-4">
						<label className="grid gap-2 text-sm font-bold">
							Your display name
							<input className="min-h-12 rounded-2xl border border-slate-200 px-4 text-base" placeholder="Zizou" value={name} onChange={(event) => setName(event.target.value)} />
						</label>
						<label className="grid gap-2 text-sm font-bold">
							Rounds
							<select className="min-h-12 rounded-2xl border border-slate-200 px-4 text-base" value={rounds} onChange={(event) => setRounds(Number(event.target.value))}>
								<option value="5">5 rounds</option>
								<option value="8">8 rounds</option>
								<option value="10">10 rounds</option>
							</select>
						</label>
						<div className="grid gap-2">
							<p className="text-sm font-bold">Categories</p>
							{categories.length === 0 ? <p className="rounded-2xl bg-orange-50 p-4 text-sm text-orange-900">No question categories found yet. Seed questions with <code>/api/questions/bulk</code> before creating rooms.</p> : null}
							<div className="flex flex-wrap gap-2">
								{categories.map((category) => (
									<button key={category} type="button" onClick={() => toggleCategory(category)} className={`min-h-11 rounded-full px-4 text-sm font-black ${selectedCategories.includes(category) ? "bg-orange-500 text-white" : "bg-slate-100 text-slate-700"}`}>
										{category}
									</button>
								))}
							</div>
						</div>
						<p className="rounded-2xl bg-slate-100 p-4 text-sm font-bold text-slate-600">After creating a room, share the invite link from the browser address bar. Players join as guests and can reconnect from the same device.</p>
						{error ? <p className="rounded-2xl bg-red-50 p-4 text-sm font-bold text-red-700">{error}</p> : null}
						<button type="button" disabled={isLoading} onClick={handleCreate} className="mt-2 min-h-12 rounded-2xl bg-slate-950 px-5 font-black text-white disabled:opacity-60">{isLoading ? "Creating..." : "Create invite link"}</button>
					</form>
				</div>
			</section>
		</main>
	);
}
