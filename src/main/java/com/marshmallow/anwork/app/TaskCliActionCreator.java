package com.marshmallow.anwork.app;

import com.marshmallow.anwork.app.cli.Action;
import com.marshmallow.anwork.app.cli.ActionCreator;
import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

/**
 * This is a {@link ActionCreator} for the task commands in the ANWORK app.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskCliActionCreator implements ActionCreator {

  private static class SetStateCliAction extends TaskManagerCliAction implements Action {

    private TaskState taskState;

    public SetStateCliAction(TaskState taskState) {
      super();
      this.taskState = taskState;
    }

    @Override
    public void run(AnworkAppConfig config,
                    ArgumentValues flags,
                    ArgumentValues arguments,
                    TaskManager manager) {
      String taskName = arguments.getValue(TASK_NAME_ARGUMENT, ArgumentType.STRING);
      manager.setState(taskName, taskState);
    }
  }

  @Override
  public Action createAction(String commandName) {
    switch (commandName) {
      case "create":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          ArgumentValues arguments,
                          TaskManager manager) {
            String name = arguments.getValue(TASK_NAME_ARGUMENT, ArgumentType.STRING);
            String description = (flags.containsKey("e")
                                  ? flags.getValue("e", ArgumentType.STRING)
                                  : "");
            int priority = (int)(flags.containsKey("p")
                                 ? flags.getValue("p", ArgumentType.NUMBER)
                                 : Task.DEFAULT_PRIORITY);
            manager.createTask(name, description, priority);
            config.getDebugPrinter().accept("created task '" + name + "'");
          }
        };
      case "set-waiting":
        return new SetStateCliAction(TaskState.WAITING);
      case "set-blocked":
        return new SetStateCliAction(TaskState.BLOCKED);
      case "set-running":
        return new SetStateCliAction(TaskState.RUNNING);
      case "set-finished":
        return new SetStateCliAction(TaskState.FINISHED);
      case "delete":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          ArgumentValues arguments,
                          TaskManager manager) {
            String name = arguments.getValue(TASK_NAME_ARGUMENT, ArgumentType.STRING);
            manager.deleteTask(name);
            config.getDebugPrinter().accept("deleted task '" + name + "'");
          }
        };
      case "delete-all":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          ArgumentValues arguments,
                          TaskManager manager) {
            for (Task task : manager.getTasks()) {
              manager.deleteTask(task.getName());
            }
            config.getDebugPrinter().accept("deleted all tasks");
          }
        };
      case "show":
        return new TaskManagerCliAction() {
          @Override
          public void run(AnworkAppConfig config,
                          ArgumentValues flags,
                          ArgumentValues arguments,
                          TaskManager manager) {
            for (Task task : manager.getTasks()) {
              System.out.println(task);
            }
          }
        };
      default:
        return null; // error!
    }
  }
}
