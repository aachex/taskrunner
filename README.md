# taskrunner
The app represents a service that can be used to run and interrupt long-executing tasks. Note that the tasks' execution is only a simulation.
## Usage
To run the app locally, use the ```go run``` command.
```shell
go run main.go
```
or
```shell
go run .
```
## API
Now, there are 3 endpoints available:
|method | endpoint | description
|---------------|----------------|-----------|
|POST|/task/{name}/run|Runs a new task with a given name. If there is an executing or completed task with this name, the 409 Conflict error will be returned.|
|GET|/task/{name}/status|Returns a status of a task in JSON format. The status includes: task name, when it was created, status itself (executing/completed), and execution time. If the task was not found, the 404 Not Found error will be returned.|
|DELETE|/task/{name}/rm|Removes a task with given name. If a task is running currently, it will be interrupted and then deleted. If the task was not found, the 404 Not Found error will be returned.|
