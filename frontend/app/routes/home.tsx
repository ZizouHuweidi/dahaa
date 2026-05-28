import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
	return [
		{ title: "Dahaa - Invite-link trivia" },
		{ name: "description", content: "Create a trivia room, share the link, and play instantly." },
	];
}

export default function Home() {
	return (
		<main className="min-h-svh overflow-hidden bg-orange-50 text-slate-950">
			<section className="mx-auto flex min-h-svh w-full max-w-6xl flex-col px-5 pb-[calc(2rem+env(safe-area-inset-bottom))] pt-[calc(1.25rem+env(safe-area-inset-top))] sm:px-8 lg:px-10">
				<nav className="flex items-center justify-between gap-4">
					<a href="/" className="flex items-center gap-3 font-black tracking-tight">
						<span className="grid size-11 place-items-center rounded-2xl bg-slate-950 text-lg text-orange-50 shadow-lg shadow-orange-900/10">د</span>
						<span className="text-xl">Dahaa</span>
					</a>
					<a className="rounded-full border border-slate-900/15 px-4 py-2 text-sm font-bold" href="/create">
						Create
					</a>
				</nav>

				<div className="grid flex-1 items-center gap-10 py-12 lg:grid-cols-[1.05fr_0.95fr] lg:py-16">
					<div className="max-w-2xl">
						<p className="mb-4 inline-flex rounded-full bg-white px-4 py-2 text-sm font-bold text-orange-700 shadow-sm ring-1 ring-orange-200">
							No login. Just share a link.
						</p>
						<h1 className="text-balance text-5xl font-black leading-[0.95] tracking-tight sm:text-6xl lg:text-7xl">
							A bluffing trivia game built for phones first.
						</h1>
						<p className="mt-6 max-w-xl text-lg leading-8 text-slate-700 sm:text-xl">
							Create a room, pick categories, invite friends, write fake answers, and score when people fall for yours.
						</p>
						<div className="mt-8 flex flex-col gap-3 sm:flex-row">
							<a className="inline-flex min-h-12 items-center justify-center rounded-2xl bg-slate-950 px-6 py-3 text-base font-black text-white shadow-xl shadow-slate-950/15" href="/create">
								Start a game
							</a>
							<a className="inline-flex min-h-12 items-center justify-center rounded-2xl bg-white px-6 py-3 text-base font-black text-slate-950 ring-1 ring-slate-900/10" href="/g/DEMO01">
								Preview invite link
							</a>
						</div>
					</div>

					<div className="rounded-[2rem] bg-slate-950 p-4 text-white shadow-2xl shadow-orange-950/20 sm:p-6">
						<div className="rounded-[1.5rem] bg-white/10 p-4 ring-1 ring-white/10">
							<div className="mb-5 flex items-center justify-between text-sm text-orange-100">
								<span>Round 3 of 8</span>
								<span>18s</span>
							</div>
							<h2 className="text-2xl font-black leading-tight">Which city is known as the City of a Thousand Minarets?</h2>
							<div className="mt-6 grid gap-3">
								{["Cairo", "Marrakesh", "Istanbul", "Damascus"].map((answer) => (
									<button key={answer} className="min-h-14 rounded-2xl bg-orange-100 px-4 text-left text-base font-black text-slate-950 active:scale-[0.99]">
										{answer}
									</button>
								))}
							</div>
						</div>
					</div>
				</div>
			</section>
		</main>
	);
}
