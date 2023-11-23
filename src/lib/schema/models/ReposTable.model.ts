import { pgTable, serial, text, varchar, json, integer } from 'drizzle-orm/pg-core';

export const ReposTable = pgTable('repos', {
	id: serial('id').primaryKey(),
	repoID: integer('repo_id').notNull(),
	name: varchar('name', { length: 150 }).notNull(),
	description: text('description'),
	htmlURL: varchar('html_url', { length: 255 }).notNull(),
	homepage: varchar('homepage', { length: 255 }),
	tag: json('tag'),
	createdAt: varchar('created_at', { length: 30 }),
	updatedAt: varchar('updated_at', { length: 30 })
});
