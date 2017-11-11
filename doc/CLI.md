This documentation is generated from com.marshmallow.anwork.app.cli.GithubReadmeDocumentationGenerator

`anwork [-d|--debug] [-c|--context <name:STRING>] [-o|--output <directory:STRING>] [-n|--no-persist]` ... : ANWORK CLI commands
* `anwork summary <days:NUMBER>`
* * Show a summary of the past days of work

`anwork journal` ... : Journal commands...
* `anwork journal show <task-name:STRING>`
* * Show the entries in the journal for a task
* `anwork journal show-all`
* * Show all of the entries in the journal

`anwork journal task` ... : Task commands...
* `anwork journal task create [-e|--description <description:STRING>] [-p|--priority <priority:NUMBER>] <task-name:STRING>`
* * Create a task
* `anwork journal task delete <task-name:STRING>`
* * Delete a task
* `anwork journal task delete-all`
* * Delete all tasks
* `anwork journal task note <task-name:STRING> <note:STRING>`
* * Add a note to a task
* `anwork journal task set-blocked <task-name:STRING>`
* * Set a task as blocked
* `anwork journal task set-finished <task-name:STRING>`
* * Set a task as finished
* `anwork journal task set-priority <task-name:STRING> <priority:NUMBER>`
* * Set the priority of a task
* `anwork journal task set-running <task-name:STRING>`
* * Set a task as running
* `anwork journal task set-waiting <task-name:STRING>`
* * Set a task as waiting
* `anwork journal task show [-s|--short]`
* * Show all tasks
