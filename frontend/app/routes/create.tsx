import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import {
	categoriesQueryKey,
	createGame,
	getCategories,
	makeGuestPlayer,
} from "~/lib/game-api";
import { type CreateGameForm, createGameFormSchema } from "~/lib/game-schemas";
import { useGuestSessionStore } from "~/stores/guest-session-store";
import type { Route } from "./+types/create";

export function meta(_args: Route.MetaArgs) {
	return [{ title: "Create Game - Dahaa" }];
}

export default function CreateGame() {
	const navigate = useNavigate();
	const savePlayer = useGuestSessionStore((state) => state.savePlayer);
	const { data: categories = [], error: categoriesError } = useQuery({
		queryKey: categoriesQueryKey,
		queryFn: getCategories,
	});
	const form = useForm<CreateGameForm>({
		resolver: zodResolver(createGameFormSchema),
		defaultValues: { name: "", rounds: 8, selectedCategories: [] },
	});
	const selectedCategories = form.watch("selectedCategories");
	const createMutation = useMutation({
		mutationFn: async (values: CreateGameForm) => {
			const player = makeGuestPlayer(values.name.trim());
			const game = await createGame(player, {
				rounds: values.rounds,
				max_players: 8,
				selected_categories: values.selectedCategories,
				time_limits: { category_selection: 20, answer_writing: 45, voting: 30 },
			});
			return { game, player };
		},
		onSuccess: ({ game, player }) => {
			savePlayer(game.code, player);
			navigate(`/g/${game.code}`);
		},
	});

	function toggleCategory(category: string) {
		const current = form.getValues("selectedCategories");
		form.setValue(
			"selectedCategories",
			current.includes(category)
				? current.filter((item) => item !== category)
				: [...current, category],
			{
				shouldDirty: true,
				shouldValidate: true,
			},
		);
	}

	return (
		<main className="min-h-svh bg-slate-950 px-5 py-[calc(1.25rem+env(safe-area-inset-top))] text-white sm:px-8">
			<section className="mx-auto flex min-h-[calc(100svh-2.5rem)] max-w-3xl flex-col justify-center pb-[calc(1rem+env(safe-area-inset-bottom))]">
				<a href="/" className="mb-8 text-sm font-bold text-orange-200">
					Back to Dahaa
				</a>
				<div className="rounded-[2rem] bg-white p-5 text-slate-950 shadow-2xl sm:p-8">
					<p className="text-sm font-black uppercase tracking-[0.2em] text-orange-600">
						MVP flow
					</p>
					<h1 className="mt-3 text-4xl font-black tracking-tight sm:text-5xl">
						Create a room
					</h1>
					<p className="mt-4 text-slate-600">
						Create a guest-hosted room and share the invite link. Accounts can
						come later for profiles and progression.
					</p>
					<form
						className="mt-8 grid gap-4"
						onSubmit={form.handleSubmit((values) =>
							createMutation.mutate(values),
						)}
					>
						<label className="grid gap-2 text-sm font-bold">
							Your display name
							<input
								className="min-h-12 rounded-2xl border border-slate-200 px-4 text-base"
								placeholder="Zizou"
								{...form.register("name")}
							/>
						</label>
						{form.formState.errors.name ? (
							<p className="text-sm font-bold text-red-700">
								{form.formState.errors.name.message}
							</p>
						) : null}

						<label className="grid gap-2 text-sm font-bold">
							Rounds
							<select
								className="min-h-12 rounded-2xl border border-slate-200 px-4 text-base"
								{...form.register("rounds", { valueAsNumber: true })}
							>
								<option value="5">5 rounds</option>
								<option value="8">8 rounds</option>
								<option value="10">10 rounds</option>
							</select>
						</label>

						<div className="grid gap-2">
							<p className="text-sm font-bold">Categories</p>
							{categories.length === 0 ? (
								<p className="rounded-2xl bg-orange-50 p-4 text-sm text-orange-900">
									No categories found. Run backend migrations to seed questions.
								</p>
							) : null}
							<div className="flex flex-wrap gap-2">
								{categories.map((category) => (
									<button
										key={category}
										type="button"
										onClick={() => toggleCategory(category)}
										className={`min-h-11 rounded-full px-4 text-sm font-black ${selectedCategories.includes(category) ? "bg-orange-500 text-white" : "bg-slate-100 text-slate-700"}`}
									>
										{category}
									</button>
								))}
							</div>
						</div>
						{form.formState.errors.selectedCategories ? (
							<p className="text-sm font-bold text-red-700">
								{form.formState.errors.selectedCategories.message}
							</p>
						) : null}
						<p className="rounded-2xl bg-slate-100 p-4 text-sm font-bold text-slate-600">
							After creating a room, share the invite link from the browser
							address bar. Players join as guests and can reconnect from the
							same device.
						</p>
						{categoriesError || createMutation.error ? (
							<p className="rounded-2xl bg-red-50 p-4 text-sm font-bold text-red-700">
								{(categoriesError || createMutation.error)?.message}
							</p>
						) : null}
						<button
							type="submit"
							disabled={createMutation.isPending}
							className="mt-2 min-h-12 rounded-2xl bg-slate-950 px-5 font-black text-white disabled:opacity-60"
						>
							{createMutation.isPending ? "Creating..." : "Create invite link"}
						</button>
					</form>
				</div>
			</section>
		</main>
	);
}
