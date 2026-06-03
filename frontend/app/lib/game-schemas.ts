import { z } from "zod";

export const playerSchema = z.object({
	id: z.string(),
	name: z.string(),
	score: z.number(),
	is_connected: z.boolean(),
	is_active: z.boolean(),
	last_seen: z.string(),
});

export const answerSchema = z.object({
	id: z.string(),
	player_id: z.string(),
	text: z.string(),
	kind: z.enum(["correct", "fake", "filler"]),
	votes: z.array(z.string()),
	created_at: z.string(),
});

export const gameSettingsSchema = z.object({
	rounds: z.number(),
	max_players: z.number(),
	selected_categories: z.array(z.string()),
	time_limits: z.object({
		category_selection: z.number(),
		answer_writing: z.number(),
		voting: z.number(),
	}),
});

export const roundSchema = z.object({
	number: z.number(),
	category: z.string(),
	question: z.string(),
	question_id: z.string(),
	status: z.enum(["waiting", "active", "voting", "completed"]),
	current_turn: z
		.object({
			player_id: z.string(),
			status: z.enum(["waiting", "active", "ended"]),
			category: z.string().optional(),
		})
		.optional()
		.nullable(),
	answer_pool: z.object({
		correct_answer: z.string(),
		fake_answers: z.array(answerSchema),
		filler_answers: z.array(answerSchema),
		options: z.array(answerSchema),
	}),
});

export const gameSchema = z.object({
	id: z.string(),
	code: z.string(),
	status: z.enum(["waiting", "playing", "ended"]),
	players: z.array(playerSchema),
	rounds: z.array(roundSchema),
	settings: gameSettingsSchema,
	host_id: z.string(),
});

export const createGameFormSchema = z.object({
	name: z.string().trim().min(1, "Enter your display name."),
	rounds: z.number().int().min(1).max(20),
	selectedCategories: z
		.array(z.string())
		.min(1, "Select at least one category."),
});

export const joinGameFormSchema = z.object({
	name: z.string().trim().min(1, "Enter your display name."),
});

export const answerFormSchema = z.object({
	answer: z.string().trim().min(1, "Write a fake answer."),
});

export type Player = z.infer<typeof playerSchema>;
export type Answer = z.infer<typeof answerSchema>;
export type Round = z.infer<typeof roundSchema>;
export type Game = z.infer<typeof gameSchema>;
export type GameSettings = z.infer<typeof gameSettingsSchema>;
export type CreateGameForm = z.infer<typeof createGameFormSchema>;
export type JoinGameForm = z.infer<typeof joinGameFormSchema>;
export type AnswerForm = z.infer<typeof answerFormSchema>;
