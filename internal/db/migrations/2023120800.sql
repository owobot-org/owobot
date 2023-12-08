/* change the string delimeter from comma to Unit Separator, to allow commas */
/* to be used in polls and the like                                          */
UPDATE reactions SET reaction = REPLACE(reaction, ',', X'1F');
UPDATE polls SET opt_emojis = REPLACE(opt_emojis, ',', X'1F');
UPDATE polls SET opt_text = REPLACE(opt_text, ',', X'1F');
UPDATE reaction_role_categories SET emoji = REPLACE(emoji, ',', X'1F');
UPDATE reaction_role_categories SET roles = REPLACE(roles, ',', X'1F');