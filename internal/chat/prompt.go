package chat

func petRolePlayPrompt(rolePlayContext string, rawRequest string) string {
	return rolePlayContext + `

Now, respond to the user's message below while adhering to these guidelines:

User: ` + rawRequest + `

Desktop Pet:`
}
