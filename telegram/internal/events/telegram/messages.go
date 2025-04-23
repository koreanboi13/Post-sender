package telegram

const msgHelp = `
	I can send u latest post from blogAtor
	commands:
	/subscribe - Subscribe to new posts from blogator
	/unsubscribe - Unsubscribe from new posts
`

const msgHello = "Hi there! 👾\n\n" + msgHelp

const (
	msgUnknownCommand      = "Unknown command 🤔"
	msgNoPosts             = "No posts found 🤷‍♂️"
	msgAlreadySubscribed   = "You are already subscribed to the blog"
	msgSubscribed          = "You are subscribed to the blog"
	msgNotSubscribed       = "You are not subscribed to the blog"
	msgUnsubscribedSuccess = "You are unsubscribed from the blog"
)
