This documentation is generated from com.marshmallow.anwork.app.cli.GithubReadmeDocumentationGenerator

#  *anwork*: ANWORK CLI commands
## Flags
- -c|--context (STRING name: The name of the persistence context): Set the persistence context
- -d|--debug: Turn on debug printing
- -n|--no-persist: Do not persist any task information
- -o|--output (STRING directory: The directory at which to output the persistence data): Set persistence output directory
# anwork *journal*: Journal commands...
## Commands
- anwork journal *show*: Show the entries in the journal for a task
- anwork journal *show-all*: Show all of the entries in the journal
# anwork *task*: Task commands...
## Commands
- anwork task *create*: Create a task
## Flags
- -e|--description (STRING description: The description of the task): The description of the task
- -p|--priority (NUMBER priority: The priority of the task): The priority of the task
## Commands
- anwork task *delete*: Delete a task
- anwork task *delete-all*: Delete all tasks
- anwork task *set-blocked*: Set a task as blocked
- anwork task *set-finished*: Set a task as finished
- anwork task *set-running*: Set a task as running
- anwork task *set-waiting*: Set a task as waiting
- anwork task *show*: Show all tasks
